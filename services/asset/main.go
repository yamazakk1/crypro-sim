package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "crypto-simulator/services/asset/internal/app"
)

func main() {
    ctx := context.Background()

    assetApp, err := app.New(ctx)
    if err != nil {
        log.Fatalf("failed to create app: %v", err)
    }

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-quit
        assetApp.Stop()
    }()

    if err := assetApp.Run(); err != nil {
        log.Fatalf("Failed to run app: %v", err)
    }
}