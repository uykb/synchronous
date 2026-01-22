package exchange

import (
	"crypto-sync-bot/internal/config"
	"crypto-sync-bot/internal/models"
	"fmt"
	"log"

	"github.com/hirokisan/bybit/v2"
)

type BybitExecutor struct {
	client *bybit.Client
	config *config.Config
}

func NewBybitExecutor(cfg *config.Config) *BybitExecutor {
	client := bybit.NewClient().
		WithAuth(cfg.Bybit.APIKey, cfg.Bybit.APISecret)
	
	return &BybitExecutor{
		client: client,
		config: cfg,
	}
}

func (e *BybitExecutor) Name() string {
	return "Bybit"
}

func (e *BybitExecutor) PlaceOrder(signal *models.TradingSignal) (*models.OrderResult, error) {
	// Map Side
	var side bybit.Side
	if signal.Side == "BUY" {
		side = bybit.SideBuy
	} else {
		side = bybit.SideSell
	}

	// Map Type
	var orderType bybit.OrderType
	if signal.OrderType == "MARKET" {
		orderType = bybit.OrderTypeMarket
	} else {
		orderType = bybit.OrderTypeLimit
	}

	// Map Symbol (Bybit uses BTCUSDT usually for Linear)
	symbol := bybit.SymbolUSDT(signal.Symbol) // Assuming signal.Symbol is like "BTCUSDT"

	// Create Order
	// Using Unified Margin or Linear Futures API
	res, err := e.client.V5().Order().CreateOrder(bybit.V5CreateOrderParam{
		Category: bybit.CategoryV5Linear,
		Symbol:   bybit.SymbolV5(symbol),
		Side:     bybit.Side(side),
		OrderType: bybit.OrderType(orderType),
		Qty:      fmt.Sprintf("%f", signal.Quantity),
		Price:    func() *string {
			if signal.OrderType == "LIMIT" {
				p := fmt.Sprintf("%f", signal.Price)
				return &p
			}
			return nil
		}(),
	})

	if err != nil {
		log.Printf("Bybit Order Failed: %v", err)
		return &models.OrderResult{
			Exchange:     "Bybit",
			Symbol:       signal.Symbol,
			Status:       "failed",
			ErrorMessage: err.Error(),
			Timestamp:    signal.Timestamp,
		}, err
	}

	return &models.OrderResult{
		Exchange:  "Bybit",
		Symbol:    signal.Symbol,
		Status:    "success",
		OrderID:   res.Result.OrderId,
		Timestamp: signal.Timestamp,
	}, nil
}

func (e *BybitExecutor) GetOrder(orderID, symbol string) (*models.OrderResult, error) {
	res, err := e.client.V5().Order().GetHistory(bybit.V5GetHistoryParam{
		Category: bybit.CategoryV5Linear,
		Symbol:   bybit.SymbolV5(symbol),
		OrderId:  &orderID,
	})
	if err != nil {
		return nil, err
	}
	if len(res.Result.List) == 0 {
		return nil, fmt.Errorf("order not found")
	}
	order := res.Result.List[0]
	return &models.OrderResult{
		Exchange: "Bybit",
		Symbol:   symbol,
		OrderID:  order.OrderId,
		Status:   string(order.OrderStatus),
	}, nil
}

func (e *BybitExecutor) Close() {
	// Cleanup
}
