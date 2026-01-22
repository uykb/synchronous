# Crypto-Sync-Bot

一个实时加密货币交易信号同步机器人，监听交易信号并在多个交易所（Binance, OKX, Bybit）自动执行订单。

## 架构设计

本项目采用前后端分离架构，以实现最佳性能和灵活部署：

```
┌─────────────────────────────────────────────────────────────────┐
│                        Crypto-Sync-Bot                          │
├────────────────────────────┬────────────────────────────────────┤
│      前端 (Vercel)         │         后端 (Docker)              │
│   Vue 3 + TypeScript       │       Go 1.24 + Redis              │
│   - 仪表盘 UI              │   - 信号处理                       │
│   - 设置面板               │   - 交易所适配器                   │
│   - 实时图表               │   - 风险管理                       │
│            ↓               │                ↓                   │
│      Vercel Edge           │        Docker 容器                 │
│   (代理 /api -> 后端)      │      (暴露 8080 端口)              │
└────────────────────────────┴────────────────────────────────────┘
```

## 技术栈

| 层级 | 技术 | 用途 |
|-------|------------|---------|
| **后端** | Go 1.24 | 高性能信号处理与订单执行 |
| **前端** | Vue 3 + TypeScript | 响应式仪表盘与配置界面 |
| **数据库** | MySQL + Redis | 配置/数据持久化与消息队列 |
| **部署** | Docker + Vercel | 容器化后端，Serverless 前端 |

## 前置要求

- **Docker** & **Docker Compose** (用于后端部署)
- **MySQL** (外部数据库，用于持久化配置)
- **Vercel 账号** (用于前端部署)

## 🚀 部署指南

### 1. 后端部署 (Docker)

后端已容器化，可直接从 GitHub Container Registry 拉取。

**推荐方式：使用 Docker Compose**

1. 在服务器上创建 `docker-compose.yml` 文件：

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
      - MYSQL_DSN=${MYSQL_DSN}
      - REDIS_ADDR=redis:6379
      # 初始配置 (首次启动时写入数据库)
      - BINANCE_API_KEY=${BINANCE_API_KEY}
      - BINANCE_API_SECRET=${BINANCE_API_SECRET}
      - SYMBOL=BTC-USDT
    depends_on:
      - redis

  redis:
    image: redis:7-alpine
    restart: unless-stopped
```

2. 创建 `.env` 文件并填入配置（见[配置说明](#配置)）。

3. 启动服务：
   ```bash
   docker-compose up -d
   ```

### 2. 前端部署 (Vercel)

1. Fork 或 Clone 本仓库到你的 GitHub。
2. 登录 [Vercel](https://vercel.com) 点击 **"Add New..."** -> **"Project"**.
3. 导入你的 `synchronous` 仓库。
4. **配置项目:**
   - **Framework Preset:** Vite
   - **Root Directory:** `frontend` (重要!)
5. **部署 (Deploy)**.
6. **部署后配置:**
   - 进入 Vercel 项目设置 (Settings)。
   - 在你的代码仓库中修改 `frontend/vercel.json`，将目标地址改为你后端的实际域名/IP：
     ```json
     {
       "rewrites": [
         {
           "source": "/api/:path*",
           "destination": "http://your-backend-ip:8080/api/:path*"
         }
       ]
     }
     ```
   - 推送代码更改以触发重新部署。

## 配置说明

### 环境变量 (.env)

后端通过环境变量进行**首次初始化**。启动后，配置将持久化存储在 MySQL 数据库中，后续修改请通过前端界面进行。

```bash
# 数据库连接 (必须)
MYSQL_DSN="user:password@tcp(your-mysql-host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

# 交易所 API (用于初始化)
BINANCE_API_KEY=your_key
BINANCE_API_SECRET=your_secret
BINANCE_TESTNET=false

BYBIT_API_KEY=your_key
BYBIT_API_SECRET=your_secret

# 交易配置
SYMBOL=BTC-USDT
POSITION_RATIO=1.0      # 仓位比例
MAX_POSITION=1.0        # 最大持仓限制
STOP_LOSS_RATIO=0.05    # 止损比例 (5%)
```

### 动态配置

系统启动后，你可以访问前端页面：
1. 登录设置面板。
2. 修改 API Key 或同步策略。
3. 点击保存并重启，新配置将立即生效（无需修改服务器环境变量）。

## 开发指南

### 本地开发

1. **启动后端:**
   ```bash
   cd crypto-sync-bot
   # 确保本地有 MySQL 和 Redis 运行
   export MYSQL_DSN="root:root@tcp(localhost:3306)/crypto_bot?charset=utf8mb4&parseTime=True&loc=Local"
   go run ./cmd/main.go
   ```

2. **启动前端:**
   ```bash
   cd frontend
   bun install
   bun dev
   ```
   访问 `http://localhost:5173`。前端已配置代理，将 `/api` 请求转发至 `http://localhost:8080`。

## API 接口

| 方法 | 路径 | 描述 |
|--------|----------|-------------|
| GET | `/api/status` | 查看机器人状态 |
| GET | `/api/config` | 获取当前配置 |
| POST | `/api/restart` | 重启服务 (应用新配置) |
| POST | `/api/signals` | 手动触发信号 |

## 许可证

MIT License
