package app

import (
    "context"
    "fmt"
    "log"
    "net"

    "github.com/jackc/pgx/v5/pgxpool"
    "google.golang.org/grpc"

    pb "crypto-simulator/pkg/pb/market"
    "crypto-simulator/services/market/internal/config"
    "crypto-simulator/services/market/internal/handler"
    "crypto-simulator/services/market/internal/repository"
    "crypto-simulator/services/market/internal/service"
)

type MarketApp struct {
    cfg        *config.MarketServiceConfig
    grpcServer *grpc.Server
    pool       *pgxpool.Pool
}

func New(ctx context.Context) (*MarketApp, error) {
    cfg := config.NewMarketServiceConfig()

    pool, err := pgxpool.New(ctx, cfg.DB_DSN)
    if err != nil {
        return nil, fmt.Errorf("failed to create DB pool: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("failed to ping DB: %w", err)
    }
    log.Println("market-app: connected to database")

    // Репозитории
    assetProvider := repository.NewAssetProvider(pool)
    priceSaver := repository.NewPriceSaver(pool)
    priceRepo := repository.NewPriceRepo(pool)
	redisPublisher := repository.NewRedisPublisher(cfg.RedisAddr)

    // Генератор цен
    gen := service.NewPriceGenerator(assetProvider, priceSaver, redisPublisher, 8)
    go gen.Start(ctx)
    log.Println("market-app: price generator started")

    // Сервис и хендлер
    marketSvc := service.NewMarketService(priceRepo)
    marketHandler := handler.NewMarketHandler(marketSvc)

    grpcServer := grpc.NewServer()
    pb.RegisterMarketServiceServer(grpcServer, marketHandler)

    return &MarketApp{
        cfg:        cfg,
        grpcServer: grpcServer,
        pool:       pool,
    }, nil
}

func (a *MarketApp) Run() error {
    addr := fmt.Sprintf(":%s", a.cfg.Port)
    lis, _ := net.Listen("tcp", addr)
    log.Printf("market-app: listening on %s", addr)
    return a.grpcServer.Serve(lis)
}

func (a *MarketApp) Stop() {
    a.grpcServer.GracefulStop()
    a.pool.Close()
    log.Println("market-app: stopped")
}