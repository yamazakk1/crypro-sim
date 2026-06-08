package config

import "os"

type GatewayServiceConfig struct {
	Port      string
	JWTSecret string
	AuthServiceAddr string 
	AssetServiceAddr string
	MarketServiceAddr string
	TradingServiceAddr string
}

const DefaultPort = ":8080"
const DefaultJWTSecret = "dev-secret-yo"
const DefaultAuthAddr = "localhost:8081"
const DefaultAssetAddr = "localhost:8082"
const DefaultMarketAddr = "localhost:8083" 
const DefaultTradingAddr = "localhost:8084"

func NewGatewayServiceConfig() *GatewayServiceConfig{
	return &GatewayServiceConfig{
		Port: GetEnvOrDefault("GATEWAY_PORT", DefaultPort),
		JWTSecret: GetEnvOrDefault("JWT_SECRET", DefaultJWTSecret),
		AuthServiceAddr: GetEnvOrDefault("AUTH_ADDR",DefaultAuthAddr),
		AssetServiceAddr: GetEnvOrDefault("ASSET_ADDR", DefaultAssetAddr),
		MarketServiceAddr: GetEnvOrDefault("MARKET_ADDR", DefaultMarketAddr),
		TradingServiceAddr: GetEnvOrDefault("TRADING_ADDR", DefaultTradingAddr),
	}
}

func GetEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != ""{
		return val
	}
	return defaultValue
}