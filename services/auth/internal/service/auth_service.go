package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"crypto-simulator/services/auth/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	Create(ctx context.Context, id, username, email string, passwordHash []byte) error
	EmailExists(ctx context.Context, email string) (bool, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetMeById(ctx context.Context, userId string) (*models.GetMeUser, error)
}

type AuthService struct {
	AuthRepo
	secret string
}

func NewAuthService(authRepo AuthRepo, secret string) *AuthService {
	log.Println("auth-service: initialized")
	return &AuthService{AuthRepo: authRepo, secret: secret}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (string, error) {
	log.Printf("auth-service: Register called: email=%s, username=%s", email, username)

	if isExists, err := s.AuthRepo.EmailExists(ctx, email); err != nil {
		log.Printf("auth-service: Register: EmailExists error: %v", err)
		return "", err
	} else if isExists {
		log.Printf("auth-service: Register: email already exists: %s", email)
		return "", ErrEmailAlreadyExists
	}

	if isExists, err := s.AuthRepo.UsernameExists(ctx, username); err != nil {
		log.Printf("auth-service: Register: UsernameExists error: %v", err)
		return "", err
	} else if isExists {
		log.Printf("auth-service: Register: username already exists: %s", username)
		return "", ErrUsernameAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("auth-service: Register: bcrypt error: %v", err)
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	userId := uuid.NewSHA1(uuid.NameSpaceOID, []byte(email)).String()
	log.Printf("auth-service: Register: generated user_id=%s", userId)

	err = s.AuthRepo.Create(ctx, userId, username, email, passwordHash)
	if err != nil {
		log.Printf("auth-service: Register: Create error: %v", err)
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := generateToken(userId, email, username, s.secret, models.UserRole)
	if err != nil {
		log.Printf("auth-service: Register: generateToken error: %v", err)
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	log.Printf("auth-service: Register: success, user_id=%s", userId)
	return token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	log.Printf("auth-service: Login called: email=%s", email)

	if isExists, err := s.AuthRepo.EmailExists(ctx, email); err != nil {
		log.Printf("auth-service: Login: EmailExists error: %v", err)
		return "", err
	} else if !isExists {
		log.Printf("auth-service: Login: user not found: %s", email)
		return "", ErrUserNotFound
	}

	user, err := s.AuthRepo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("auth-service: Login: GetUserByEmail error: %v", err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Passwordhash), []byte(password))
	if err != nil {
		log.Printf("auth-service: Login: wrong password for email=%s", email)
		return "", ErrWrongPassword
	}

	token, err := generateToken(user.UserID, user.Email, user.Username, s.secret, user.Role)
	if err != nil {
		log.Printf("auth-service: Login: generateToken error: %v", err)
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	log.Printf("auth-service: Login: success, user_id=%s", user.UserID)
	return token, nil
}

func (s *AuthService) GetMe(ctx context.Context, userId string) (*models.GetMeUser, error) {
	log.Printf("auth-service: GetMe called: user_id=%s", userId)

	user, err := s.AuthRepo.GetMeById(ctx, userId)
	if err != nil {
		log.Printf("auth-service: GetMe: GetMeById error: %v", err)
		return nil, err
	}

	log.Printf("auth-service: GetMe: success, username=%s", user.Username)
	return user, nil
}

func generateToken(userId, email, username, secret string, role models.Role) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userId,
		"role":     role,
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}