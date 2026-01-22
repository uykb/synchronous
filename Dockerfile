# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and ssl certificates
RUN apk add --no-cache git ca-certificates

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o crypto-sync-bot ./cmd/main.go

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/crypto-sync-bot .

# Copy config if exists (optional if using ENV)
# COPY config.yaml . 

CMD ["./crypto-sync-bot"]
