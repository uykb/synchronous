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
	bybitCfg := cfg.GetBybit()
	client := bybit.NewClient().
		WithAuth(bybitCfg.APIKey, bybitCfg.APISecret)
	
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
	symbolStr := bybit.SymbolV5(signal.Symbol) // Assuming signal.Symbol is like "BTCUSDT"

	// Create Order
	// Using Unified Margin or Linear Futures API
	res, err := e.client.V5().Order().CreateOrder(bybit.V5CreateOrderParam{
		Category: bybit.CategoryV5Linear,
		Symbol:   symbolStr,
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
		OrderID:   res.Result.OrderID,
		Timestamp: signal.Timestamp,
	}, nil
}

func (e *BybitExecutor) GetOrder(orderID, symbol string) (*models.OrderResult, error) {
	// Use GetOpenOrders to check status (Note: Filled orders might move to History)
	symbolStr := bybit.SymbolV5(symbol)
	res, err := e.client.V5().Order().GetOpenOrders(bybit.V5GetOpenOrdersParam{
		Category: bybit.CategoryV5Linear,
		Symbol:   &symbolStr,
		OrderID:  &orderID,
	})
	if err != nil {
		return nil, err
	}
	if len(res.Result.List) == 0 {
		// If not in open orders, check history
		// Note: GetHistory might be needed here if supported
		return nil, fmt.Errorf("order not found in open orders")
	}
	order := res.Result.List[0]
	return &models.OrderResult{
		Exchange: "Bybit",
		Symbol:   symbol,
		OrderID:  order.OrderID,
		Status:   string(order.OrderStatus),
	}, nil
}

func (e *BybitExecutor) Close() {
	// Cleanup
}
