# ── Stage 1: Build ──────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependency downloads as a separate layer
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Produce a statically linked binary with no CGO dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /subito-backend-challenge ./cmd/server

# ── Stage 2: Minimal runtime image ──────────────────────────────────────────
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /subito-backend-challenge .

EXPOSE 8080

CMD ["./subito-backend-challenge"]
