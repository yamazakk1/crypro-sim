package repository

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "crypto-simulator/services/trading/internal/models"
)

type PortfolioRepo struct {
    pool *pgxpool.Pool
}

func NewPortfolioRepo(pool *pgxpool.Pool) *PortfolioRepo {
    return &PortfolioRepo{pool: pool}
}

func (r *PortfolioRepo) GetBalance(ctx context.Context, userID string) (float64, error) {
    var balance float64
    err := r.pool.QueryRow(ctx, `SELECT balance_usdt FROM users WHERE id = $1`, userID).Scan(&balance)
    return balance, err
}

func (r *PortfolioRepo) UpdateBalance(ctx context.Context, tx pgx.Tx, userID string, delta float64) error {
    _, err := tx.Exec(ctx, `UPDATE users SET balance_usdt = balance_usdt + $1 WHERE id = $2`, delta, userID)
    return err
}

func (r *PortfolioRepo) GetUserAsset(ctx context.Context, userID, assetID string) (quantity, avgPrice float64, err error) {
    err = r.pool.QueryRow(ctx,
        `SELECT quantity, avg_buy_price FROM user_assets WHERE user_id = $1 AND asset_id = $2`,
        userID, assetID).Scan(&quantity, &avgPrice)
    if err == pgx.ErrNoRows {
        return 0, 0, nil
    }
    return
}

func (r *PortfolioRepo) UpsertAsset(ctx context.Context, tx pgx.Tx, userID, assetID string, quantity, avgPrice float64) error {
    _, err := tx.Exec(ctx,
        `INSERT INTO user_assets (id, user_id, asset_id, quantity, avg_buy_price)
         VALUES (gen_random_uuid(), $1, $2, $3, $4)
         ON CONFLICT (user_id, asset_id) DO UPDATE
         SET quantity = $3, avg_buy_price = $4, updated_at = NOW()`,
        userID, assetID, quantity, avgPrice)
    return err
}

func (r *PortfolioRepo) DeleteAsset(ctx context.Context, tx pgx.Tx, userID, assetID string) error {
    _, err := tx.Exec(ctx, `DELETE FROM user_assets WHERE user_id = $1 AND asset_id = $2`, userID, assetID)
    return err
}

func (r *PortfolioRepo) GetPortfolio(ctx context.Context, userID string) (*models.Portfolio, error) {
    balance, err := r.GetBalance(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("get balance: %w", err)
    }

    rows, err := r.pool.Query(ctx,
        `SELECT ua.asset_id, a.symbol, ua.quantity, ua.avg_buy_price,
                COALESCE(
                    (SELECT price_usdt FROM market_prices WHERE asset_id = ua.asset_id ORDER BY recorded_at DESC LIMIT 1),
                    (SELECT initial_price FROM assets WHERE id = ua.asset_id), 0
                ) as current_price
         FROM user_assets ua
         JOIN assets a ON ua.asset_id = a.id
         WHERE ua.user_id = $1 AND ua.quantity > 0`, userID)
    if err != nil {
        return nil, fmt.Errorf("query portfolio: %w", err)
    }
    defer rows.Close()

    portfolio := &models.Portfolio{BalanceUSDT: balance}
    for rows.Next() {
        var item models.PortfolioItem
        if err := rows.Scan(&item.AssetID, &item.Symbol, &item.Quantity, &item.AvgBuyPrice, &item.CurrentPrice); err != nil {
            return nil, fmt.Errorf("scan portfolio item: %w", err)
        }
        item.TotalValue = item.Quantity * item.CurrentPrice
        item.ProfitLoss = item.TotalValue - (item.Quantity * item.AvgBuyPrice)
        if item.AvgBuyPrice > 0 {
            item.ProfitLossPct = ((item.CurrentPrice - item.AvgBuyPrice) / item.AvgBuyPrice) * 100
        }
        portfolio.Items = append(portfolio.Items, item)
        portfolio.TotalValue += item.TotalValue
        portfolio.TotalProfitLoss += item.ProfitLoss
    }
    return portfolio, nil
}