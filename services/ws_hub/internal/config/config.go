package config

import "os"

const DefaultWSHubPort = "8085"
const DefaultRedisAddr = "localhost:6379"

type WSHubConfig struct {
    Port      string
    RedisAddr string
}

func NewWSHubConfig() *WSHubConfig {
    return &WSHubConfig{
        Port:      GetEnvOrDefault("WS_HUB_PORT", DefaultWSHubPort),
        RedisAddr: GetEnvOrDefault("REDIS_ADDR", DefaultRedisAddr),
    }
}

func GetEnvOrDefault(key, defaultValue string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultValue
}