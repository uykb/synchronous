# Crypto-Sync-Bot

A real-time cryptocurrency trading signal synchronization bot that listens to trading signals and executes orders across multiple exchanges (Binance, OKX, Bybit).

## Architecture

The project uses a decoupled architecture for optimal performance and deployment flexibility:

```
┌─────────────────────────────────────────────────────────────────┐
│                        Crypto-Sync-Bot                          │
├────────────────────────────┬────────────────────────────────────┤
│      Frontend (Vercel)     │         Backend (Docker)           │
│   Vue 3 + TypeScript       │       Go 1.24 + Redis              │
│   - Dashboard UI           │   - Signal Processing              │
│   - Settings Panel         │   - Exchange Adapters              │
│   - Real-time Charts       │   - Risk Management                │
│            ↓               │                ↓                   │
│      Vercel Edge           │        Docker Container            │
│   (Proxy /api -> Backend)  │      (Exposes Port 8080)           │
└────────────────────────────┴────────────────────────────────────┘
```

## Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| **Backend** | Go 1.24 | High-performance signal processing & execution |
| **Frontend** | Vue 3 + TypeScript | Reactive dashboard & configuration UI |
| **Database** | Redis + SQLite | Message queue (Streams) & persistence |
| **Deployment** | Docker + Vercel | Containerized backend, Serverless frontend |

## Prerequisites

- **Docker** & **Docker Compose** (for backend)
- **Node.js** / **Bun** (for local frontend development)
- **Vercel Account** (for frontend deployment)

## 🚀 Deployment Guide

### 1. Backend Deployment (Docker)

The backend is containerized and available via GitHub Container Registry.

**Option A: Using Docker Compose (Recommended)**

1. Create a `docker-compose.yml` file on your server:

```yaml
version: '3.8'

services:
  crypto-sync-bot:
    image: ghcr.io/uykb/synchronous-backend:latest
    container_name: crypto-sync-bot
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - BINANCE_API_KEY=${BINANCE_API_KEY}
      - BINANCE_API_SECRET=${BINANCE_API_SECRET}
      - BINANCE_TESTNET=false
      - BYBIT_API_KEY=${BYBIT_API_KEY}
      - BYBIT_API_SECRET=${BYBIT_API_SECRET}
      - SYMBOL=BTC-USDT
      - POSITION_RATIO=1.0
    volumes:
      - ./data:/app/data
    depends_on:
      - redis

  redis:
    image: redis:7-alpine
    restart: unless-stopped
```

2. Create a `.env` file with your credentials (see [Configuration](#configuration)).

3. Start the service:
   ```bash
   docker-compose up -d
   ```

**Option B: Manual Docker Run**

```bash
docker run -d \
  --name crypto-sync-bot \
  -p 8080:8080 \
  --env-file .env \
  ghcr.io/uykb/synchronous-backend:latest
```

### 2. Frontend Deployment (Vercel)

1. Fork or clone this repository to your GitHub.
2. Log in to [Vercel](https://vercel.com) and click **"Add New..."** -> **"Project"**.
3. Import your `synchronous` repository.
4. **Configure Project:**
   - **Framework Preset:** Vite
   - **Root Directory:** `frontend` (Important!)
5. **Deploy**.
6. **Post-Deployment Configuration:**
   - Go to your Vercel Project Settings.
   - Update `frontend/vercel.json` in your repo to point to your actual backend URL:
     ```json
     {
       "rewrites": [
         {
           "source": "/api/:path*",
           "destination": "https://your-backend-domain.com/api/:path*"
         }
       ]
     }
     ```
   - Push the change to trigger a redeploy.

## Configuration

Create a `.env` file for the backend:

```bash
# Binance API
BINANCE_API_KEY=your_key
BINANCE_API_SECRET=your_secret
BINANCE_TESTNET=false

# Bybit API
BYBIT_API_KEY=your_key
BYBIT_API_SECRET=your_secret

# OKX API (Optional/Coming Soon)
OKX_API_KEY=
OKX_API_SECRET=
OKX_API_PASSPHRASE=

# Trading Config
SYMBOL=BTC-USDT
POSITION_RATIO=1.0      # Multiplier for order size
MAX_POSITION=1.0        # Max position size allowed
STOP_LOSS_RATIO=0.05    # 5% stop loss
```

## Development

### Local Development

1. **Start Backend:**
   ```bash
   cd crypto-sync-bot
   go mod tidy
   go run ./cmd/main.go
   ```

2. **Start Frontend:**
   ```bash
   cd frontend
   bun install
   bun dev
   ```
   Access at `http://localhost:5173`. The frontend proxies `/api` requests to `http://localhost:8080`.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/status` | Check bot status |
| GET | `/api/config` | Get current configuration |
| POST | `/api/signals` | Manually trigger a signal |

## License

MIT License
