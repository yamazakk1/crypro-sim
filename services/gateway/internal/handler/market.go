package handler

import (
	"log"
	"net/http"
	pbMarket "crypto-simulator/pkg/pb/market"
	"google.golang.org/grpc/status"
)

func (h *GatewayHandler) HandleGetCurrentPrices(w http.ResponseWriter, r *http.Request) {
    resp, err := h.MarketServiceClient.GetCurrentPrices(r.Context(), &pbMarket.GetCurrentPricesRequest{})
    if err != nil {
        st := status.Convert(err)
        writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
        return
    }

    var prices []PriceUpdate
    for _, p := range resp.GetPrices() {
        prices = append(prices, PriceUpdate{
            AssetID:       p.GetAssetId(),
            Symbol:        p.GetSymbol(),
            PriceUsdt:     p.GetPriceUsdt(),
            ChangeUsdt:    p.GetChangeUsdt(),
            ChangePercent: p.GetChangePercent(),
            Timestamp:     p.GetTimestamp(),
        })
    }
    writeJson(w, http.StatusOK, GetCurrentPricesResponse{Prices: prices})
}

func (h *GatewayHandler) HandleGetPriceHistory(w http.ResponseWriter, r *http.Request) {
	assetID := r.PathValue("id")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	log.Printf("gateway: HandleGetPriceHistory called: asset=%s", assetID)
	resp, err := h.MarketServiceClient.GetPriceHistory(r.Context(), &pbMarket.GetPriceHistoryRequest{
		AssetId: assetID,
		From:    from,
		To:      to,
	})
	if err != nil {
		st := status.Convert(err)
		writeJson(w, grpcToHTTP(st.Code()), ErrorResponse{Error: st.Message()})
		return
	}

	writeJson(w, http.StatusOK, resp)
}
