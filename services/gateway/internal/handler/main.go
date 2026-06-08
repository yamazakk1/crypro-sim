package handler

import (
	pbAsset "crypto-simulator/pkg/pb/asset"
	pbAuth "crypto-simulator/pkg/pb/auth"
	pbmarket "crypto-simulator/pkg/pb/market"
	pbTrading "crypto-simulator/pkg/pb/trading"
)

type GatewayHandler struct {
	AuthServiceClient   pbAuth.AuthServiceClient
	AssetServiceClient  pbAsset.AssetServiceClient
	MarketServiceClient pbmarket.MarketServiceClient
	TradingServiceClient       pbTrading.TradingServiceClient
}

func NewGatewayHandler(authService pbAuth.AuthServiceClient, assetService pbAsset.AssetServiceClient, marketService pbmarket.MarketServiceClient, tradingService pbTrading.TradingServiceClient) *GatewayHandler {
	return &GatewayHandler{
		AuthServiceClient:   authService,
		AssetServiceClient:  assetService,
		MarketServiceClient: marketService,
		TradingServiceClient: tradingService,
	}
}
