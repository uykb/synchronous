# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and ssl certificates
RUN apk add --no-cache git ca-certificates

# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and ssl certificates
RUN apk add --no-cache git ca-certificates

# Copy source code first to allow go mod tidy to see imports
COPY . .

# Initialize/Tidy the module to resolve dependencies
# We ignore the existing go.mod versions and let tidy resolve them based on imports
RUN go mod tidy

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
