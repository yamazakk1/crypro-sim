package models

type PortfolioItem struct {
    AssetID        string  `json:"asset_id"`
    Symbol         string  `json:"symbol"`
    Quantity       float64 `json:"quantity"`
    AvgBuyPrice    float64 `json:"avg_buy_price"`
    CurrentPrice   float64 `json:"current_price"`
    TotalValue     float64 `json:"total_value"`
    ProfitLoss     float64 `json:"profit_loss"`
    ProfitLossPct  float64 `json:"profit_loss_percent"`
}

type Portfolio struct {
    Items           []PortfolioItem `json:"items"`
    TotalValue      float64         `json:"total_value_usdt"`
    BalanceUSDT     float64         `json:"balance_usdt"`
    TotalProfitLoss float64         `json:"total_profit_loss"`
}