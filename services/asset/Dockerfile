FROM golang:1.26.2-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o asset_service ./services/asset

FROM scratch

COPY --from=builder /src/asset_service .

CMD ["./asset_service"]