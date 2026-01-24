package exchange

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"crypto-sync-bot/internal/config"
	"crypto-sync-bot/internal/models"
)

const lighterBaseURL = "https://mainnet.zklighter.elliot.ai"

type LighterExecutor struct {
	config     *config.Config
	httpClient *http.Client
}

func NewLighterExecutor(cfg *config.Config) *LighterExecutor {
	return &LighterExecutor{
		config:     cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (e *LighterExecutor) Name() string {
	return "Lighter"
}

func (e *LighterExecutor) PlaceOrder(signal *models.TradingSignal) (*models.OrderResult, error) {
	lighterCfg := e.config.GetLighter()
	
	// Map side: Lighter uses IsAsk=0 for Buy, IsAsk=1 for Sell
	isAsk := 0
	if signal.Side == "SELL" {
		isAsk = 1
	}
	
	// Map order type: 0=Limit, 1=Market
	orderType := 1 // Market by default
	if signal.OrderType == "LIMIT" {
		orderType = 0
	}
	
	// Build order request
	// Note: Lighter uses market_id, we'll need to map symbols
	// For simplicity, assume symbol is already a market_id or use a mapping
	orderReq := map[string]interface{}{
		"tx_type": "CreateOrder",
		"tx_info": map[string]interface{}{
			"market_id":     getMarketID(signal.Symbol),
			"amount":        fmt.Sprintf("%f", signal.Quantity),
			"price":         fmt.Sprintf("%f", signal.Price),
			"is_ask":        isAsk,
			"type":          orderType,
			"account_index": lighterCfg.AccountIndex,
			"nonce":         time.Now().UnixNano(),
		},
	}
	
	respBody, err := e.signedRequest("POST", "/api/v1/sendTx", orderReq)
	if err != nil {
		log.Printf("Lighter Order Failed: %v", err)
		return &models.OrderResult{
			Exchange:     "Lighter",
			Symbol:       signal.Symbol,
			Status:       "failed",
			ErrorMessage: err.Error(),
			Timestamp:    signal.Timestamp,
		}, err
	}
	
	var resp struct {
		Result struct {
			TxHash string `json:"tx_hash"`
		} `json:"result"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	if resp.Error != "" {
		return &models.OrderResult{
			Exchange:     "Lighter",
			Symbol:       signal.Symbol,
			Status:       "failed",
			ErrorMessage: resp.Error,
			Timestamp:    signal.Timestamp,
		}, fmt.Errorf("lighter error: %s", resp.Error)
	}
	
	return &models.OrderResult{
		Exchange:  "Lighter",
		Symbol:    signal.Symbol,
		Status:    "success",
		OrderID:   resp.Result.TxHash,
		Timestamp: signal.Timestamp,
	}, nil
}

func (e *LighterExecutor) GetOrder(orderID, symbol string) (*models.OrderResult, error) {
	// Lighter uses tx_hash for order tracking
	// This would query the transaction status
	return &models.OrderResult{
		Exchange: "Lighter",
		Symbol:   symbol,
		OrderID:  orderID,
		Status:   "unknown",
	}, nil
}

func (e *LighterExecutor) Close() {
	// Cleanup if needed
}

func (e *LighterExecutor) signedRequest(method, path string, body interface{}) ([]byte, error) {
	lighterCfg := e.config.GetLighter()
	
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	
	// Create signature: HMAC-SHA256 of timestamp + body
	mac := hmac.New(sha256.New, []byte(lighterCfg.APISecret))
	mac.Write([]byte(timestamp))
	mac.Write(jsonBody)
	signature := hex.EncodeToString(mac.Sum(nil))
	
	req, err := http.NewRequest(method, lighterBaseURL+path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Lighter-API-Key", lighterCfg.APIKey)
	req.Header.Set("X-Lighter-Signature", signature)
	req.Header.Set("X-Lighter-Timestamp", timestamp)
	
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("lighter API error %d: %s", resp.StatusCode, string(respBody))
	}
	
	return respBody, nil
}

// getMarketID maps symbol to Lighter market_id
// Common mappings - in production this should query /api/v1/orderBooks
func getMarketID(symbol string) int {
	marketMap := map[string]int{
		"BTCUSDT": 1,
		"ETHUSDT": 2,
		"SOLUSDT": 3,
		"BTC_USDT": 1,
		"ETH_USDT": 2,
		"SOL_USDT": 3,
	}
	if id, ok := marketMap[symbol]; ok {
		return id
	}
	return 1 // Default to BTC market
}
