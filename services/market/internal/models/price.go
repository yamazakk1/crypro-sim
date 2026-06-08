package models

type PriceUpdate struct {
    AssetID       string
    Symbol        string
    Price         float64
    Change        float64
    ChangePercent float64
    Timestamp     string
}

type PricePoint struct {
    Timestamp string
    Open      float64
    High      float64
    Low       float64
    Close     float64
}