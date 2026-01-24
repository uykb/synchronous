package main

import (
	"context"
	"crypto-sync-bot/internal/api"
	"crypto-sync-bot/internal/config"
	"crypto-sync-bot/internal/database"
	"crypto-sync-bot/internal/exchange"
	"crypto-sync-bot/internal/models"
	"crypto-sync-bot/internal/processor"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
)

func main() {
	// 0. Initialize MySQL (for Config)
	if err := database.InitMySQL(); err != nil {
		log.Printf("Warning: MySQL initialization failed (using Env/File config): %v", err)
	}

	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Redis (optional)
	if err := database.InitRedis(); err != nil {
		log.Printf("Warning: Redis initialization failed: %v", err)
	} else if database.RDB == nil {
		log.Println("Note: Redis not configured (REDIS_ADDR not set)")
	}

	// Initialize SQLite
	if err := database.InitSQLite("./trading.db"); err != nil {
		log.Printf("Warning: SQLite initialization failed: %v", err)
	}

	// 2. Ensure Consumer Group exists (only if Redis is available)
	if database.RDB != nil {
		err = database.RDB.XGroupCreateMkStream(context.Background(), "signals:trading", "trading-group", "$").Err()
		if err != nil {
			log.Printf("Note: Consumer group setup: %v (usually means it already exists)", err)
		}
	}

	// 3. Initialize Executors and Processor
	okxRaw := exchange.NewOKXExecutor(cfg)
	bybitRaw := exchange.NewBybitExecutor(cfg)
	backpackRaw, err := exchange.NewBackpackExecutor(cfg)
	if err != nil {
		log.Printf("Warning: Backpack executor disabled: %v", err)
	} else if backpackRaw == nil {
		log.Println("Note: Backpack executor not configured (API keys can be set via admin panel)")
	}
	lighterRaw := exchange.NewLighterExecutor(cfg)

	okxCB := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:         "OKX",
		IsSuccessful: exchange.IsSuccessful,
	})
	bybitCB := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:         "Bybit",
		IsSuccessful: exchange.IsSuccessful,
	})
	lighterCB := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:         "Lighter",
		IsSuccessful: exchange.IsSuccessful,
	})

	okxExecutor := exchange.NewResilientExecutor(okxRaw, okxCB)
	bybitExecutor := exchange.NewResilientExecutor(bybitRaw, bybitCB)
	lighterExecutor := exchange.NewResilientExecutor(lighterRaw, lighterCB)

	// Backpack executor is optional
	var backpackExecutor models.ExchangeExecutor
	if backpackRaw != nil {
		backpackCB := gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:         "Backpack",
			IsSuccessful: exchange.IsSuccessful,
		})
		backpackExecutor = exchange.NewResilientExecutor(backpackRaw, backpackCB)
	}

	proc := processor.NewSignalProcessor(cfg, okxExecutor, bybitExecutor, backpackExecutor, lighterExecutor)

	// 4. Start Binance Listener (Produces to Redis)
	binanceListener := exchange.NewBinanceListener(cfg)
	if err := binanceListener.Start(); err != nil {
		log.Printf("Warning: Failed to start Binance listener: %v", err)
	}

	// 5. Start Signal Processor (Consumes from Redis)
	if err := proc.Start(); err != nil {
		log.Printf("Warning: Failed to start processor: %v", err)
	}

	// 6. Start Reconciler
	executors := []models.ExchangeExecutor{okxExecutor, bybitExecutor, lighterExecutor}
	if backpackExecutor != nil {
		executors = append(executors, backpackExecutor)
	}
	reconciler := processor.NewReconciler(executors)
	ctx, cancel := context.WithCancel(context.Background())
	go reconciler.Start(ctx)

	// 7. Initialize API
	r := gin.Default()
	
	// Add CORS middleware
	r.Use(CORSMiddleware())

	apiHandler := api.NewAPI(cfg, proc)
	apiHandler.SetupRoutes(r)

	// Run API in background
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to run API server: %v", err)
		}
	}()

	// 7. Wait for Shutdown Signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	cancel()
	binanceListener.Stop()
	proc.Stop()
	log.Println("Shutdown complete")
}

func CORSMiddleware() gin.HandlerFunc {
	allowedOrigin := os.Getenv("CORS_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173" // Development default
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow the configured origin or match it
		if origin == allowedOrigin || allowedOrigin == "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
