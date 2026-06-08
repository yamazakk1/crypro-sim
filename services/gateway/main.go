package main

import (
	"context"
	"crypto-simulator/services/gateway/internal/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, _ := app.NewGatewayApp()
	defer app.Stop(context.Background())

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		app.Stop(context.Background())
	}()

	if err := app.Run(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
