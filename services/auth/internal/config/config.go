package config

import (
	"fmt"
	"os"
)

const DefaultAuthServicePort = "8081"
const DefaultJWTSecret = "dev-secret-yo"

type AuthServiceConfig struct {
	Port      string
	DB_DSN    string
	JWTSecret string
}

func NewAuthServiceConfig() *AuthServiceConfig {
	return &AuthServiceConfig{
		Port:      GetEnvOrDefault("AUTH_PORT", DefaultAuthServicePort),
		JWTSecret: GetEnvOrDefault("JWT_SECRET", DefaultJWTSecret),
		DB_DSN:    BuildDSN(),
	}
}

func BuildDSN() string {
	if dsn := os.Getenv("DB_DSN"); dsn != ""{
		return dsn
	}
	return fmt.Sprintf(
		"postgres://%s:%s@postgres:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
}

func GetEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
