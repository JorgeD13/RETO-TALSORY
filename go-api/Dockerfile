# ── Stage 1: build ──────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar archivos de módulos primero para aprovechar caché de capas.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar el binario estático para Linux.
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-api ./cmd/main.go

# ── Stage 2: runtime ─────────────────────────────────────────────────────────
FROM alpine:3.20

WORKDIR /app

# Certificados para llamadas HTTPS externas.
RUN apk --no-cache add ca-certificates

COPY --from=builder /go-api .
COPY --from=builder /app/.env.example .env

EXPOSE 8080

ENTRYPOINT ["/app/go-api"]
