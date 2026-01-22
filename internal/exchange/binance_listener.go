package exchange

import (
	"context"
	"crypto-sync-bot/internal/config"
	"crypto-sync-bot/internal/models"
	"crypto-sync-bot/internal/processor"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2/futures"
)

type BinanceListener struct {
	client   *futures.Client // Using Futures Client for API calls if needed
	config   *config.Config
	mu       sync.Mutex
	running  bool
	stopChan chan struct{}
}

func NewBinanceListener(cfg *config.Config) *BinanceListener {
	return &BinanceListener{
		config:   cfg,
		stopChan: make(chan struct{}),
	}
}

func (b *BinanceListener) Start() error {
	b.mu.Lock()
	if b.running {
		b.mu.Unlock()
		return nil
	}
	b.running = true
	b.mu.Unlock()

	// Configure Testnet if needed
	if b.config.Binance.Testnet {
		futures.UseTestnet = true
	}

	// Initialize the client (optional, mostly for REST calls)
	b.client = futures.NewClient(b.config.Binance.APIKey, b.config.Binance.APISecret)

	go b.connectWebSocket()

	return nil
}

func (b *BinanceListener) connectWebSocket() {
	errHandler := func(err error) {
		log.Printf("Binance WebSocket Error: %v", err)
	}

	for {
		select {
		case <-b.stopChan:
			return
		default:
			listenKey, err := b.client.NewStartUserStreamService().Do(context.Background())
			if err != nil {
				log.Printf("Error getting listenKey: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			log.Printf("Starting Binance User Stream with ListenKey: %s", listenKey)

			// KeepAlive for ListenKey
			go func(lk string) {
				ticker := time.NewTicker(30 * time.Minute)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						err := b.client.NewKeepaliveUserStreamService().ListenKey(lk).Do(context.Background())
						if err != nil {
							log.Printf("Error keeping alive listenKey: %v", err)
						}
					case <-b.stopChan:
						return
					}
				}
			}(listenKey)

			doneC, stopC, err := futures.WsFuturesUserDataServe(listenKey, func(event *futures.WsFuturesUserDataEvent) {
				if event.Event == "ORDER_TRADE_UPDATE" {
					trade := event.OrderTradeUpdate
					if trade.Status == "FILLED" && trade.Symbol == b.config.Sync.Symbol {
						signal := &models.TradingSignal{
							SignalID:  strconv.FormatInt(trade.OrderID, 10),
							Symbol:    trade.Symbol,
							Side:      string(trade.Side),
							OrderType: string(trade.Type),
							Quantity:  mustParseFloat(trade.AccumulatedFilledQty),
							Price:     mustParseFloat(trade.AveragePrice),
							Timestamp: event.Time,
							Source:    "binance",
						}
						if err := processor.ProduceSignal(context.Background(), signal); err != nil {
							log.Printf("Error producing signal from Binance: %v", err)
						}
					}
				}
			}, errHandler)

			if err != nil {
				log.Printf("Error connecting to WebSocket: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			<-doneC
			close(stopC)
			log.Println("Binance WebSocket disconnected, reconnecting...")
			time.Sleep(2 * time.Second)
		}
	}
}

func (b *BinanceListener) Stop() {
	close(b.stopChan)
}

func mustParseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
