#!/bin/sh
echo "Waiting for DB..."
sleep 5

cd /app
for f in migrations/*.up.sql; do
    psql $DB_DSN -f "$f"
done


/app/seed

export REDIS_ADDR=$(echo "$REDIS_ADDR" | sed 's|redis://||' | sed 's|.*@||')
echo "Redis addr: $REDIS_ADDR"

export AUTH_ADDR=localhost:8081
export ASSET_ADDR=localhost:8082
export MARKET_ADDR=localhost:8083
export TRADING_ADDR=localhost:8084

PORT=8081 /app/auth &
PORT=8082 /app/asset &
PORT=8083 /app/market &
PORT=8084 /app/trading &
PORT=8085 /app/ws-hub &
PORT=8080 /app/gateway &

wait