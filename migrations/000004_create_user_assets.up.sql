CREATE TABLE user_assets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    quantity DECIMAL(20,8) NOT NULL,
    avg_buy_price DECIMAL(15,6) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, asset_id)
);

CREATE INDEX idx_user_assets_user_id ON user_assets(user_id);