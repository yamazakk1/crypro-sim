FROM golang:1.26.2-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o market_service ./services/market

FROM scratch

COPY --from=builder /src/market_service .

CMD ["./market_service"]