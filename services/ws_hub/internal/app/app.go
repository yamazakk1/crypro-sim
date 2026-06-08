package app

import (
    "context"
    "log"
    "net/http"

    "github.com/redis/go-redis/v9"
    "crypto-simulator/services/ws_hub/internal/config"
    "crypto-simulator/services/ws_hub/internal/hub"
)

type WSHubApp struct {
    cfg    *config.WSHubConfig
    server *http.Server
    hub    *hub.Hub
}

func New(ctx context.Context) (*WSHubApp, error) {
    cfg := config.NewWSHubConfig()

    rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
    h := hub.NewHub()
    go h.Run()

    // Подписка на Redis
    go h.SubscribeRedis(ctx, rdb)

    mux := http.NewServeMux()
    mux.HandleFunc("/ws", h.ServeWS)

    server := &http.Server{Addr: ":" + cfg.Port, Handler: mux}

    return &WSHubApp{cfg: cfg, server: server, hub: h}, nil
}

func (a *WSHubApp) Run() error {
    log.Printf("ws-hub: listening on :%s", a.cfg.Port)
    return a.server.ListenAndServe()
}

func (a *WSHubApp) Stop(ctx context.Context) {
    a.server.Shutdown(ctx)
}