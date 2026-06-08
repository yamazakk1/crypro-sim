package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    dsn := os.Getenv("DB_DSN")
    if dsn == "" {
        dsn = fmt.Sprintf(
            "postgres://%s:%s@%s:%s/%s?sslmode=disable",
            getEnv("POSTGRES_USER", "myuser"),
            getEnv("POSTGRES_PASSWORD", "vova123"),
            getEnv("DB_SERVICE_NAME", "localhost"),
            getEnv("POSTGRES_PORT", "5432"),
            getEnv("POSTGRES_DB", "maindb"),
        )
    }

    pool, err := pgxpool.New(context.Background(), dsn)
    if err != nil {
        log.Fatalf("Failed to connect to DB: %v", err)
    }
    defer pool.Close()

    if err := pool.Ping(context.Background()); err != nil {
        log.Fatalf("Failed to ping DB: %v", err)
    }

    ctx := context.Background()

    // Хешируем пароли
    adminHash, _ := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
    userHash, _ := bcrypt.GenerateFromPassword([]byte("User1234!"), bcrypt.DefaultCost)

    // Создаём админа
    _, err = pool.Exec(ctx,
        `INSERT INTO users (id, username, email, password_hash, role, balance_usdt) 
         VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO NOTHING`,
        "00000000-0000-0000-0000-000000000001",
        "admin",
        "admin@crypto.local",
        string(adminHash),
        "admin",
        1000000.00,
    )
    if err != nil {
        log.Printf("Admin: %v", err)
    } else {
        log.Println("Admin created (admin@crypto.local / Admin123!)")
    }

    // Создаём тестового пользователя
    _, err = pool.Exec(ctx,
        `INSERT INTO users (id, username, email, password_hash, role, balance_usdt) 
         VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO NOTHING`,
        "00000000-0000-0000-0000-000000000002",
        "testuser",
        "user@crypto.local",
        string(userHash),
        "user",
        10000.00,
    )
    if err != nil {
        log.Printf("Test user: %v", err)
    } else {
        log.Println("Test user created (user@crypto.local / User1234!)")
    }

    // Создаём валюты
    assets := []struct {
        id       string
        symbol   string
        fullname string
        initPrice float64
    }{
        {"10000000-0000-0000-0000-000000000001", "BCT", "ByteCoin", 50000},
        {"10000000-0000-0000-0000-000000000002", "ETH", "Etherium", 3000},
        {"10000000-0000-0000-0000-000000000003", "DGE", "DogeCoin", 1},
        {"10000000-0000-0000-0000-000000000004", "SHB", "ShibaToken", 0.01},
    }

    for _, a := range assets {
        _, err = pool.Exec(ctx,
            `INSERT INTO assets (id, symbol, full_name, initial_price, is_active) 
             VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING`,
            a.id, a.symbol, a.fullname, a.initPrice, true,
        )
        if err != nil {
            log.Printf("Asset %s: %v", a.symbol, err)
        } else {
            log.Printf("Asset created: %s (%s)", a.symbol, a.fullname)
        }
    }

    log.Println("✅ Seed completed: 2 users, 4 assets")
}

func getEnv(key, defaultValue string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultValue
}