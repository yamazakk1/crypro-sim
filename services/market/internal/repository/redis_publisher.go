package repository

import (
    "context"
    "encoding/json"
    "log"

    "github.com/redis/go-redis/v9"
)

type RedisPublisher struct {
    client *redis.Client
}

func NewRedisPublisher(addr string) *RedisPublisher {
    client := redis.NewClient(&redis.Options{Addr: addr})
    return &RedisPublisher{client: client}
}

func (p *RedisPublisher) PublishPrice(ctx context.Context, assetID, symbol string, price, change, changePercent float64) error {
    msg, _ := json.Marshal(map[string]interface{}{
        "asset_id":       assetID,
        "symbol":         symbol,
        "price_usdt":     price,
        "change_usdt":    change,
        "change_percent": changePercent,
    })
    err := p.client.Publish(ctx, "prices", msg).Err()
    if err != nil {
        log.Printf("redis: publish error: %v", err)
    }
    return err
}