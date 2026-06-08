package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/grpc/status"

	pbAuth "crypto-simulator/pkg/pb/auth"
	"crypto-simulator/services/gateway/internal/middleware"
)

func (h *GatewayHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	log.Println("gateway: HandleRegister called")
	var req RegisterRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.Username == "" {
		log.Println("gateway: HandleRegister: empty username")
		writeError(w, "empty username", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		log.Println("gateway: HandleRegister: empty email")
		writeError(w, "empty email", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		log.Println("gateway: HandleRegister: empty password")
		writeError(w, "empty password", http.StatusBadRequest)
		return
	}

	log.Printf("gateway: HandleRegister: calling Auth Service, email=%s, username=%s", req.Email, req.Username)
	resp, err := h.AuthServiceClient.Register(r.Context(), &pbAuth.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		st := status.Convert(err)
		log.Printf("gateway: HandleRegister: Auth Service error: %v", st.Message())
		writeError(w, st.Message(), grpcToHTTP(st.Code()))
		return
	}

	log.Printf("gateway: HandleRegister: success, email=%s", req.Email)
	writeJson(w, http.StatusCreated, RegisterResponse{Token: resp.GetToken()})
}

func (h *GatewayHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("gateway: HandleLogin called")
	var req LoginRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.Email == "" {
		log.Println("gateway: HandleLogin: empty email")
		writeError(w, "empty email", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		log.Println("gateway: HandleLogin: empty password")
		writeError(w, "empty password", http.StatusBadRequest)
		return
	}

	log.Printf("gateway: HandleLogin: calling Auth Service, email=%s", req.Email)
	resp, err := h.AuthServiceClient.Login(r.Context(), &pbAuth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		st := status.Convert(err)
		log.Printf("gateway: HandleLogin: Auth Service error: %v", st.Message())
		writeError(w, st.Message(), grpcToHTTP(st.Code()))
		return
	}

	log.Printf("gateway: HandleLogin: success, email=%s", req.Email)
	writeJson(w, http.StatusOK, LoginResponse{Token: resp.GetToken()})
}

func (h *GatewayHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	log.Println("gateway: GetMe called")
	ctx := r.Context()
	userId := ctx.Value(middleware.UserIDKey)

	if userId == nil || userId.(string) == "" {
		log.Println("gateway: GetMe: empty user_id in context")
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("gateway: GetMe: calling Auth Service, user_id=%s", userId.(string))
	resp, err := h.AuthServiceClient.GetMe(ctx, &pbAuth.GetMeRequest{UserId: userId.(string)})

	if err != nil {
		st := status.Convert(err)
		log.Printf("gateway: GetMe: Auth Service error: %v", st.Message())
		writeError(w, st.Message(), grpcToHTTP(st.Code()))
		return
	}

	log.Printf("gateway: GetMe: success, username=%s", resp.GetUsername())
	writeJson(w, http.StatusOK, GetMeResponse{
		Username:    resp.GetUsername(),
		Email:       resp.GetEmail(),
		BalanceUsdt: float64(resp.GetBalanceUsdt()),
		Role:        resp.GetRole(),
	})
}
