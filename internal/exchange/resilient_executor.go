package exchange

import (
	"crypto-sync-bot/internal/models"
	"github.com/sony/gobreaker"
	"net"
	"strings"
)

// IsSuccessful returns true if the error is NOT a network error or a 5xx server error.
func IsSuccessful(err error) bool {
	if err == nil {
		return true
	}

	// Check for network errors
	if _, ok := err.(net.Error); ok {
		return false
	}

	// Check for 5xx errors by looking at common error message patterns
	// This is a heuristic as different exchanges return different error formats
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "500") || strings.Contains(errStr, "502") || 
	   strings.Contains(errStr, "503") || strings.Contains(errStr, "504") ||
	   strings.Contains(errStr, "internal server error") || 
	   strings.Contains(errStr, "service unavailable") {
		return false
	}

	return true
}

type ResilientExecutor struct {
	executor models.ExchangeExecutor
	cb       *gobreaker.CircuitBreaker
}

func NewResilientExecutor(executor models.ExchangeExecutor, cb *gobreaker.CircuitBreaker) *ResilientExecutor {
	return &ResilientExecutor{
		executor: executor,
		cb:       cb,
	}
}

func (r *ResilientExecutor) Name() string {
	return r.executor.Name()
}

func (r *ResilientExecutor) PlaceOrder(signal *models.TradingSignal) (*models.OrderResult, error) {
	result, err := r.cb.Execute(func() (interface{}, error) {
		return r.executor.PlaceOrder(signal)
	})
	if err != nil {
		return nil, err
	}
	return result.(*models.OrderResult), nil
}

func (r *ResilientExecutor) GetOrder(orderID, symbol string) (*models.OrderResult, error) {
	result, err := r.cb.Execute(func() (interface{}, error) {
		return r.executor.GetOrder(orderID, symbol)
	})
	if err != nil {
		return nil, err
	}
	return result.(*models.OrderResult), nil
}

func (r *ResilientExecutor) Close() {
	r.executor.Close()
}
