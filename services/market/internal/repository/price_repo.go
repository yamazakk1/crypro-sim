package repository

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "crypto-simulator/services/market/internal/models"
)

type PriceRepo struct {
    pool *pgxpool.Pool
}

func NewPriceRepo(pool *pgxpool.Pool) *PriceRepo {
    return &PriceRepo{pool: pool}
}

func (r *PriceRepo) GetCurrentPrices(ctx context.Context) ([]models.PriceUpdate, error) {
    query := `
        SELECT DISTINCT ON (mp.asset_id) 
            mp.asset_id, a.symbol, mp.price_usdt, mp.change_usdt, mp.change_percent, mp.recorded_at
        FROM market_prices mp
        JOIN assets a ON mp.asset_id = a.id
        WHERE a.is_active = true
        ORDER BY mp.asset_id, mp.recorded_at DESC
    `
    rows, err := r.pool.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var prices []models.PriceUpdate
    for rows.Next() {
        var p models.PriceUpdate
        var ts time.Time
        if err := rows.Scan(&p.AssetID, &p.Symbol, &p.Price, &p.Change, &p.ChangePercent, &ts); err != nil {
            return nil, err
        }
        p.Timestamp = ts.Format(time.RFC3339)
        prices = append(prices, p)
    }
    return prices, nil
}

func (r *PriceRepo) GetPriceHistory(ctx context.Context, assetID, from, to string) ([]models.PricePoint, error) {
    query := `
        SELECT price_usdt, price_usdt, price_usdt, price_usdt, recorded_at
        FROM market_prices
        WHERE asset_id = $1 AND recorded_at BETWEEN $2 AND $3
        ORDER BY recorded_at ASC
    `
    rows, err := r.pool.Query(ctx, query, assetID, from, to)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var points []models.PricePoint
    for rows.Next() {
        var p models.PricePoint
        var ts time.Time
        rows.Scan(&p.Open, &p.High, &p.Low, &p.Close, &ts)
        p.Timestamp = ts.Format(time.RFC3339)
        points = append(points, p)
    }
    return points, nil
}