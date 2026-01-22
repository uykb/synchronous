package models

type TradingSignal struct {
	Symbol          string  `json:"symbol"`
	Side            string  `json:"side"`        // "BUY" or "SELL"
	OrderType       string  `json:"order_type"`  // "MARKET" or "LIMIT"
	Quantity        float64 `json:"quantity"`
	Price           float64 `json:"price"`       // Limit price
	Leverage        int     `json:"leverage"`    // Leverage multiplier
	StopLossPrice   float64 `json:"stop_loss"`   // Stop loss price
	TakeProfitPrice float64 `json:"take_profit"` // Take profit price
	Timestamp       int64   `json:"timestamp"`
	SignalID        string  `json:"signal_id"`
	Source          string  `json:"source"` // "binance"
}

type OrderResult struct {
	Exchange     string `json:"exchange"`
	Symbol       string `json:"symbol"`
	OrderID      string `json:"order_id"`
	Status       string `json:"status"` // "success" or "failed"
	ErrorMessage string `json:"error_message"`
	Timestamp    int64  `json:"timestamp"`
}
