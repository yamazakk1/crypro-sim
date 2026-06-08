package handler

import (
	"context"
	"log"

	pb "crypto-simulator/pkg/pb/asset"
	"crypto-simulator/services/asset/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AssetService interface {
    List(ctx context.Context) ([]*models.Asset, error)
    GetByID(ctx context.Context, id string) (*models.Asset, error)
    Create(ctx context.Context, symbol, fullname string, initPrice float64) (*models.Asset, error)
    Update(ctx context.Context, id, fullname string) error
    Deactivate(ctx context.Context, id string) error
    Activate(ctx context.Context, id string) error
    Delete(ctx context.Context, id string) error
}

type AssetHandler struct {
    pb.UnimplementedAssetServiceServer
    service AssetService
}

func NewAssetHandler(svc AssetService) *AssetHandler {
    return &AssetHandler{service: svc}
}

func (h *AssetHandler) ListAssets(ctx context.Context, req *pb.ListAssetsRequest) (*pb.ListAssetsResponse, error) {
    log.Println("asset-handler: ListAssets called")
    assets, err := h.service.List(ctx)
    if err != nil {
        log.Printf("asset-handler: ListAssets: error: %v", err)
        return nil, status.Error(codes.Internal, err.Error())
    }
    var pbAssets []*pb.Asset
    for _, a := range assets {
        pbAssets = append(pbAssets, &pb.Asset{Id: a.ID, Symbol: a.Symbol, Fullname: a.Fullname, IsActive: a.IsActive})
    }
    log.Printf("asset-handler: ListAssets: success, count=%d", len(pbAssets))
    return &pb.ListAssetsResponse{Assets: pbAssets}, nil
}

func (h *AssetHandler) GetAsset(ctx context.Context, req *pb.GetAssetRequest) (*pb.Asset, error) {
    log.Printf("asset-handler: GetAsset called: id=%s", req.Id)
    a, err := h.service.GetByID(ctx, req.Id)
    if err != nil {
        log.Printf("asset-handler: GetAsset: not found: id=%s", req.Id)
        return nil, status.Error(codes.NotFound, "asset not found")
    }
    log.Printf("asset-handler: GetAsset: success, symbol=%s", a.Symbol)
    return &pb.Asset{Id: a.ID, Symbol: a.Symbol, Fullname: a.Fullname, IsActive: a.IsActive}, nil
}

func (h *AssetHandler) CreateAsset(ctx context.Context, req *pb.CreateAssetRequest) (*pb.Asset, error) {
    log.Printf("asset-handler: CreateAsset called: symbol=%s, fullname=%s", req.Symbol, req.Fullname)
    a, err := h.service.Create(ctx, req.Symbol, req.Fullname, float64(req.InitialPrice))
    if err != nil {
        log.Printf("asset-handler: CreateAsset: error: %v", err)
        return nil, status.Error(codes.InvalidArgument, err.Error())
    }
    log.Printf("asset-handler: CreateAsset: success, id=%s", a.ID)
    return &pb.Asset{Id: a.ID, Symbol: a.Symbol, Fullname: a.Fullname, IsActive: a.IsActive}, nil
}

func (h *AssetHandler) UpdateAsset(ctx context.Context, req *pb.UpdateAssetRequest) (*pb.Asset, error) {
    log.Printf("asset-handler: UpdateAsset called: id=%s, fullname=%s", req.Id, req.Fullname)
    err := h.service.Update(ctx, req.Id, req.Fullname)
    if err != nil {
        log.Printf("asset-handler: UpdateAsset: error: %v", err)
        return nil, status.Error(codes.Internal, err.Error())
    }
    log.Printf("asset-handler: UpdateAsset: success, id=%s", req.Id)
    return &pb.Asset{Id: req.Id, Fullname: req.Fullname}, nil
}

func (h *AssetHandler) DeactivateAsset(ctx context.Context, req *pb.DeactivateAssetRequest) (*pb.Asset, error) {
    log.Printf("asset-handler: DeactivateAsset called: id=%s", req.Id)
    err := h.service.Deactivate(ctx, req.Id)
    if err != nil {
        log.Printf("asset-handler: DeactivateAsset: error: %v", err)
        return nil, status.Error(codes.Internal, err.Error())
    }
    log.Printf("asset-handler: DeactivateAsset: success, id=%s", req.Id)
    return &pb.Asset{Id: req.Id, IsActive: false}, nil
}

func (h *AssetHandler) ActivateAsset(ctx context.Context, req *pb.ActivateAssetRequest) (*pb.Asset, error) {
    log.Printf("asset-handler: ActivateAsset called: id=%s", req.Id)
    err := h.service.Activate(ctx, req.Id)
    if err != nil {
        log.Printf("asset-handler: ActivateAsset: error: %v", err)
        return nil, status.Error(codes.Internal, err.Error())
    }
    log.Printf("asset-handler: ActivateAsset: success, id=%s", req.Id)
    return &pb.Asset{Id: req.Id, IsActive: true}, nil
}

func (h *AssetHandler) DeleteAsset(ctx context.Context, req *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
    log.Printf("asset-handler: DeleteAsset called: id=%s", req.Id)
    err := h.service.Delete(ctx, req.Id)
    if err != nil {
        log.Printf("asset-handler: DeleteAsset: error: %v", err)
        return nil, status.Error(codes.Internal, err.Error())
    }
    log.Printf("asset-handler: DeleteAsset: success, id=%s", req.Id)
    return &pb.DeleteAssetResponse{Success: true}, nil
}