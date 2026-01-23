package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
	"crypto-sync-bot/internal/auth"
	"crypto-sync-bot/internal/database"
)

type AuthConfig struct {
	TOTPSecret   string `json:"totp_secret" mapstructure:"totp_secret"`
	IsConfigured bool   `json:"is_configured" mapstructure:"is_configured"`
}

type SyncItem struct {
	ID      string   `json:"id" mapstructure:"id"`
	Name    string   `json:"name" mapstructure:"name"`
	Enabled bool     `json:"enabled" mapstructure:"enabled"`
	Source  string   `json:"source" mapstructure:"source"`
	Targets []string `json:"targets" mapstructure:"targets"`
	Symbol  string   `json:"symbol" mapstructure:"symbol"`
}

type BinanceConfig struct {
	APIKey    string `json:"api_key" mapstructure:"api_key"`
	APISecret string `json:"api_secret" mapstructure:"api_secret"`
	Testnet   bool   `json:"testnet" mapstructure:"testnet"`
}

type OKXConfig struct {
	APIKey     string `json:"api_key" mapstructure:"api_key"`
	APISecret  string `json:"api_secret" mapstructure:"api_secret"`
	Passphrase string `json:"passphrase" mapstructure:"passphrase"`
}

type BybitConfig struct {
	APIKey    string `json:"api_key" mapstructure:"api_key"`
	APISecret string `json:"api_secret" mapstructure:"api_secret"`
}

type SyncConfig struct {
	Symbol        string  `json:"symbol" mapstructure:"symbol"`
	PositionRatio float64 `json:"position_ratio" mapstructure:"position_ratio"`
	MaxPosition   float64 `json:"max_position" mapstructure:"max_position"`
	StopLossRatio float64 `json:"stop_loss_ratio" mapstructure:"stop_loss_ratio"`
	OrderTimeout  int     `json:"order_timeout" mapstructure:"order_timeout"`
	MaxRetries    int     `json:"max_retries" mapstructure:"max_retries"`
}

type Config struct {
	Auth          AuthConfig `json:"auth" mapstructure:"auth"`
	WebhookSecret string     `json:"webhook_secret" mapstructure:"webhook_secret"`
	SyncItems     []SyncItem `json:"sync_items" mapstructure:"sync_items"`

	Binance BinanceConfig `json:"binance" mapstructure:"binance"`
	OKX     OKXConfig     `json:"okx" mapstructure:"okx"`
	Bybit   BybitConfig   `json:"bybit" mapstructure:"bybit"`
	Sync    SyncConfig    `json:"sync" mapstructure:"sync"`

	mu sync.RWMutex `json:"-"`
}

func LoadConfig() (*Config, error) {
	// 1. Try to load from Database first
	if database.MySQLDB != nil {
		data, err := database.LoadConfigRaw()
		if err == nil && len(data) > 0 {
			var cfg Config
			if err := json.Unmarshal(data, &cfg); err == nil {
				log.Println("Loaded config from MySQL")
				return &cfg, nil
			}
		}
	}

	// 2. Fallback to Environment Variables (Viper)
	viper.AutomaticEnv()

	// Bind environment variables
	viper.BindEnv("binance.api_key", "BINANCE_API_KEY")
	viper.BindEnv("binance.api_secret", "BINANCE_API_SECRET")
	viper.BindEnv("binance.testnet", "BINANCE_TESTNET")
	viper.BindEnv("okx.api_key", "OKX_API_KEY")
	viper.BindEnv("okx.api_secret", "OKX_API_SECRET")
	viper.BindEnv("okx.passphrase", "OKX_API_PASSPHRASE")
	viper.BindEnv("bybit.api_key", "BYBIT_API_KEY")
	viper.BindEnv("bybit.api_secret", "BYBIT_API_SECRET")
	viper.BindEnv("sync.symbol", "SYMBOL")
	viper.BindEnv("sync.position_ratio", "POSITION_RATIO")
	viper.BindEnv("sync.max_position", "MAX_POSITION")
	viper.BindEnv("sync.stop_loss_ratio", "STOP_LOSS_RATIO")

	viper.SetDefault("binance.testnet", false)
	viper.SetDefault("sync.position_ratio", 1.0)
	viper.SetDefault("sync.max_position", 1.0)
	viper.SetDefault("sync.stop_loss_ratio", 0.05)
	viper.SetDefault("sync.order_timeout", 30)
	viper.SetDefault("sync.max_retries", 3)

	var cfg Config
	// Viper unmarshal from Env
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// 3. Save initialized config to DB for next time
	if database.MySQLDB != nil {
		if err := database.SaveConfig(&cfg); err != nil {
			log.Printf("Warning: Failed to save initial config to DB: %v", err)
		} else {
			log.Println("Initialized config in MySQL from Environment")
		}
	}

	if key := os.Getenv("ENCRYPTION_KEY"); key != "" {
		cfg.decryptFields(key)
	}

	return &cfg, nil
}

