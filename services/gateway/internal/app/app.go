package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"net/http/httputil"
    "net/url"

	pbAsset "crypto-simulator/pkg/pb/asset"
	pbAuth "crypto-simulator/pkg/pb/auth"
	pbMarket "crypto-simulator/pkg/pb/market"
	pbTrading "crypto-simulator/pkg/pb/trading"
	"crypto-simulator/services/gateway/internal/config"
	"crypto-simulator/services/gateway/internal/handler"
	"crypto-simulator/services/gateway/internal/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GatewayApp struct {
	cfg         *config.GatewayServiceConfig
	httpServer  *http.Server
	authConn    *grpc.ClientConn
	assetConn   *grpc.ClientConn
	marketConn  *grpc.ClientConn
	tradingConn *grpc.ClientConn
}

func NewGatewayApp() (*GatewayApp, error) {
	log.Println("gateway: starting initialization")
	cfg := config.NewGatewayServiceConfig()
	log.Printf("gateway: config loaded, port=%s, auth_addr=%s, asset_addr=%s", cfg.Port, cfg.AuthServiceAddr, cfg.AssetServiceAddr)

	authConn, err := grpc.NewClient(cfg.AuthServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("gateway: failed to connect to auth service: %v", err)
		return nil, fmt.Errorf("failed to connect to auth: %w", err)
	}
	log.Printf("gateway: connected to auth service at %s", cfg.AuthServiceAddr)

	assetConn, err := grpc.NewClient(cfg.AssetServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		authConn.Close()
		log.Printf("gateway: failed to connect to asset service: %v", err)
		return nil, fmt.Errorf("failed to connect to asset: %w", err)
	}
	log.Printf("gateway: connected to asset service at %s", cfg.AssetServiceAddr)

	marketConn, err := grpc.NewClient(cfg.MarketServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		marketConn.Close()
		log.Printf("gateway: failed to connect to market service: %v", err)
		return nil, fmt.Errorf("failed to connect to market: %w", err)
	}

	tradingConn, err := grpc.NewClient(cfg.TradingServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		tradingConn.Close()
		log.Printf("gateway: failed to connect to trading service: %v", err)
		return nil, fmt.Errorf("failed to connect to trading: %w", err)
	}
	gatewayHandler := handler.NewGatewayHandler(
		pbAuth.NewAuthServiceClient(authConn),
		pbAsset.NewAssetServiceClient(assetConn),
		pbMarket.NewMarketServiceClient(marketConn),
		pbTrading.NewTradingServiceClient(tradingConn),
	)

	authMiddleware := middleware.JWTAuth(cfg.JWTSecret)
	adminMiddleware := middleware.RequireAdmin

	mux := http.NewServeMux()

	// ──────────────────────────────────────────────
	// ПУБЛИЧНЫЕ СТРАНИЦЫ (без middleware)
	// ──────────────────────────────────────────────
	mux.HandleFunc("/register.html", serveStatic("./static/register.html"))
	mux.HandleFunc("/login.html", serveStatic("./static/login.html"))
	mux.HandleFunc("/assets.html", serveStatic("./static/assets.html"))
	mux.HandleFunc("/admin.html", serveStatic("./static/admin.html"))
	mux.HandleFunc("/portfolio.html", serveStatic("./static/portfolio.html"))

	// ──────────────────────────────────────────────
	// AUTH API
	// ──────────────────────────────────────────────
	// Публичные
	mux.HandleFunc("POST /register", gatewayHandler.HandleRegister)
	mux.HandleFunc("POST /login", gatewayHandler.HandleLogin)
	// Защищённые
	mux.Handle("GET /getme", authMiddleware(http.HandlerFunc(gatewayHandler.GetMe)))

	// ──────────────────────────────────────────────
	// ASSET API
	// ──────────────────────────────────────────────
	// Публичные
	mux.HandleFunc("GET /api/assets", gatewayHandler.HandleListAssets)
	mux.HandleFunc("GET /api/assets/{id}", gatewayHandler.HandleGetAsset)
	// Админские
	mux.Handle("POST /api/assets", authMiddleware(adminMiddleware(http.HandlerFunc(gatewayHandler.HandleCreateAsset))))
	mux.Handle("PUT /api/assets/{id}", authMiddleware(adminMiddleware(http.HandlerFunc(gatewayHandler.HandleUpdateAsset))))
	mux.Handle("POST /api/assets/{id}/deactivate", authMiddleware(adminMiddleware(http.HandlerFunc(gatewayHandler.HandleDeactivateAsset))))
	mux.Handle("POST /api/assets/{id}/activate", authMiddleware(adminMiddleware(http.HandlerFunc(gatewayHandler.HandleActivateAsset))))
	mux.Handle("DELETE /api/assets/{id}", authMiddleware(adminMiddleware(http.HandlerFunc(gatewayHandler.HandleDeleteAsset))))

	// ──────────────────────────────────────────────
	// MARKET API
	// ──────────────────────────────────────────────
	// Публичные
	mux.HandleFunc("GET /api/market/prices", gatewayHandler.HandleGetCurrentPrices)
	mux.HandleFunc("GET /api/market/prices/{id}", gatewayHandler.HandleGetPriceHistory)
	server := &http.Server{
		Addr:    cfg.Port,
		Handler: mux,
	}
	wsURL, _ := url.Parse("http://localhost:8085")
	mux.Handle("/ws", httputil.NewSingleHostReverseProxy(wsURL))
	
	mux.Handle("GET /api/trading/portfolio", authMiddleware(http.HandlerFunc(gatewayHandler.HandleGetPortfolio)))
	mux.Handle("POST /api/trading/buy", authMiddleware(http.HandlerFunc(gatewayHandler.HandleBuy)))
	mux.Handle("POST /api/trading/sell", authMiddleware(http.HandlerFunc(gatewayHandler.HandleSell)))
	mux.Handle("GET /api/trading/transactions", authMiddleware(http.HandlerFunc(gatewayHandler.HandleGetTransactions)))
	mux.Handle("POST /api/trading/balance", authMiddleware(http.HandlerFunc(gatewayHandler.HandleAddBalance)))

	log.Println("gateway: initialization complete")
	return &GatewayApp{
		cfg:        cfg,
		httpServer: server,
		authConn:   authConn,
		assetConn:  assetConn,
	}, nil
}

func (a *GatewayApp) Run() error {
	log.Printf("gateway: listening on %s", a.cfg.Port)
	return a.httpServer.ListenAndServe()
}

func (a *GatewayApp) Stop(ctx context.Context) {
	log.Println("gateway: shutting down")
	a.httpServer.Shutdown(ctx)
	a.authConn.Close()
	a.assetConn.Close()
	log.Println("gateway: stopped")
}

func serveStatic(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
