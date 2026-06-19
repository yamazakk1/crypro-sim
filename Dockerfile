FROM golang:1.26.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .


RUN CGO_ENABLED=0 go build -o /gateway ./services/gateway
RUN CGO_ENABLED=0 go build -o /auth ./services/auth
RUN CGO_ENABLED=0 go build -o /asset ./services/asset
RUN CGO_ENABLED=0 go build -o /market ./services/market
RUN CGO_ENABLED=0 go build -o /trading ./services/trading
RUN CGO_ENABLED=0 go build -o /ws-hub ./services/ws_hub
RUN CGO_ENABLED=0 go build -o /seed ./cmd/seed

FROM alpine:latest
RUN apk add --no-cache postgresql-client
COPY --from=builder /gateway /auth /asset /market /trading /ws_hub /seed /app/
COPY services/gateway/static /app/static
COPY migrations /app/migrations
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh
CMD ["/app/start.sh"]