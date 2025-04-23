# Стадия сборки
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o go-ai-service main.go

# Финальный образ
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/go-ai-service .

# Для доступа к .env, если нужен
COPY .env .

ENV GIN_MODE=release

CMD ["./go-ai-service"]