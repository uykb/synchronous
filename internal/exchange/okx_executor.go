package exchange

import (
	"crypto-sync-bot/internal/config"
	"crypto-sync-bot/internal/models"
	"fmt"
	"log"

	"github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/options"
)

type OKXExecutor struct {
	client *goex.OKX
	config *config.Config
}

func NewOKXExecutor(cfg *config.Config) *OKXExecutor {
	// Initialize OKX client
	// Note: goex/v2 structure might vary slightly, adapting to common usage
	client := goex.NewOKX(
		options.WithApiKey(cfg.OKX.APIKey),
		options.WithApiSecretKey(cfg.OKX.APISecret),
		options.WithPassphrase(cfg.OKX.Passphrase),
	)
	
	return &OKXExecutor{
		client: client,
		config: cfg,
	}
}

func (e *OKXExecutor) Name() string {
	return "OKX"
}

func (e *OKXExecutor) PlaceOrder(signal *models.TradingSignal) (*models.OrderResult, error) {
	// Map Side
	var side goex.OrderSide
	if signal.Side == "BUY" {
		side = goex.Spot_Buy // For Futures, might need mapping to specific Futures constants if different
	} else {
		side = goex.Spot_Sell
	}

	// Map Type
	var orderType goex.OrderType
	if signal.OrderType == "MARKET" {
		orderType = goex.OrderType_Market
	} else {
		orderType = goex.OrderType_Limit
	}
	
	// Create Order
	// Assuming BTC-USDT-SWAP for Futures, you might need to adjust symbol mapping
	symbol := signal.Symbol // e.g., "BTC-USDT"
	// OKX Futures often use "BTC-USDT-SWAP"
	
	// Execute Order (Swap/Futures)
	// goex/v2 organizes by Spot/Swap/Future. We assume Swap/Perpetual for "Futures" here
	result, err := e.client.Swap.CreateOrder(
		symbol, 
		side, 
		orderType, 
		signal.Quantity, 
		signal.Price, 
		options.WithOrderType(orderType),
	)

	if err != nil {
		log.Printf("OKX Order Failed: %v", err)
		return &models.OrderResult{
			Exchange:     "OKX",
			Symbol:       signal.Symbol,
			Status:       "failed",
			ErrorMessage: err.Error(),
			Timestamp:    signal.Timestamp,
		}, err
	}

	return &models.OrderResult{
		Exchange:  "OKX",
		Symbol:    signal.Symbol,
		Status:    "success",
		OrderID:   result.OrderID,
		Timestamp: signal.Timestamp,
	}, nil
}

func (e *OKXExecutor) GetOrder(orderID, symbol string) (*models.OrderResult, error) {
	// goex/v2 GetOrder for Swap
	result, err := e.client.Swap.GetOrder(symbol, orderID)
	if err != nil {
		return nil, err
	}

	return &models.OrderResult{
		Exchange: "OKX",
		Symbol:   symbol,
		OrderID:  result.OrderID,
		Status:   string(result.Status),
	}, nil
}

func (e *OKXExecutor) Close() {
	// Cleanup if needed
}
