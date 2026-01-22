# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and certificates
RUN apk add --no-cache git ca-certificates

# Copy only go.mod initially
COPY go.mod ./

# Generate a fresh go.sum and download dependencies
# We don't copy the old go.sum to avoid checksum conflicts
RUN go mod tidy

# Copy the rest of the source code
COPY . .

# Final tidy and build
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o crypto-sync-bot ./cmd/main.go

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
