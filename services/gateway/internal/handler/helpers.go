package handler

import (
	"crypto-simulator/services/gateway/internal/middleware"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func writeJson(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func grpcToHTTP(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func (h *GatewayHandler) isAdmin(r *http.Request) bool {
	role, _ := r.Context().Value(middleware.UserRoleKey).(string)
	return role == "admin"
}
