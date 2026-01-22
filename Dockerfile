# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and ssl certificates
RUN apk add --no-cache git ca-certificates

# Copy dependency files
COPY go.mod go.sum ./

# Ensure dependencies are tidy and downloaded
# We run tidy to fix any missing checksums that might occur due to environment differences
RUN go mod tidy

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

# The application listens on 8080
EXPOSE 8080

CMD ["./crypto-sync-bot"]
