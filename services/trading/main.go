package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "crypto-simulator/services/trading/internal/app"
)

func main() {
    ctx := context.Background()
    tradingApp, err := app.New(ctx)
    if err != nil {
        log.Fatalf("failed to create app: %v", err)
    }

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    go func() { <-quit; tradingApp.Stop() }()

    if err := tradingApp.Run(); err != nil {
        log.Fatalf("failed to run app: %v", err)
    }
}