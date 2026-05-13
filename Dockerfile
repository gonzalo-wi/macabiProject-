# ── Stage 1: build ───────────────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Descargar dependencias primero (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código y compilar
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/

# ── Stage 2: imagen final mínima ─────────────────────────────────────────────
FROM alpine:3.21

# Certificados TLS (necesarios para sslmode=require con la DB)
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8081

CMD ["./server"]
