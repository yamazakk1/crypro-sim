package hub

import (
    "context"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
    "github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
    conn *websocket.Conn
    send chan []byte
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            log.Printf("hub: client connected, total=%d", len(h.clients))
        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()
            log.Printf("hub: client disconnected, total=%d", len(h.clients))
        case msg := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- msg:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("hub: upgrade error: %v", err)
        return
    }

    client := &Client{conn: conn, send: make(chan []byte, 256)}
    h.register <- client

    go h.writePump(client)
    go h.readPump(client)
}

func (h *Hub) writePump(c *Client) {
    defer c.conn.Close()
    for msg := range c.send {
        if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            return
        }
    }
}

func (h *Hub) readPump(c *Client) {
    defer func() { h.unregister <- c }()
    for {
        _, _, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
    }
}

func (h *Hub) SubscribeRedis(ctx context.Context, rdb *redis.Client) {
    pubsub := rdb.Subscribe(ctx, "prices")
    defer pubsub.Close()

    log.Println("hub: subscribed to Redis channel 'prices'")
    for msg := range pubsub.Channel() {
        h.broadcast <- []byte(msg.Payload)
    }
}