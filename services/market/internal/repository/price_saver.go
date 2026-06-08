package repository

import (
    "context"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/google/uuid"
)

type PriceSaver struct {
    pool *pgxpool.Pool
}

func NewPriceSaver(pool *pgxpool.Pool) *PriceSaver {
    return &PriceSaver{pool: pool}
}

func (s *PriceSaver) SaveNewPrice(ctx context.Context, assetID string, price, change, changePercent float64) error {
    _, err := s.pool.Exec(ctx,
        `INSERT INTO market_prices (id, asset_id, price_usdt, change_usdt, change_percent)
         VALUES ($1, $2, $3, $4, $5)`,
        uuid.New().String(), assetID, price, change, changePercent,
    )
    return err
}