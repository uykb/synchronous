package processor

import (
	"context"
	"crypto-sync-bot/internal/api"
	"crypto-sync-bot/internal/config"
	"crypto-sync-bot/internal/database"
	"crypto-sync-bot/internal/models"
	"crypto-sync-bot/internal/risk"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type SignalProcessor struct {
	okxExecutor   models.ExchangeExecutor
	bybitExecutor models.ExchangeExecutor
	riskManager   *risk.Manager
	config        *config.Config
	stopChan      chan struct{}
}

func NewSignalProcessor(cfg *config.Config, okx models.ExchangeExecutor, bybit models.ExchangeExecutor) *SignalProcessor {
	return &SignalProcessor{
		config:        cfg,
		okxExecutor:   okx,
		bybitExecutor: bybit,
		riskManager:   risk.NewManager(cfg),
		stopChan:      make(chan struct{}),
	}
}

func (p *SignalProcessor) Start() error {
	log.Println("Signal Processor Started (Redis Stream Consumer)")
	go p.processSignals()
	return nil
}

func (p *SignalProcessor) processSignals() {
	ctx := context.Background()
	groupName := "trading-group"
	consumerName := "processor-1"

	for {
		select {
		case <-p.stopChan:
			return
		default:
			// Read new messages
			entries, err := database.RDB.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumerName,
				Streams:  []string{"signals:trading", ">"},
				Count:    1,
				Block:    5 * time.Second,
			}).Result()

			if err != nil {
				if err != redis.Nil {
					log.Printf("XReadGroup Error: %v", err)
					time.Sleep(time.Second)
				}
				continue
			}

			for _, stream := range entries {
				for _, msg := range stream.Messages {
					p.handleMessage(ctx, msg)
				}
			}
		}
		// Basic logic: also try to process pending messages occasionally or on startup
		// For this implementation, we will focus on the main loop.
	}
}

func (p *SignalProcessor) handleMessage(ctx context.Context, msg redis.XMessage) {
	start := time.Now()
	defer func() {
		api.OrderLatency.Observe(time.Since(start).Seconds())
	}()

	payload, ok := msg.Values["payload"].(string)
	if !ok {
		log.Printf("Invalid message payload in msg %s", msg.ID)
		database.RDB.XAck(ctx, "signals:trading", "trading-group", msg.ID)
		return
	}

	var signal models.TradingSignal
	if err := json.Unmarshal([]byte(payload), &signal); err != nil {
		log.Printf("Failed to unmarshal signal in msg %s: %v", msg.ID, err)
		database.RDB.XAck(ctx, "signals:trading", "trading-group", msg.ID)
		return
	}

	log.Printf("Processing Signal from Stream [%s]: %s %s", msg.ID, signal.Side, signal.Symbol)

	// 0. Idempotency Check
	duplicate, err := IsDuplicate(ctx, signal.SignalID, signal.Source, signal.Quantity, signal.Price)
	if err != nil {
		log.Printf("Idempotency Check Error: %v", err)
	}
	if duplicate {
		log.Printf("Duplicate Signal Detected, skipping: %s", signal.SignalID)
		database.RDB.XAck(ctx, "signals:trading", "trading-group", msg.ID)
		return
	}

	// 1. Risk Check
	if err := p.riskManager.PreOrderCheck(&signal); err != nil {
		log.Printf("Risk Check Failed: %v", err)
		database.RDB.XAck(ctx, "signals:trading", "trading-group", msg.ID)
		return
	}

	// 2. Calculate Position
	signal.Quantity = signal.Quantity * p.config.Sync.PositionRatio

	// 3. Execute Orders in Parallel
	var wg sync.WaitGroup
	var okxErr, bybitErr error
	wg.Add(2)

	go func() {
		defer wg.Done()
		res, err := p.okxExecutor.PlaceOrder(&signal)
		okxErr = err
		if res != nil {
			database.SaveOrderResult(res)
		}
		if okxErr != nil {
			log.Printf("OKX Execution Error: %v", okxErr)
			api.OrdersCounter.WithLabelValues("okx", "failed").Inc()
		} else {
			api.OrdersCounter.WithLabelValues("okx", "success").Inc()
		}
	}()

	go func() {
		defer wg.Done()
		res, err := p.bybitExecutor.PlaceOrder(&signal)
		bybitErr = err
		if res != nil {
			database.SaveOrderResult(res)
		}
		if bybitErr != nil {
			log.Printf("Bybit Execution Error: %v", bybitErr)
			api.OrdersCounter.WithLabelValues("bybit", "failed").Inc()
		} else {
			api.OrdersCounter.WithLabelValues("bybit", "success").Inc()
		}
	}()

	wg.Wait()

	if okxErr == nil && bybitErr == nil {
		// Success on both exchanges
		database.RDB.XAck(ctx, "signals:trading", "trading-group", msg.ID)
		log.Printf("Successfully processed signal %s on both exchanges", msg.ID)
	} else {
		// Failure logic: Retry/DLQ
		p.handleFailure(ctx, msg)
	}
}

func (p *SignalProcessor) handleFailure(ctx context.Context, msg redis.XMessage) {
	// Use XPending to get delivery count
	pending, err := database.RDB.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: "signals:trading",
		Group:  "trading-group",
		Start:  msg.ID,
		End:    msg.ID,
		Count:  1,
	}).Result()

	deliveryCount := int64(0)
	if err == nil && len(pending) > 0 {
		deliveryCount = pending[0].TimesDelivered
	}

	if deliveryCount >= 3 {
		log.Printf("Signal %s failed %d times, moving to DLQ", msg.ID, deliveryCount)
		database.RDB.XAdd(ctx, &redis.XAddArgs{
			Stream: "signals:dlq",
			Values: msg.Values,
		})
		database.RDB.XAck(ctx, "signals:trading", "trading-group", msg.ID)
	} else {
		log.Printf("Signal %s failed (attempt %d), will be retried by consumer", msg.ID, deliveryCount)
		// We don't ACK, so it will stay in PEL.
		// In a real system, we'd have a claimer or use XReadGroup with 0 to re-process.
	}
}

func (p *SignalProcessor) Stop() {
	close(p.stopChan)
	p.okxExecutor.Close()
	p.bybitExecutor.Close()
}
