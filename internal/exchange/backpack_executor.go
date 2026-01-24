package exchange

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"crypto-sync-bot/internal/config"
	"crypto-sync-bot/internal/models"
)

const backpackBaseURL = "https://api.backpack.exchange"

type BackpackExecutor struct {
	config     *config.Config
	httpClient *http.Client
	privateKey ed25519.PrivateKey
}

func NewBackpackExecutor(cfg *config.Config) (*BackpackExecutor, error) {
	backpackCfg := cfg.GetBackpack()
	
	// If API secret is not configured, return nil (will be configured later via admin panel)
	if backpackCfg.APISecret == "" {
		return nil, nil
	}
	
	// Decode the secret key (Base64 encoded Ed25519 private key)
	privateKeyBytes, err := base64.StdEncoding.DecodeString(backpackCfg.APISecret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode backpack secret: %w", err)
	}
	
	// Ed25519 private key should be 64 bytes (or 32 byte seed)
	var privateKey ed25519.PrivateKey
	if len(privateKeyBytes) == 64 {
		privateKey = ed25519.PrivateKey(privateKeyBytes)
	} else if len(privateKeyBytes) == 32 {
		privateKey = ed25519.NewKeyFromSeed(privateKeyBytes)
	} else {
		return nil, fmt.Errorf("invalid backpack secret key length: %d", len(privateKeyBytes))
	}
	
	return &BackpackExecutor{
		config:     cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		privateKey: privateKey,
	}, nil
}

func (e *BackpackExecutor) Name() string {
	return "Backpack"
}

func (e *BackpackExecutor) PlaceOrder(signal *models.TradingSignal) (*models.OrderResult, error) {
	// Map side: Backpack uses "Bid" for buy, "Ask" for sell
	side := "Bid"
	if signal.Side == "SELL" {
		side = "Ask"
	}
	
	// Map order type
	orderType := "Market"
	if signal.OrderType == "LIMIT" {
		orderType = "Limit"
	}
	
	// Convert symbol format: BTCUSDT -> BTC_USDT
	symbol := convertSymbol(signal.Symbol)
	
	// Build order params
	params := map[string]string{
		"symbol":    symbol,
		"side":      side,
		"orderType": orderType,
		"quantity":  fmt.Sprintf("%f", signal.Quantity),
	}
	
	if orderType == "Limit" {
		params["price"] = fmt.Sprintf("%f", signal.Price)
		params["timeInForce"] = "GTC"
	}
	
	// Make request
	respBody, err := e.signedRequest("POST", "/api/v1/order", "orderExecute", params)
	if err != nil {
		log.Printf("Backpack Order Failed: %v", err)
		return &models.OrderResult{
			Exchange:     "Backpack",
			Symbol:       signal.Symbol,
			Status:       "failed",
			ErrorMessage: err.Error(),
			Timestamp:    signal.Timestamp,
		}, err
	}
	
	// Parse response
	var resp struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &models.OrderResult{
		Exchange:  "Backpack",
		Symbol:    signal.Symbol,
		Status:    "success",
		OrderID:   resp.ID,
		Timestamp: signal.Timestamp,
	}, nil
}

func (e *BackpackExecutor) GetOrder(orderID, symbol string) (*models.OrderResult, error) {
	params := map[string]string{
		"orderId": orderID,
		"symbol":  convertSymbol(symbol),
	}
	
	respBody, err := e.signedRequest("GET", "/api/v1/order", "orderQuery", params)
	if err != nil {
		return nil, err
	}
	
	var resp struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}
	
	return &models.OrderResult{
		Exchange: "Backpack",
		Symbol:   symbol,
		OrderID:  resp.ID,
		Status:   resp.Status,
	}, nil
}

func (e *BackpackExecutor) Close() {
	// Cleanup if needed
}

// signedRequest makes an authenticated request to Backpack API
func (e *BackpackExecutor) signedRequest(method, path, instruction string, params map[string]string) ([]byte, error) {
	backpackCfg := e.config.GetBackpack()
	
	timestamp := time.Now().UnixMilli()
	window := int64(5000)
	
	// Build signing string: instruction=X&param1=v1&param2=v2&timestamp=T&window=W
	// Parameters must be sorted alphabetically
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	
	var signParts []string
	signParts = append(signParts, fmt.Sprintf("instruction=%s", instruction))
	for _, k := range keys {
		signParts = append(signParts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	signParts = append(signParts, fmt.Sprintf("timestamp=%d", timestamp))
	signParts = append(signParts, fmt.Sprintf("window=%d", window))
	
	signString := strings.Join(signParts, "&")
	
	// Sign with Ed25519
	signature := ed25519.Sign(e.privateKey, []byte(signString))
	signatureB64 := base64.StdEncoding.EncodeToString(signature)
	
	// Build request
	var req *http.Request
	var err error
	
	url := backpackBaseURL + path
	
	if method == "GET" {
		// For GET, params go in query string
		url += "?" + buildQueryString(params)
		req, err = http.NewRequest(method, url, nil)
	} else {
		// For POST, params go in JSON body
		jsonBody, _ := json.Marshal(params)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	}
	
	if err != nil {
		return nil, err
	}
	
	// Set auth headers
	req.Header.Set("X-API-Key", backpackCfg.APIKey)
	req.Header.Set("X-Signature", signatureB64)
	req.Header.Set("X-Timestamp", strconv.FormatInt(timestamp, 10))
	req.Header.Set("X-Window", strconv.FormatInt(window, 10))
	
	// Execute
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("backpack API error %d: %s", resp.StatusCode, string(body))
	}
	
	return body, nil
}

func buildQueryString(params map[string]string) string {
	var parts []string
	for k, v := range params {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(parts, "&")
}

// convertSymbol converts BTCUSDT to BTC_USDT format
func convertSymbol(symbol string) string {
	// Common patterns: BTCUSDT -> BTC_USDT, SOLUSDC -> SOL_USDC
	quotes := []string{"USDT", "USDC", "USD", "BTC", "ETH"}
	for _, quote := range quotes {
		if strings.HasSuffix(symbol, quote) {
			base := strings.TrimSuffix(symbol, quote)
			return base + "_" + quote
		}
	}
	return symbol
}
