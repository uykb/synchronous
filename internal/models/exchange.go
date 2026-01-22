package models

type ExchangeExecutor interface {
	Name() string
	PlaceOrder(signal *TradingSignal) (*OrderResult, error)
	GetOrder(orderID, symbol string) (*OrderResult, error)
	Close()
}
