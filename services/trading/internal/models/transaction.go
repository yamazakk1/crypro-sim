package models

type Transaction struct {
    ID          string  `json:"id"`
    UserID      string  `json:"user_id"`
    AssetID     string  `json:"asset_id"`
    AssetSymbol string  `json:"asset_symbol"`
    Type        string  `json:"type"`
    Quantity    float64 `json:"quantity"`
    AssetPrice  float64 `json:"asset_price"`
    TotalUSDT   float64 `json:"total_usdt"`
    CreatedAt   string  `json:"created_at"`
}