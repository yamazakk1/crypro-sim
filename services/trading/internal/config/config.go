package config

import (
    "fmt"
    "os"
)

const DefaultTradingServicePort = "8084"

type TradingServiceConfig struct {
    Port   string
    DB_DSN string
}

func NewTradingServiceConfig() *TradingServiceConfig {
    return &TradingServiceConfig{
        Port:   GetEnvOrDefault("TRADING_PORT", DefaultTradingServicePort),
        DB_DSN: BuildDSN(),
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