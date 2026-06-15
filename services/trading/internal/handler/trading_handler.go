package handler

import (
	"context"
	"fmt"

	pb "crypto-simulator/pkg/pb/trading"
	"crypto-simulator/services/trading/internal/models"
	"crypto-simulator/services/trading/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TradingHandler struct {
	pb.UnimplementedTradingServiceServer
	svc *service.TradingService
}

func NewTradingHandler(svc *service.TradingService) *TradingHandler {
	return &TradingHandler{svc: svc}
}

func (h *TradingHandler) GetPortfolio(ctx context.Context, req *pb.GetPortfolioRequest) (*pb.PortfolioResponse, error) {
	p, err := h.svc.GetPortfolio(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var items []*pb.PortfolioItem
	for _, item := range p.Items {
		items = append(items, &pb.PortfolioItem{
			AssetId:           item.AssetID,
			Symbol:            item.Symbol,
			Quantity:          fmt.Sprintf("%.8f", item.Quantity),
			AvgBuyPrice:       fmt.Sprintf("%.6f", item.AvgBuyPrice),
			CurrentPrice:      fmt.Sprintf("%.6f", item.CurrentPrice),
			TotalValue:        fmt.Sprintf("%.2f", item.TotalValue),
			ProfitLoss:        fmt.Sprintf("%.2f", item.ProfitLoss),
			ProfitLossPercent: fmt.Sprintf("%.2f", item.ProfitLossPct),
		})
	}

	return &pb.PortfolioResponse{
		Items:           items,
		TotalValueUsdt:  fmt.Sprintf("%.2f", p.TotalValue),
		BalanceUsdt:     fmt.Sprintf("%.2f", p.BalanceUSDT),
		TotalProfitLoss: fmt.Sprintf("%.2f", p.TotalProfitLoss),
	}, nil
}

func (h *TradingHandler) Buy(ctx context.Context, req *pb.BuyRequest) (*pb.TransactionResponse, error) {
	qty := parseFloat(req.Quantity)
	t, err := h.svc.Buy(ctx, req.UserId, req.AssetId, qty)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toTxResponse(t), nil
}

func (h *TradingHandler) Sell(ctx context.Context, req *pb.SellRequest) (*pb.TransactionResponse, error) {
	qty := parseFloat(req.Quantity)
	t, err := h.svc.Sell(ctx, req.UserId, req.AssetId, qty)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return toTxResponse(t), nil
}

func (h *TradingHandler) GetTransactions(ctx context.Context, req *pb.GetTransactionsRequest) (*pb.TransactionsResponse, error) {
	txs, err := h.svc.GetTransactions(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var items []*pb.TransactionItem
	for _, t := range txs {
		items = append(items, &pb.TransactionItem{
			Id:          t.ID,
			Type:        t.Type,
			AssetSymbol: t.AssetSymbol,
			Quantity:    fmt.Sprintf("%.8f", t.Quantity),
			AssetPrice:  fmt.Sprintf("%.6f", t.AssetPrice),
			TotalUsdt:   fmt.Sprintf("%.2f", t.TotalUSDT),
			CreatedAt:   t.CreatedAt,
		})
	}
	return &pb.TransactionsResponse{Transactions: items}, nil
}

func (h *TradingHandler) AddBalance(ctx context.Context, req *pb.AddBalanceRequest) (*pb.BalanceResponse, error) {
	amount := parseFloat(req.Amount)
	balance, err := h.svc.AddBalance(ctx, req.UserId, amount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.BalanceResponse{
		UserId:      req.UserId,
		BalanceUsdt: fmt.Sprintf("%.2f", balance),
	}, nil
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func toTxResponse(t *models.Transaction) *pb.TransactionResponse {
	return &pb.TransactionResponse{
		Id:          t.ID,
		Type:        t.Type,
		AssetSymbol: t.AssetSymbol,
		Quantity:    fmt.Sprintf("%.8f", t.Quantity),
		AssetPrice:  fmt.Sprintf("%.6f", t.AssetPrice),
		TotalUsdt:   fmt.Sprintf("%.2f", t.TotalUSDT),
		CreatedAt:   t.CreatedAt,
	}
}
