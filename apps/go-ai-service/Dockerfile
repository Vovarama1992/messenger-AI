# ========================
# Build stage
# ========================
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o go-ai-service -trimpath -ldflags="-s -w" main.go

# ========================
# Final stage
# ========================
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/go-ai-service .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/go-ai-service"]