# Crypto-Sync-Bot

A real-time cryptocurrency trading signal synchronization bot that listens to trading signals and executes orders across multiple exchanges (Binance, OKX, Bybit).

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Crypto-Sync-Bot                          │
├────────────────────────────┬────────────────────────────────────┤
│       Frontend (Bun)       │         Backend (Go)               │
│   Vue 3 + TypeScript       │       Gin / Standard Lib           │
│   - Dashboard UI           │   - Signal Processing              │
│   - Settings Panel         │   - Exchange Adapters              │
│   - Real-time Charts       │   - Risk Management                │
│            ↓               │                ↓                   │
│   Nginx (Production)       │         Direct Execution           │
└────────────────────────────┴────────────────────────────────────┘
```

## Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| Backend | Go 1.21+ | Signal processing, order execution, exchange adapters |
| Frontend | Vue 3 + TypeScript | Dashboard UI, settings management |
| Package Manager | Bun 1.3+ | Fast frontend dependency management |
| Container | Docker + Nginx | Production deployment |
| Config | Viper + .env | Environment-based configuration |

## Prerequisites

- **Go** 1.21 or higher (for backend)
- **Bun** 1.3 or higher (for frontend)
- **Docker** 24+ and **Docker Compose** v2
- **Git**

### Installing Bun (if not installed)

```bash
curl -fsSL https://bun.sh/install | bash
```

## Project Structure

```
crypto-sync-bot/
├── cmd/
│   └── main.go                 # Backend entry point
├── internal/
│   ├── config/                 # Configuration loading
│   ├── models/                 # Data structures
│   ├── exchange/               # Exchange implementations
│   ├── processor/              # Business logic
│   └── risk/                   # Risk management
├── frontend/                   # Vue 3 + TypeScript frontend
│   ├── src/
│   │   ├── views/              # Page components
│   │   ├── stores/             # Pinia state management
│   │   ├── router/             # Vue Router configuration
│   │   └── api/                # API client
│   ├── Dockerfile              # Frontend container
│   └── vite.config.ts          # Vite configuration
├── pkg/                        # Shared utilities
├── docker-compose.yml          # Multi-container orchestration
├── .env.example                # Environment template
└── AGENTS.md                   # Developer guidelines
```

## Environment Configuration

Copy the example environment file and configure your credentials:

```bash
cp .env.example .env
```

### Required Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `BINANCE_API_KEY` | Binance API Key | Yes* |
| `BINANCE_API_SECRET` | Binance API Secret | Yes* |
| `BINANCE_TESTNET` | Use testnet (true/false) | No (default: false) |
| `OKX_API_KEY` | OKX API Key | Yes* |
| `OKX_API_SECRET` | OKX API Secret | Yes* |
| `OKX_API_PASSPHRASE` | OKX API Passphrase | Yes* |
| `BYBIT_API_KEY` | Bybit API Key | Yes* |
| `BYBIT_API_SECRET` | Bybit API Secret | Yes* |
| `SYMBOL` | Trading pair (e.g., BTC-USDT) | No (default: BTC-USDT) |
| `POSITION_RATIO` | Position size ratio | No (default: 1.0) |
| `MAX_POSITION` | Maximum position limit | No (default: 1.0) |
| `STOP_LOSS_RATIO` | Stop loss percentage | No (default: 0.05) |

*At least one exchange configuration is required.

## Development

### Backend (Go)

```bash
# Navigate to project root
cd crypto-sync-bot

# Run directly (requires .env configuration)
go run ./cmd/main.go

# Build binary
go build -o crypto-sync-bot ./cmd/main.go

# Run tests
go test ./...

# Run with race detection
go test -race ./...
```

### Frontend (Bun)

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies (first time only)
bun install

# Start development server
bun dev

# Build for production
bun run build

# Preview production build
bun run preview
```

The frontend development server runs on `http://localhost:49617` by default.

### Full Stack Development

To run both frontend and backend simultaneously:

**Terminal 1 - Backend:**
```bash
go run ./cmd/main.go
```

**Terminal 2 - Frontend:**
```bash
cd frontend && bun dev
```

Access the dashboard at `http://localhost:49617`. The frontend is configured to proxy API requests to `http://localhost:8080`.

## Docker Deployment

### Quick Start

```bash
# Build and start all services
docker-compose up -d --build

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### Services

| Service | Port | Description |
|---------|------|-------------|
| `crypto-sync-bot` | 8080 | Backend API (internal) |
| `frontend` | 80 | Web Dashboard |

### Manual Frontend Build

```bash
# Build frontend only
cd frontend
docker build -t crypto-sync-bot-frontend .

# Run frontend container
docker run -p 80:80 crypto-sync-bot-frontend
```

## Production Deployment

### Using Docker Compose (Recommended)

```bash
# Set production environment variables
export BINANCE_API_KEY="your_production_key"
export BINANCE_API_SECRET="your_production_secret"
# ... set other variables

# Deploy
docker-compose up -d --build
```

### Environment File for Production

Create a `.env.prod` file with production credentials:

```bash
BINANCE_API_KEY=prod_key_here
BINANCE_API_SECRET=prod_secret_here
BINANCE_TESTNET=false
OKX_API_KEY=
OKX_API_SECRET=
OKX_API_PASSPHRASE=
BYBIT_API_KEY=
BYBIT_API_SECRET=
SYMBOL=BTC-USDT
POSITION_RATIO=1.0
MAX_POSITION=1.0
STOP_LOSS_RATIO=0.05
```

Then deploy with:
```bash
docker-compose --env-file .env.prod up -d --build
```

### Direct Binary Deployment

```bash
# Build binary
go build -o crypto-sync-bot ./cmd/main.go

# Run with environment variables
BINANCE_API_KEY=xxx BINANCE_API_SECRET=yyy ./crypto-sync-bot
```

## API Endpoints

When running locally or via Docker, the backend exposes:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/signals` | Get recent signals |
| POST | `/api/signals` | Submit a signal |
| GET | `/api/status` | Bot status |
| POST | `/api/start` | Start the bot |
| POST | `/api/stop` | Stop the bot |

## Supported Exchanges

- **Binance** - Spot trading
- **OKX** - Spot trading
- **Bybit** - Spot trading

## Risk Management

The bot includes built-in risk management features:

- **Position Limits**: Maximum position size enforcement
- **Stop Loss**: Automatic stop loss at configurable ratio
- **Symbol Validation**: Only trades configured trading pairs

## Troubleshooting

### Frontend not connecting to backend

Ensure Vite proxy is configured in `frontend/vite.config.ts`:

```typescript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    }
  }
}
```

### Docker build fails

- Ensure Docker daemon is running
- Check that all environment variables are set
- Verify no port conflicts (80, 8080)

### Bot not executing orders

- Verify API keys have correct permissions (Spot Trading)
- Check that testnet mode is disabled for production
- Review logs: `docker-compose logs crypto-sync-bot`

## License

MIT License
