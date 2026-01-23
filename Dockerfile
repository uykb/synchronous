# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o crypto-sync-bot ./cmd/main.go

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates
RUN apk add --no-cache ca-certificates

# Copy binary
COPY --from=builder /app/crypto-sync-bot .

# Expose API and Metrics ports
EXPOSE 8080

CMD ["./crypto-sync-bot"]
