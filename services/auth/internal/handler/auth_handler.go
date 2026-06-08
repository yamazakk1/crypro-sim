package handler

import (
	"context"
	"errors"
	"log"

	pb "crypto-simulator/pkg/pb/auth"
	"crypto-simulator/services/auth/helpers"
	"crypto-simulator/services/auth/internal/models"
	"crypto-simulator/services/auth/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetMe(ctx context.Context, userId string) (*models.GetMeUser, error)
}

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("auth-handler: Register called: email=%s, username=%s", req.GetEmail(), req.GetUsername())

	username := req.GetUsername()
	if username == "" {
		log.Println("auth-handler: Register: empty username")
		return nil, status.Error(codes.InvalidArgument, "username can not be empty")
	}

	email := req.GetEmail()
	if !helpers.IsEmailValid(email) {
		log.Printf("auth-handler: Register: invalid email: %s", email)
		return nil, status.Error(codes.InvalidArgument, "invalid email input")
	}

	password := req.GetPassword()
	if err := helpers.IsPasswordValid(password); err != nil {
		log.Printf("auth-handler: Register: weak password")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Println("auth-handler: Register: calling service")
	resp, err := h.AuthService.Register(ctx, username, email, password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyExists), errors.Is(err, service.ErrUsernameAlreadyExists):
			log.Printf("auth-handler: Register: already exists: %v", err)
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		log.Printf("auth-handler: Register: internal error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("auth-handler: Register: success, email=%s", email)
	return &pb.RegisterResponse{Token: resp}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("auth-handler: Login called: email=%s", req.GetEmail())

	email := req.GetEmail()
	if !helpers.IsEmailValid(email) {
		log.Printf("auth-handler: Login: invalid email: %s", email)
		return nil, status.Error(codes.InvalidArgument, "invalid email input")
	}

	password := req.GetPassword()
	if err := helpers.IsPasswordValid(password); err != nil {
		log.Printf("auth-handler: Login: weak password")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Println("auth-handler: Login: calling service")
	resp, err := h.AuthService.Login(ctx, email, password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			log.Printf("auth-handler: Login: user not found: %s", email)
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, service.ErrWrongPassword):
			log.Printf("auth-handler: Login: wrong password for email=%s", email)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		log.Printf("auth-handler: Login: internal error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("auth-handler: Login: success, email=%s", email)
	return &pb.LoginResponse{Token: resp}, nil
}

func (h *AuthHandler) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.GetMeResponse, error) {
	log.Printf("auth-handler: GetMe called: user_id=%s", req.GetUserId())

	userId := req.GetUserId()
	resp, err := h.AuthService.GetMe(ctx, userId)
	if err != nil {
		log.Printf("auth-handler: GetMe: error: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("auth-handler: GetMe: success, username=%s", resp.Username)
	return &pb.GetMeResponse{
		Username:    resp.Username,
		Email:       resp.Email,
		BalanceUsdt: float32(resp.Balance_usdt),
		Role: string(resp.Role),
	}, nil
}