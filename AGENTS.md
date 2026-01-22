# Agent Guidelines for Crypto-Sync-Bot

This document provides essential instructions for AI agents and developers working on the `crypto-sync-bot` repository. It covers build processes, testing procedures, and coding standards to ensure consistency and reliability.

## 1. Project Overview

Crypto-Sync-Bot is a real-time cryptocurrency trading signal synchronization bot with:
- **Backend**: Go-based signal processing and order execution
- **Frontend**: Vue 3 + TypeScript managed with Bun
- **Deployment**: Docker Compose with Nginx

## 2. Build and Execution

### Backend (Go)

```bash
# Navigate to project root
cd crypto-sync-bot

# Build the binary
go build -o crypto-sync-bot ./cmd/main.go

# Run directly (requires .env configuration)
go run ./cmd/main.go

# Run tests
go test ./...

# Run with race detection
go test -race ./...
```

### Frontend (Bun)

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
bun install

# Development server
bun dev

# Production build
bun run build
```

### Docker (Full Stack)

```bash
# Build and start all services
docker-compose up -d --build

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 3. Testing

### Run All Tests
Execute all unit tests across the project:
```bash
go test ./...
```

### Run Specific Test
To run a single test function within a specific package:
```bash
go test -v ./internal/processor -run TestSignalProcessor_Process
```

### Run with Race Detection
Since this project relies heavily on concurrency (goroutines/channels), always use race detection during development:
```bash
go test -race ./...
```

## 4. Code Style & Conventions

### Backend (Go)
- **Tooling**: Always run `gofmt` or `goimports` before committing.
- **Indentation**: Use tabs for indentation, not spaces.
- **Files**: Use snake_case for filenames (e.g., `signal_processor.go`, `binance_listener.go`).

### Frontend (Vue 3 + TypeScript + Bun)
- **Tooling**: ESLint and Prettier (configured automatically with Bun)
- **Formatting**: Run `bun run lint` before committing
- **Components**: Use Composition API with `<script setup>` syntax
- **State Management**: Pinia for global state
- **Routing**: Vue Router 4

### Naming Conventions
- **Go**: PascalCase for exported, camelCase for unexported
- **Vue/Pinia**: PascalCase for components and store names, camelCase for state/actions
- **Files**: snake_case for Go, kebab-case for Vue components

## 5. Project Structure

```
crypto-sync-bot/
├── cmd/
│   └── main.go                 # Backend entry point
├── internal/
│   ├── config/                 # Configuration loading (Viper)
│   ├── models/                 # Data structures (TradingSignal, OrderResult)
│   ├── exchange/               # Exchange implementations
│   ├── processor/              # Business logic and orchestration
│   └── risk/                   # Risk management logic
├── frontend/                   # Vue 3 + TypeScript frontend
│   ├── src/
│   │   ├── views/              # Page components (Dashboard, Settings)
│   │   ├── stores/             # Pinia state management
│   │   ├── router/             # Vue Router configuration
│   │   └── api/                # API client (Axios)
│   ├── Dockerfile              # Multi-stage build with Bun + Nginx
│   └── vite.config.ts          # Vite + proxy configuration
├── pkg/                        # Public library code (utils)
├── docker-compose.yml          # Multi-container orchestration
├── .env.example                # Environment template
└── .gitignore                  # Git ignore rules
```

## 6. Imports

### Go
- Group imports: standard library, third-party, internal project
- Use full module path `crypto-sync-bot/...` for internal imports

### TypeScript
- Group imports: Vue/React ecosystem, third-party, internal
- Use absolute imports when configured

## 7. Error Handling

### Go
- **Check Errors**: Never ignore errors using `_`. Always handle or propagate them.
- **Context**: Wrap errors with context: `fmt.Errorf("failed to place order: %w", err)`.
- **Logging**: Use `log` or `logrus` to log errors.

### TypeScript
- Use `try/catch` with proper error logging
- Propagate errors with context: `throw new Error(\`failed to fetch: \${error.message}\`)`
- Use Result/Either patterns where appropriate for async operations

## 8. Concurrency

### Go
- **Mutexes**: Use `sync.Mutex` or `sync.RWMutex` to protect shared state.
- **WaitGroups**: Use `sync.WaitGroup` when spawning multiple goroutines.
- **Channels**: Use channels for signal passing. Ensure channels are closed properly.

### TypeScript
- **Async/Await**: Prefer async/await over raw promises
- **Web Workers**: Use for heavy computations (e.g., chart data processing)
- **AbortController**: Use for cancellable requests

## 9. Dependency Management

### Go
- **Viper**: Configuration is managed via `viper`. Do not hardcode credentials.
- **SDKs**: Use established SDKs (`go-binance`, `goex`, `bybit`).

### Frontend (Bun)
- **Package Manager**: Always use `bun install` (not npm or yarn)
- **Lockfile**: Commit `bun.lock` for deterministic builds
- **Build**: Use `bun run build` for production builds

## 10. Specific Implementation Details

### Exchange Executors
- Implement the `ExchangeExecutor` interface for new exchanges.
- Ensure `PlaceOrder` is thread-safe.
- Normalize `TradingSignal` (e.g., "BUY"/"SELL") to exchange-specific types.

### WebSocket Listeners
- Implement reconnection logic for WebSockets (`for` loop with `select`).
- Handle keep-alive mechanisms for listen keys.

### Risk Management
- Always check `RiskManager.PreOrderCheck` before executing any trade.
- Critical checks (max position, symbol match) must pass.

### Frontend Integration
- Frontend proxies API requests to backend via Vite dev server
- Production: Nginx serves static files, backend runs on separate port
- Use Pinia for signal state management across components

## 11. Linting

### Go
While not enforced by CI yet, use `golangci-lint` locally:
```bash
golangci-lint run
```

### Frontend
Run linting before committing:
```bash
cd frontend
bun run lint
```

## 12. Docker Development Workflow

### Development with Hot Reload

1. **Backend**: Run directly with `go run ./cmd/main.go`
2. **Frontend**: Run with `cd frontend && bun dev`
3. Access dashboard at `http://localhost:49617`

### Production Build

```bash
docker-compose up -d --build
```

The frontend container uses a multi-stage build:
- **Build stage**: `oven/bun:latest` for fast dependency install and build
- **Production stage**: `nginx:stable-alpine` for serving static files

### Frontend-Only Build

```bash
cd frontend
docker build -t crypto-sync-bot-frontend .
docker run -p 80:80 crypto-sync-bot-frontend
```
