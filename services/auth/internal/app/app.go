package app

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "crypto-simulator/pkg/pb/auth"
	"crypto-simulator/services/auth/internal/config"
	"crypto-simulator/services/auth/internal/handler"
	"crypto-simulator/services/auth/internal/repository"
	"crypto-simulator/services/auth/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type AuthApp struct {
	cfg        *config.AuthServiceConfig
	grpcServer *grpc.Server
	pool       *pgxpool.Pool
}

func New(ctx context.Context) (*AuthApp, error) {
	log.Println("auth-app: starting initialization")
	cfg := config.NewAuthServiceConfig()
	log.Printf("auth-app: config loaded, port=%s", cfg.Port)

	pool, err := pgxpool.New(ctx, cfg.DB_DSN)
	if err != nil {
		log.Printf("auth-app: failed to create DB pool: %v", err)
		return nil, fmt.Errorf("failed to create DB pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		log.Printf("auth-app: failed to ping DB: %v", err)
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}
	log.Println("auth-app: connected to database")

	authRepo := repository.NewAuthRepo(pool)
	authService := service.NewAuthService(authRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	log.Println("auth-app: initialization complete")
	return &AuthApp{
		cfg:        cfg,
		grpcServer: grpcServer,
		pool:       pool,
	}, nil
}

func (a *AuthApp) Run() error {
	addr := fmt.Sprintf(":%s", a.cfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("auth-app: failed to listen on %s: %v", addr, err)
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Printf("auth-app: listening on %s", addr)
	return a.grpcServer.Serve(listener)
}

func (a *AuthApp) Stop() {
	log.Println("auth-app: shutting down")
	a.grpcServer.GracefulStop()
	a.pool.Close()
	log.Println("auth-app: stopped")
}