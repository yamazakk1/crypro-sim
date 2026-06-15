package service

import (
	"context"
	"fmt"
	"log"

	"crypto-simulator/services/trading/internal/models"
	"crypto-simulator/services/trading/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TradingService struct {
	pool          *pgxpool.Pool
	portfolioRepo *repository.PortfolioRepo
	txRepo        *repository.TransactionRepo
}

func NewTradingService(pool *pgxpool.Pool, pRepo *repository.PortfolioRepo, tRepo *repository.TransactionRepo) *TradingService {
	return &TradingService{pool: pool, portfolioRepo: pRepo, txRepo: tRepo}
}

func (s *TradingService) GetCurrentPrice(ctx context.Context, tx pgx.Tx, assetID string) (float64, error) {
	var price float64
	err := tx.QueryRow(ctx,
		`SELECT COALESCE(
            (SELECT price_usdt FROM market_prices WHERE asset_id = $1 ORDER BY recorded_at DESC LIMIT 1),
            (SELECT initial_price FROM assets WHERE id = $1), 0)`,
		assetID).Scan(&price)
	return price, err
}

func (s *TradingService) Buy(ctx context.Context, userID, assetID string, quantity float64) (*models.Transaction, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	price, err := s.GetCurrentPrice(ctx, tx, assetID)
	if err != nil || price == 0 {
		return nil, fmt.Errorf("no price for asset %s", assetID)
	}

	totalCost := quantity * price

	balance, err := s.portfolioRepo.GetBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}
	if balance < totalCost {
		return nil, repository.ErrInsufficientBalance
	}

	if err := s.portfolioRepo.UpdateBalance(ctx, tx, userID, -totalCost); err != nil {
		return nil, fmt.Errorf("update balance: %w", err)
	}

	oldQty, oldAvg, _ := s.portfolioRepo.GetUserAsset(ctx, userID, assetID)
	newQty := oldQty + quantity
	newAvg := ((oldQty * oldAvg) + totalCost) / newQty

	if err := s.portfolioRepo.UpsertAsset(ctx, tx, userID, assetID, newQty, newAvg); err != nil {
		return nil, fmt.Errorf("upsert asset: %w", err)
	}

	t, err := s.txRepo.Create(ctx, tx, userID, assetID, "buy", quantity, price, totalCost)
	if err != nil {
		return nil, fmt.Errorf("create tx: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	log.Printf("trading: BUY user=%s asset=%s qty=%.4f price=%.2f total=%.2f", userID, assetID, quantity, price, totalCost)
	return t, nil
}

func (s *TradingService) Sell(ctx context.Context, userID, assetID string, quantity float64) (*models.Transaction, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	oldQty, oldAvg, err := s.portfolioRepo.GetUserAsset(ctx, userID, assetID)
	if err != nil {
		return nil, fmt.Errorf("get user asset: %w", err)
	}
	if oldQty < quantity {
		return nil, repository.ErrInsufficientAssets
	}

	price, err := s.GetCurrentPrice(ctx, tx, assetID)
	if err != nil || price == 0 {
		return nil, fmt.Errorf("no price for asset %s", assetID)
	}

	totalRevenue := quantity * price

	if err := s.portfolioRepo.UpdateBalance(ctx, tx, userID, totalRevenue); err != nil {
		return nil, fmt.Errorf("update balance: %w", err)
	}

	newQty := oldQty - quantity
	if newQty <= 0 {
		if err := s.portfolioRepo.DeleteAsset(ctx, tx, userID, assetID); err != nil {
			return nil, fmt.Errorf("delete asset: %w", err)
		}
	} else {
		if err := s.portfolioRepo.UpsertAsset(ctx, tx, userID, assetID, newQty, oldAvg); err != nil {
			return nil, fmt.Errorf("upsert asset: %w", err)
		}
	}

	t, err := s.txRepo.Create(ctx, tx, userID, assetID, "sell", quantity, price, totalRevenue)
	if err != nil {
		return nil, fmt.Errorf("create tx: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	log.Printf("trading: SELL user=%s asset=%s qty=%.4f price=%.2f total=%.2f", userID, assetID, quantity, price, totalRevenue)
	return t, nil
}

func (s *TradingService) GetPortfolio(ctx context.Context, userID string) (*models.Portfolio, error) {
	return s.portfolioRepo.GetPortfolio(ctx, userID)
}

func (s *TradingService) GetTransactions(ctx context.Context, userID string) ([]models.Transaction, error) {
	return s.txRepo.GetByUserID(ctx, userID)
}

func (s *TradingService) AddBalance(ctx context.Context, userID string, amount float64) (float64, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	if err := s.portfolioRepo.UpdateBalance(ctx, tx, userID, amount); err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return s.portfolioRepo.GetBalance(ctx, userID)
}
