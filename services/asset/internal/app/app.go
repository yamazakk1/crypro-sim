package app

import (
    "context"
    "fmt"
    "log"
    "net"

    "github.com/jackc/pgx/v5/pgxpool"
    "google.golang.org/grpc"

    pb "crypto-simulator/pkg/pb/asset"
    "crypto-simulator/services/asset/internal/config"
    "crypto-simulator/services/asset/internal/handler"
    "crypto-simulator/services/asset/internal/repository"
    "crypto-simulator/services/asset/internal/service"
)

type AssetApp struct {
    cfg        *config.AssetServiceConfig
    grpcServer *grpc.Server
    pool       *pgxpool.Pool
}

func New(ctx context.Context) (*AssetApp, error) {
    log.Println("asset-app: starting initialization")
    cfg := config.NewAssetServiceConfig()
    log.Printf("asset-app: config loaded, port=%s", cfg.Port)

    pool, err := pgxpool.New(ctx, cfg.DB_DSN)
    if err != nil {
        log.Printf("asset-app: failed to create DB pool: %v", err)
        return nil, fmt.Errorf("failed to create DB pool: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        log.Printf("asset-app: failed to ping DB: %v", err)
        return nil, fmt.Errorf("failed to ping DB: %w", err)
    }
    log.Printf("asset-app: connected to database, link: %s", cfg.DB_DSN)

    assetRepo := repository.NewAssetRepo(pool)
    assetService := service.NewAssetService(assetRepo)
    assetHandler := handler.NewAssetHandler(assetService)

    grpcServer := grpc.NewServer()
    pb.RegisterAssetServiceServer(grpcServer, assetHandler)

    log.Println("asset-app: initialization complete")
    return &AssetApp{
        cfg:        cfg,
        grpcServer: grpcServer,
        pool:       pool,
    }, nil
}

func (a *AssetApp) Run() error {
    addr := fmt.Sprintf(":%s", a.cfg.Port)
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        log.Printf("asset-app: failed to listen on %s: %v", addr, err)
        return fmt.Errorf("failed to listen: %w", err)
    }

    log.Printf("asset-app: listening on %s", addr)
    return a.grpcServer.Serve(listener)
}

func (a *AssetApp) Stop() {
    log.Println("asset-app: shutting down")
    a.grpcServer.GracefulStop()
    a.pool.Close()
    log.Println("asset-app: stopped")
}