CREATE TABLE market_prices (
    id UUID PRIMARY KEY,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    price_usdt DECIMAL(15,6) NOT NULL,
    change_usdt DECIMAL(15,6) NOT NULL DEFAULT 0,
    change_percent DECIMAL(10,4) NOT NULL DEFAULT 0,
    recorded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_market_prices_asset_time ON market_prices(asset_id, recorded_at DESC);