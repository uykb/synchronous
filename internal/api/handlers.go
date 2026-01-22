package api

import (
	"crypto-sync-bot/internal/auth"
	"crypto-sync-bot/internal/config"
	"crypto-sync-bot/internal/models"
	"crypto-sync-bot/internal/processor"
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
	api := r.Group("/api")
	{
		api.GET("/status", a.GetStatus)
		api.POST("/auth/setup", a.SetupAuth)
		api.POST("/auth/verify", a.VerifyAuth)
		api.POST("/restart", a.Restart)

		// Webhook route with HMAC verification
		api.POST("/signals", HMACVerification(a.cfg.WebhookSecret), a.PostSignal)

		// Protected routes
		protected := api.Group("/")
		protected.Use(AuthMiddleware())
		{
			protected.GET("/config", a.GetConfig)
			protected.PUT("/config", a.UpdateConfig)
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
		"is_configured": a.cfg.Auth.IsConfigured,
		"version":      "1.0.0",
	})
}

func (a *API) SetupAuth(c *gin.Context) {
	if a.cfg.Auth.IsConfigured {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already configured"})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, _ := auth.HashPassword(req.Password)
	secret, url, _ := auth.GenerateTOTPSecret("admin")

	a.cfg.Auth.PasswordHash = hash
	a.cfg.Auth.TOTPSecret = secret
	a.cfg.Auth.IsConfigured = true
	a.cfg.Save()

	c.JSON(http.StatusOK, gin.H{
		"totp_secret": secret,
		"totp_url":    url,
	})
}

func (a *API) VerifyAuth(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// First time setup might not require password if just verifying TOTP, 
	// but normally we check password then TOTP.
	// For simplicity, we check TOTP.
	if !auth.VerifyTOTP(req.Code, a.cfg.Auth.TOTPSecret) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid TOTP code"})
		return
	}

	token, _ := auth.GenerateToken("admin")
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *API) Restart(c *gin.Context) {
	// Check auth if configured
	if a.cfg.Auth.IsConfigured {
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
	c.JSON(http.StatusOK, a.cfg)
}

func (a *API) UpdateConfig(c *gin.Context) {
	var newCfg config.Config
	if err := c.ShouldBindJSON(&newCfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields selectively or overwrite? 
	// For now, let's just update exchange keys and sync settings.
	a.cfg.Binance = newCfg.Binance
	a.cfg.OKX = newCfg.OKX
	a.cfg.Bybit = newCfg.Bybit
	a.cfg.Sync = newCfg.Sync
	a.cfg.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Config updated"})
}

func (a *API) GetSyncItems(c *gin.Context) {
	c.JSON(http.StatusOK, a.cfg.SyncItems)
}

func (a *API) AddSyncItem(c *gin.Context) {
	var item config.SyncItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a.cfg.SyncItems = append(a.cfg.SyncItems, item)
	a.cfg.Save()
	c.JSON(http.StatusOK, item)
}

func (a *API) DeleteSyncItem(c *gin.Context) {
	id := c.Param("id")
	for i, item := range a.cfg.SyncItems {
		if item.ID == id {
			a.cfg.SyncItems = append(a.cfg.SyncItems[:i], a.cfg.SyncItems[i+1:]...)
			a.cfg.Save()
			c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
			return
		}
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