func (c *Config) decryptFields(key string) {
	if c.Binance.APIKey != "" {
		if decrypted, err := auth.Decrypt(c.Binance.APIKey, key); err == nil {
			c.Binance.APIKey = decrypted
		}
	}
	if c.Binance.APISecret != "" {
		if decrypted, err := auth.Decrypt(c.Binance.APISecret, key); err == nil {
			c.Binance.APISecret = decrypted
		}
	}
	if c.OKX.APIKey != "" {
		if decrypted, err := auth.Decrypt(c.OKX.APIKey, key); err == nil {
			c.OKX.APIKey = decrypted
		}
	}
	if c.OKX.APISecret != "" {
		if decrypted, err := auth.Decrypt(c.OKX.APISecret, key); err == nil {
			c.OKX.APISecret = decrypted
		}
	}
	if c.OKX.Passphrase != "" {
		if decrypted, err := auth.Decrypt(c.OKX.Passphrase, key); err == nil {
			c.OKX.Passphrase = decrypted
		}
	}
	if c.Bybit.APIKey != "" {
		if decrypted, err := auth.Decrypt(c.Bybit.APIKey, key); err == nil {
			c.Bybit.APIKey = decrypted
		}
	}
	if c.Bybit.APISecret != "" {
		if decrypted, err := auth.Decrypt(c.Bybit.APISecret, key); err == nil {
			c.Bybit.APISecret = decrypted
		}
	}
}

func (c *Config) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if database.MySQLDB != nil {
		return database.SaveConfig(c)
	}

	// Fallback to file if DB not available (e.g. local dev without DB)
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("config.json", data, 0644)
}

func (c *Config) GetAuth() AuthConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Auth
}

func (c *Config) GetWebhookSecret() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.WebhookSecret
}

func (c *Config) GetSyncItems() []SyncItem {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Return a copy to avoid race on slice elements if modified elsewhere
	items := make([]SyncItem, len(c.SyncItems))
	copy(items, c.SyncItems)
	return items
}

func (c *Config) GetBinance() BinanceConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Binance
}

func (c *Config) GetOKX() OKXConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.OKX
}

func (c *Config) GetBybit() BybitConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Bybit
}

func (c *Config) GetSync() SyncConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Sync
}

func (c *Config) Update(binance BinanceConfig, okx OKXConfig, bybit BybitConfig, sync SyncConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Binance = binance
	c.OKX = okx
	c.Bybit = bybit
	c.Sync = sync
}

func (c *Config) UpdateAuth(auth AuthConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Auth = auth
}

func (c *Config) SetSyncItems(items []SyncItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.SyncItems = items
}

func (c *Config) AddSyncItem(item SyncItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.SyncItems = append(c.SyncItems, item)
}

func (c *Config) DeleteSyncItem(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, item := range c.SyncItems {
		if item.ID == id {
			c.SyncItems = append(c.SyncItems[:i], c.SyncItems[i+1:]...)
			return true
		}
	}
	return false
}
