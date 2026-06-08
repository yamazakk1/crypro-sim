package service

import (
    "context"
    "log"

    "crypto-simulator/services/asset/internal/models"
)

type AssetRepo interface {
    List(ctx context.Context) ([]*models.Asset, error)
    GetByID(ctx context.Context, id string) (*models.Asset, error)
    Create(ctx context.Context, symbol, fullname string, initPrice float64) (*models.Asset, error)
    Update(ctx context.Context, id, fullname string) error
    SetActive(ctx context.Context, id string, active bool) error
    Delete(ctx context.Context, id string) error
}

type AssetService struct {
    repo AssetRepo
}

func NewAssetService(repo AssetRepo) *AssetService {
    log.Println("asset-service: initialized")
    return &AssetService{repo: repo}
}

func (s *AssetService) List(ctx context.Context) ([]*models.Asset, error) {
    log.Println("asset-service: List called")
    assets, err := s.repo.List(ctx)
    if err != nil {
        log.Printf("asset-service: List: error: %v", err)
        return nil, err
    }
    log.Printf("asset-service: List: success, count=%d", len(assets))
    return assets, nil
}

func (s *AssetService) GetByID(ctx context.Context, id string) (*models.Asset, error) {
    log.Printf("asset-service: GetByID called: id=%s", id)
    asset, err := s.repo.GetByID(ctx, id)
    if err != nil {
        log.Printf("asset-service: GetByID: error: %v", err)
        return nil, err
    }
    log.Printf("asset-service: GetByID: success, symbol=%s", asset.Symbol)
    return asset, nil
}

func (s *AssetService) Create(ctx context.Context, symbol, fullname string, initPrice float64) (*models.Asset, error) {
    log.Printf("asset-service: Create called: symbol=%s, fullname=%s", symbol, fullname)
    if symbol == "" || fullname == "" {
        log.Println("asset-service: Create: empty fields")
        return nil, ErrEmptyFields
    }
    asset, err := s.repo.Create(ctx, symbol, fullname, initPrice)
    if err != nil {
        log.Printf("asset-service: Create: error: %v", err)
        return nil, err
    }
    log.Printf("asset-service: Create: success, id=%s", asset.ID)
    return asset, nil
}

func (s *AssetService) Update(ctx context.Context, id, fullname string) error {
    log.Printf("asset-service: Update called: id=%s, fullname=%s", id, fullname)
    err := s.repo.Update(ctx, id, fullname)
    if err != nil {
        log.Printf("asset-service: Update: error: %v", err)
        return err
    }
    log.Printf("asset-service: Update: success, id=%s", id)
    return nil
}

func (s *AssetService) Deactivate(ctx context.Context, id string) error {
    log.Printf("asset-service: Deactivate called: id=%s", id)
    err := s.repo.SetActive(ctx, id, false)
    if err != nil {
        log.Printf("asset-service: Deactivate: error: %v", err)
        return err
    }
    log.Printf("asset-service: Deactivate: success, id=%s", id)
    return nil
}

func (s *AssetService) Activate(ctx context.Context, id string) error {
    log.Printf("asset-service: Activate called: id=%s", id)
    err := s.repo.SetActive(ctx, id, true)
    if err != nil {
        log.Printf("asset-service: Activate: error: %v", err)
        return err
    }
    log.Printf("asset-service: Activate: success, id=%s", id)
    return nil
}

func (s *AssetService) Delete(ctx context.Context, id string) error {
    log.Printf("asset-service: Delete called: id=%s", id)
    err := s.repo.Delete(ctx, id)
    if err != nil {
        log.Printf("asset-service: Delete: error: %v", err)
        return err
    }
    log.Printf("asset-service: Delete: success, id=%s", id)
    return nil
}