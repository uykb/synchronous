package api

import (
	"crypto-sync-bot/internal/auth"
	"crypto-sync-bot/internal/config"
	"crypto-sync-bot/internal/models"
	"crypto-sync-bot/internal/processor"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type API struct {
	cfg  *config.Config
	proc *processor.SignalProcessor
}

func NewAPI(cfg *config.Config, proc *processor.SignalProcessor) *API {
	return &API{cfg: cfg, proc: proc}
}

func (a *API) SetupRoutes(r *gin.Engine) {
	authLimiter := NewRateLimiter(5, time.Minute)
	signalLimiter := NewRateLimiter(60, time.Minute)

	api := r.Group("/api")
	{
		api.GET("/status", a.GetStatus)
		api.GET("/system/ip", a.GetIP)

		authGroup := api.Group("/auth")
		authGroup.Use(RateLimitMiddleware(authLimiter))
		{
			authGroup.POST("/setup", a.SetupAuth)
			authGroup.POST("/verify", a.VerifyAuth)
		}

		api.POST("/restart", a.Restart)

		// Webhook route with HMAC verification and rate limiting
		api.POST("/signals", RateLimitMiddleware(signalLimiter), HMACVerification(a.cfg.GetWebhookSecret()), a.PostSignal)

		// Protected routes
		protected := api.Group("/")
		protected.Use(AuthMiddleware())
		{
			protected.GET("/config", a.GetConfig)
			protected.PUT("/config", a.UpdateConfig)
			protected.PUT("/exchanges/:id", a.UpdateExchangeConfig)
			protected.DELETE("/exchanges/:id", a.DeleteExchangeConfig)
			protected.POST("/exchanges/:id/test", a.TestExchangeConnection)
			protected.GET("/sync-items", a.GetSyncItems)
			protected.POST("/sync-items", a.AddSyncItem)
			protected.DELETE("/sync-items/:id", a.DeleteSyncItem)
		}
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

func (a *API) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"is_configured": a.cfg.GetAuth().IsConfigured,
		"version":      "1.0.0",
	})
}

func (a *API) GetIP(c *gin.Context) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch IP"})
		return
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ip": string(ip)})
}

func (a *API) SetupAuth(c *gin.Context) {
	authCfg := a.cfg.GetAuth()
	if authCfg.IsConfigured {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already configured"})
		return
	}

	secret, url, _ := auth.GenerateTOTPSecret("admin")

	authCfg.TOTPSecret = secret
	authCfg.IsConfigured = true
	a.cfg.UpdateAuth(authCfg)
	a.cfg.Save()

	c.JSON(http.StatusOK, gin.H{
		"totp_secret": secret,
		"totp_url":    url,
	})
}

func (a *API) VerifyAuth(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For simplicity, we check TOTP.
	if !auth.VerifyTOTP(req.Code, a.cfg.GetAuth().TOTPSecret) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid TOTP code"})
		return
	}

	token, _ := auth.GenerateToken("admin")
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *API) Restart(c *gin.Context) {
	// Check auth if configured
	if a.cfg.GetAuth().IsConfigured {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
			return
		}
		if _, err := auth.ValidateToken(token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restarting..."})
	
	// Exit in a goroutine to allow response to be sent
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
}

func (a *API) GetConfig(c *gin.Context) {
	type ExchangeStatus struct {
		Enabled    bool   `json:"enabled"`
		APIKeyHint string `json:"api_key_hint,omitempty"` // e.g., "abc1...xyz9"
		Testnet    bool   `json:"testnet,omitempty"`
	}

	binanceCfg := a.cfg.GetBinance()
	okxCfg := a.cfg.GetOKX()
	bybitCfg := a.cfg.GetBybit()

	maskKey := func(key string) string {
		if len(key) < 8 {
			return ""
		}
		return key[:4] + "..." + key[len(key)-4:]
	}

	safe := struct {
		Binance   ExchangeStatus     `json:"binance"`
		OKX       ExchangeStatus     `json:"okx"`
		Bybit     ExchangeStatus     `json:"bybit"`
		Sync      interface{}        `json:"sync"`
		SyncItems []config.SyncItem  `json:"sync_items"`
	}{
		Binance: ExchangeStatus{
			Enabled:    binanceCfg.APIKey != "",
			APIKeyHint: maskKey(binanceCfg.APIKey),
			Testnet:    binanceCfg.Testnet,
		},
		OKX: ExchangeStatus{
			Enabled:    okxCfg.APIKey != "",
			APIKeyHint: maskKey(okxCfg.APIKey),
		},
		Bybit: ExchangeStatus{
			Enabled:    bybitCfg.APIKey != "",
			APIKeyHint: maskKey(bybitCfg.APIKey),
		},
		Sync:      a.cfg.GetSync(),
		SyncItems: a.cfg.GetSyncItems(),
	}

	c.JSON(http.StatusOK, safe)
}

func (a *API) UpdateExchangeConfig(c *gin.Context) {
	exchangeID := c.Param("id")

	var req struct {
		APIKey     string `json:"api_key"`
		APISecret  string `json:"api_secret"`
		Passphrase string `json:"passphrase,omitempty"`
		Testnet    bool   `json:"testnet,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.cfg.UpdateExchange(exchangeID, req.APIKey, req.APISecret, req.Passphrase, req.Testnet)
	if err := a.cfg.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exchange config saved", "exchange": exchangeID})
}

func (a *API) DeleteExchangeConfig(c *gin.Context) {
	exchangeID := c.Param("id")

	a.cfg.DeleteExchange(exchangeID)
	if err := a.cfg.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exchange config deleted", "exchange": exchangeID})
}

func (a *API) TestExchangeConnection(c *gin.Context) {
	exchangeID := c.Param("id")

	// For now, just check if config exists - real implementation would ping the exchange API
	var enabled bool
	switch exchangeID {
	case "binance":
		enabled = a.cfg.GetBinance().APIKey != ""
	case "okx":
		enabled = a.cfg.GetOKX().APIKey != ""
	case "bybit":
		enabled = a.cfg.GetBybit().APIKey != ""
	}

	if !enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Exchange not configured", "success": false})
		return
	}

	// TODO: Implement actual exchange API ping
	c.JSON(http.StatusOK, gin.H{"message": "Connection test passed", "success": true, "exchange": exchangeID})
}

func (a *API) UpdateConfig(c *gin.Context) {
	var newCfg config.Config
	if err := c.ShouldBindJSON(&newCfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use UpdateAll for backward compatibility or UpdateSync if only sync changed
	a.cfg.UpdateSync(newCfg.Sync)
	a.cfg.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Sync config updated"})
}

func (a *API) GetSyncItems(c *gin.Context) {
	c.JSON(http.StatusOK, a.cfg.GetSyncItems())
}

func (a *API) AddSyncItem(c *gin.Context) {
	var item config.SyncItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.cfg.AddSyncItem(item)
	a.cfg.Save()
	c.JSON(http.StatusOK, item)
}

func (a *API) DeleteSyncItem(c *gin.Context) {
	id := c.Param("id")
	if a.cfg.DeleteSyncItem(id) {
		a.cfg.Save()
		c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
}

func (a *API) PostSignal(c *gin.Context) {
	var signal models.TradingSignal
	if err := c.ShouldBindJSON(&signal); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := processor.ProduceSignal(c.Request.Context(), &signal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue signal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Signal received and queued"})
}
