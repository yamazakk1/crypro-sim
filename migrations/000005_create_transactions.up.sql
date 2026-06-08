CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL,
    quantity DECIMAL(20,8) NOT NULL,
    asset_price DECIMAL(15,6) NOT NULL,
    total_usdt DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_user_time ON transactions(user_id, created_at DESC);
CREATE INDEX idx_transactions_asset_time ON transactions(asset_id, created_at DESC);