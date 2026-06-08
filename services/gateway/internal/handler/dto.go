package handler

// ─── Auth ───

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type GetMeResponse struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	BalanceUsdt float64   `json:"balance_usdt"`
	Role        string `json:"role"`
}

// ─── Asset ───

type Asset struct {
	Id       string `json:"id"`
	Symbol   string `json:"symbol"`
	FullName string `json:"full_name"`
	IsActive bool   `json:"is_active"`
}

type ListAssetsResponse struct {
	Assets []Asset `json:"assets"`
}

type GetAssetRequest struct {
	Id string `json:"id"`
}

type CreateAssetRequest struct {
    Symbol       string  `json:"symbol"`
    FullName     string  `json:"full_name"`
    InitialPrice float64 `json:"initial_price"`
}

type CreateAssetResponse struct {
	Id       string `json:"id"`
	Symbol   string `json:"symbol"`
	FullName string `json:"full_name"`
	IsActive bool   `json:"is_active"`
}

type UpdateAssetRequest struct {
	Id       string `json:"id"`
	FullName string `json:"full_name"`
}

type DeactivateAssetRequest struct {
	Id string `json:"id"`
}

type ActivateAssetRequest struct {
	Id string `json:"id"`
}

type DeleteAssetRequest struct {
	Id string `json:"id"`
}

type DeleteAssetResponse struct {
	Success bool `json:"success"`
}

type PriceUpdate struct {
    AssetID       string `json:"asset_id"`
    Symbol        string `json:"symbol"`
    PriceUsdt     string `json:"price_usdt"`
    ChangeUsdt    string `json:"change_usdt"`
    ChangePercent string `json:"change_percent"`
    Timestamp     string `json:"timestamp"`
}

type GetCurrentPricesResponse struct {
    Prices []PriceUpdate `json:"prices"`
}