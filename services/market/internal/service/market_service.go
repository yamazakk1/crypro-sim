package service

import (
    "context"
    "crypto-simulator/services/market/internal/models"
)

type PriceRepo interface {
    GetCurrentPrices(ctx context.Context) ([]models.PriceUpdate, error)
    GetPriceHistory(ctx context.Context, assetID, from, to string) ([]models.PricePoint, error)
}

type MarketService struct {
    repo PriceRepo
}

func NewMarketService(repo PriceRepo) *MarketService {
    return &MarketService{repo: repo}
}

func (s *MarketService) GetCurrentPrices(ctx context.Context) ([]models.PriceUpdate, error) {
    return s.repo.GetCurrentPrices(ctx)
}

func (s *MarketService) GetPriceHistory(ctx context.Context, assetID, from, to string) ([]models.PricePoint, error) {
    return s.repo.GetPriceHistory(ctx, assetID, from, to)
}