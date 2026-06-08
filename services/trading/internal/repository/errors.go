package repository

import "errors"

var (
    ErrInsufficientBalance = errors.New("insufficient balance")
    ErrInsufficientAssets  = errors.New("insufficient assets")
    ErrAssetNotFound       = errors.New("asset not found")
)