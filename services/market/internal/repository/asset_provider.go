package repository

import (
    "context"

    "github.com/jackc/pgx/v5/pgxpool"
    "crypto-simulator/services/market/internal/service"
)

type AssetProvider struct {
    pool *pgxpool.Pool
}

func NewAssetProvider(pool *pgxpool.Pool) *AssetProvider {
    return &AssetProvider{pool: pool}
}

func (p *AssetProvider) GetLastPrices(ctx context.Context) ([]*service.CachedPrice, error) {
    query := `
        SELECT a.id, a.symbol, COALESCE(mp.price_usdt, a.initial_price, 100.0)
        FROM assets a
        LEFT JOIN LATERAL (
            SELECT price_usdt FROM market_prices 
            WHERE asset_id = a.id 
            ORDER BY recorded_at DESC 
            LIMIT 1
        ) mp ON true
        WHERE a.is_active = true
    `
    rows, err := p.pool.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var prices []*service.CachedPrice
    for rows.Next() {
        var cp service.CachedPrice
        rows.Scan(&cp.Id, &cp.Symbol, &cp.Price)
        cp.Trend = service.StableTrend
        prices = append(prices, &cp)
    }
    return prices, nil
}