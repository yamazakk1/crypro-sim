package app

import (
    "context"
    "fmt"
    "log"
    "net"

    "github.com/jackc/pgx/v5/pgxpool"
    "google.golang.org/grpc"

    pb "crypto-simulator/pkg/pb/trading"
    "crypto-simulator/services/trading/internal/config"
    "crypto-simulator/services/trading/internal/handler"
    "crypto-simulator/services/trading/internal/repository"
    "crypto-simulator/services/trading/internal/service"
)

type TradingApp struct {
    cfg        *config.TradingServiceConfig
    grpcServer *grpc.Server
    pool       *pgxpool.Pool
}

func New(ctx context.Context) (*TradingApp, error) {
    cfg := config.NewTradingServiceConfig()

    pool, err := pgxpool.New(ctx, cfg.DB_DSN)
    if err != nil {
        return nil, fmt.Errorf("failed to create DB pool: %w", err)
    }
    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("failed to ping DB: %w", err)
    }
    log.Println("trading-app: connected to database")

    portfolioRepo := repository.NewPortfolioRepo(pool)
    txRepo := repository.NewTransactionRepo(pool)
    tradingSvc := service.NewTradingService(pool, portfolioRepo, txRepo)
    tradingHandler := handler.NewTradingHandler(tradingSvc)

    grpcServer := grpc.NewServer()
    pb.RegisterTradingServiceServer(grpcServer, tradingHandler)

    return &TradingApp{cfg: cfg, grpcServer: grpcServer, pool: pool}, nil
}

func (a *TradingApp) Run() error {
    addr := fmt.Sprintf(":%s", a.cfg.Port)
    lis, _ := net.Listen("tcp", addr)
    log.Printf("trading-app: listening on %s", addr)
    return a.grpcServer.Serve(lis)
}

func (a *TradingApp) Stop() {
    a.grpcServer.GracefulStop()
    a.pool.Close()
    log.Println("trading-app: stopped")
}