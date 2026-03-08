# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Gera os docs do swagger antes de compilar
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/api/main.go -o docs

RUN CGO_ENABLED=0 GOOS=linux go build -o /workshop-api ./cmd/api

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=America/Sao_Paulo

WORKDIR /app

COPY --from=builder /workshop-api .

EXPOSE 8080

ENTRYPOINT ["/app/workshop-api"]
