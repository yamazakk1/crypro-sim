FROM golang:1.26.2-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o trading_service ./services/trading

FROM scratch

COPY --from=builder /src/trading_service .

CMD ["./trading_service"]