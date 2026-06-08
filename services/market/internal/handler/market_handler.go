package handler

import (
    "context"
    "log"
	"fmt"
    pb "crypto-simulator/pkg/pb/market"
    "crypto-simulator/services/market/internal/models"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type MarketService interface {
    GetCurrentPrices(ctx context.Context) ([]models.PriceUpdate, error)
    GetPriceHistory(ctx context.Context, assetID, from, to string) ([]models.PricePoint, error)
}

type MarketHandler struct {
    pb.UnimplementedMarketServiceServer
    service MarketService
}

func NewMarketHandler(svc MarketService) *MarketHandler {
    return &MarketHandler{service: svc}
}

func (h *MarketHandler) GetCurrentPrices(ctx context.Context, req *pb.GetCurrentPricesRequest) (*pb.GetCurrentPricesResponse, error) {
    prices, err := h.service.GetCurrentPrices(ctx)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    var pbPrices []*pb.PriceUpdate
    for _, p := range prices {
        pbPrices = append(pbPrices, &pb.PriceUpdate{
            AssetId:       p.AssetID,
            Symbol:        p.Symbol,
            PriceUsdt:     fmt.Sprintf("%.6f", p.Price),
            ChangeUsdt:    fmt.Sprintf("%.6f", p.Change),
            ChangePercent: fmt.Sprintf("%.4f", p.ChangePercent),
            Timestamp:     p.Timestamp,
        })
    }
    return &pb.GetCurrentPricesResponse{Prices: pbPrices}, nil
}

func (h *MarketHandler) GetPriceHistory(ctx context.Context, req *pb.GetPriceHistoryRequest) (*pb.GetPriceHistoryResponse, error) {
    log.Printf("market-handler: GetPriceHistory called: asset=%s", req.AssetId)
    points, err := h.service.GetPriceHistory(ctx, req.AssetId, req.From, req.To)
    if err != nil {
        return nil, status.Error(codes.Internal, err.Error())
    }

    var pbPoints []*pb.PricePoint
    for _, p := range points {
        pbPoints = append(pbPoints, &pb.PricePoint{
            Timestamp: p.Timestamp,
            Open:      fmt.Sprintf("%.6f", p.Open),
            High:      fmt.Sprintf("%.6f", p.High),
            Low:       fmt.Sprintf("%.6f", p.Low),
            Close:     fmt.Sprintf("%.6f", p.Close),
        })
    }

    return &pb.GetPriceHistoryResponse{Points: pbPoints}, nil
}