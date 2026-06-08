package handler

import (
	"context"
	"log"
	"strconv"

	pb "crypto-simulator/pkg/pb/trading"
	"crypto-simulator/services/trading/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TradingHandler struct {
	pb.UnimplementedTradingServiceServer
	service *service.TradingService
}

func NewTradingHandler(svc *service.TradingService) *TradingHandler {
	return &TradingHandler{service: svc}
}

func (h *TradingHandler) GetPortfolio(ctx context.Context, req *pb.GetPortfolioRequest) (*pb.PortfolioResponse, error) {
	log.Printf("trading-handler: GetPortfolio: user=%s", req.UserId)
	p, err := h.service.GetPortfolio(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var items []*pb.PortfolioItem
	for _, item := range p.Items {
		items = append(items, &pb.PortfolioItem{
			AssetId:           item.AssetID,
			Symbol:            item.Symbol,
			Quantity:          strconv.FormatFloat(item.Quantity, 'f', 8, 64),
			AvgBuyPrice:       strconv.FormatFloat(item.AvgBuyPrice, 'f', 6, 64),
			CurrentPrice:      strconv.FormatFloat(item.CurrentPrice, 'f', 6, 64),
			TotalValue:        strconv.FormatFloat(item.TotalValue, 'f', 2, 64),
			ProfitLoss:        strconv.FormatFloat(item.ProfitLoss, 'f', 2, 64),
			ProfitLossPercent: strconv.FormatFloat(item.ProfitLossPct, 'f', 2, 64),
		})
	}

	return &pb.PortfolioResponse{
		Items:           items,
		TotalValueUsdt:  strconv.FormatFloat(p.TotalValue, 'f', 2, 64),
		BalanceUsdt:     strconv.FormatFloat(p.BalanceUSDT, 'f', 2, 64),
		TotalProfitLoss: strconv.FormatFloat(p.TotalProfitLoss, 'f', 2, 64),
	}, nil
}

func (h *TradingHandler) Buy(ctx context.Context, req *pb.BuyRequest) (*pb.TransactionResponse, error) {
	qty, _ := strconv.ParseFloat(req.Quantity, 64)
	log.Printf("trading-handler: Buy: user=%s asset=%s qty=%.4f", req.UserId, req.AssetId, qty)

	t, err := h.service.Buy(ctx, req.UserId, req.AssetId, qty)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TransactionResponse{
		Id:          t.ID,
		Type:        t.Type,
		AssetSymbol: t.AssetSymbol,
		Quantity:    strconv.FormatFloat(t.Quantity, 'f', 8, 64),
		AssetPrice:  strconv.FormatFloat(t.AssetPrice, 'f', 6, 64),
		TotalUsdt:   strconv.FormatFloat(t.TotalUSDT, 'f', 2, 64),
		CreatedAt:   t.CreatedAt,
	}, nil
}

func (h *TradingHandler) Sell(ctx context.Context, req *pb.SellRequest) (*pb.TransactionResponse, error) {
	qty, _ := strconv.ParseFloat(req.Quantity, 64)
	log.Printf("trading-handler: Sell: user=%s asset=%s qty=%.4f", req.UserId, req.AssetId, qty)

	t, err := h.service.Sell(ctx, req.UserId, req.AssetId, qty)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.TransactionResponse{
		Id:          t.ID,
		Type:        t.Type,
		AssetSymbol: t.AssetSymbol,
		Quantity:    strconv.FormatFloat(t.Quantity, 'f', 8, 64),
		AssetPrice:  strconv.FormatFloat(t.AssetPrice, 'f', 6, 64),
		TotalUsdt:   strconv.FormatFloat(t.TotalUSDT, 'f', 2, 64),
		CreatedAt:   t.CreatedAt,
	}, nil
}

func (h *TradingHandler) GetTransactions(ctx context.Context, req *pb.GetTransactionsRequest) (*pb.TransactionsResponse, error) {
	log.Printf("trading-handler: GetTransactions: user=%s", req.UserId)
	txs, err := h.service.GetTransactions(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var items []*pb.TransactionItem
	for _, t := range txs {
		items = append(items, &pb.TransactionItem{
			Id:          t.ID,
			Type:        t.Type,
			AssetSymbol: t.AssetSymbol,
			Quantity:    strconv.FormatFloat(t.Quantity, 'f', 8, 64),
			AssetPrice:  strconv.FormatFloat(t.AssetPrice, 'f', 6, 64),
			TotalUsdt:   strconv.FormatFloat(t.TotalUSDT, 'f', 2, 64),
			CreatedAt:   t.CreatedAt,
		})
	}

	return &pb.TransactionsResponse{Transactions: items}, nil
}

func (h *TradingHandler) AddBalance(ctx context.Context, req *pb.AddBalanceRequest) (*pb.BalanceResponse, error) {
	amount, _ := strconv.ParseFloat(req.Amount, 64)
	log.Printf("trading-handler: AddBalance: user=%s amount=%.2f", req.UserId, amount)

	balance, err := h.service.AddBalance(ctx, req.UserId, amount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.BalanceResponse{
		UserId:      req.UserId,
		BalanceUsdt: strconv.FormatFloat(balance, 'f', 2, 64),
	}, nil
}
