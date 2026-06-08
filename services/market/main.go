package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "crypto-simulator/services/market/internal/app"
)

func main() {
    ctx := context.Background()

    marketApp, err := app.New(ctx)
    if err != nil {
        log.Fatalf("failed to create app: %v", err)
    }

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-quit
        marketApp.Stop()
    }()

    if err := marketApp.Run(); err != nil {
        log.Fatalf("failed to run app: %v", err)
    }
}