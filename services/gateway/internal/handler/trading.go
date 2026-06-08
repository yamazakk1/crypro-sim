package handler

import (
	"crypto-simulator/services/gateway/internal/middleware"
	"encoding/json"
	"net/http"

	pbTrading "crypto-simulator/pkg/pb/trading"

	"google.golang.org/grpc/status"
)

func (h *GatewayHandler) HandleGetPortfolio(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	resp, err := h.TradingServiceClient.GetPortfolio(r.Context(), &pbTrading.GetPortfolioRequest{UserId: userID})
	if err != nil {
		st := status.Convert(err)
		writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
		return
	}
	writeJson(w, http.StatusOK, resp)
}

func (h *GatewayHandler) HandleBuy(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	var req struct {
		AssetID  string `json:"asset_id"`
		Quantity string `json:"quantity"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.TradingServiceClient.Buy(r.Context(), &pbTrading.BuyRequest{
		UserId:   userID,
		AssetId:  req.AssetID,
		Quantity: req.Quantity,
	})
	if err != nil {
		st := status.Convert(err)
		writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
		return
	}
	writeJson(w, http.StatusOK, resp)
}

func (h *GatewayHandler) HandleSell(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	var req struct {
		AssetID  string `json:"asset_id"`
		Quantity string `json:"quantity"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.TradingServiceClient.Sell(r.Context(), &pbTrading.SellRequest{
		UserId:   userID,
		AssetId:  req.AssetID,
		Quantity: req.Quantity,
	})
	if err != nil {
		st := status.Convert(err)
		writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
		return
	}
	writeJson(w, http.StatusOK, resp)
}

func (h *GatewayHandler) HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	resp, err := h.TradingServiceClient.GetTransactions(r.Context(), &pbTrading.GetTransactionsRequest{UserId: userID})
	if err != nil {
		st := status.Convert(err)
		writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
		return
	}
	writeJson(w, http.StatusOK, resp)
}

func (h *GatewayHandler) HandleAddBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	var req struct {
		Amount string `json:"amount"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	resp, err := h.TradingServiceClient.AddBalance(r.Context(), &pbTrading.AddBalanceRequest{
		UserId: userID,
		Amount: req.Amount,
	})
	if err != nil {
		st := status.Convert(err)
		writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
		return
	}
	writeJson(w, http.StatusOK, resp)
}
