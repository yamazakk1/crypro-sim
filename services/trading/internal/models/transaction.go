package models

type Transaction struct {
    ID          string
    UserID      string
    AssetID     string
    AssetSymbol string
    Type        string 
    Quantity    float64
    AssetPrice  float64
    TotalUSDT   float64
    CreatedAt   string
}