FROM golang:1.26.2-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ws_hub ./services/ws_hub

FROM scratch

COPY --from=builder /src/ws_hub .

CMD ["./ws_hub"]