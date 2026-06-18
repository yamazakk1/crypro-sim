#!/bin/sh
echo "Waiting for DB..."
sleep 5

cd /app
for f in migrations/*.up.sql; do
    psql $DB_DSN -f "$f"
done


/app/seed

/app/auth &
/app/asset &
/app/market &
/app/trading &
/app/ws-hub &
/app/gateway &

wait