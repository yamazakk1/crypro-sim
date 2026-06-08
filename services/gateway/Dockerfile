FROM golang:1.26.2-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./services/gateway

FROM scratch

COPY --from=builder /src/gateway .
COPY --from=builder /src/services/gateway/static ./static

CMD ["./gateway"]