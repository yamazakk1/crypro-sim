package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserRoleKey contextKey = "role"

func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				log.Println("middleware: JWTAuth: missing or invalid Authorization header")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(header, "Bearer ")

			log.Printf("middleware: JWTAuth: parsing token, prefix=%s...", tokenString[:min(20, len(tokenString))])
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				log.Printf("middleware: JWTAuth: invalid token: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Println("middleware: JWTAuth: failed to parse claims")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userID, _ := claims["user_id"].(string)
			role, _ := claims["role"].(string)

			log.Printf("middleware: JWTAuth: authenticated user_id=%s, role=%s", userID, role)
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(UserRoleKey).(string)
		if role != "admin" {
			log.Printf("middleware: RequireAdmin: access denied, role=%s", role)
			http.Error(w, `{"error":"admin only"}`, http.StatusForbidden)
			return
		}
		log.Println("middleware: RequireAdmin: access granted")
		next.ServeHTTP(w, r)
	})
}

