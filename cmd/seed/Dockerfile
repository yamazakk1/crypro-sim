FROM golang:1.26.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /seed ./cmd/seed

FROM alpine:latest
WORKDIR /app
COPY --from=builder /seed .
CMD ["./seed"]