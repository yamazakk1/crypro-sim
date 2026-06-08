package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"crypto-simulator/services/auth/internal/app"
)

func main() {
	ctx := context.Background()
	authApp, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to create app: %w", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		authApp.Stop()
	}()

	if err := authApp.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
