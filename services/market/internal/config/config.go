package config

import (
	"fmt"
	"os"
)

const DefaultMarketServicePort = "8083"

type MarketServiceConfig struct {
	Port      string
	DB_DSN    string
	RedisAddr string
}

func NewMarketServiceConfig() *MarketServiceConfig {
	return &MarketServiceConfig{
		Port:      GetEnvOrDefault("MARKET_PORT", DefaultMarketServicePort),
		DB_DSN:    BuildDSN(),
		RedisAddr: GetEnvOrDefault("REDIS_ADDR", "localhost:6379"),
	}
}

func BuildDSN() string {
	if dsn := os.Getenv("DB_DSN"); dsn != "" {
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
