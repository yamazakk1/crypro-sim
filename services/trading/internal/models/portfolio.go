package models

type PortfolioItem struct {
    AssetID        string
    Symbol         string
    Quantity       float64
    AvgBuyPrice    float64
    CurrentPrice   float64
    TotalValue     float64
    ProfitLoss     float64
    ProfitLossPct  float64
}

type Portfolio struct {
    Items          []PortfolioItem
    TotalValue     float64
    BalanceUSDT    float64
    TotalProfitLoss float64
}