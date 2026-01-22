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
	PasswordHash string `json:"password_hash" mapstructure:"password_hash"`
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

type Config struct {
	Auth          AuthConfig `json:"auth" mapstructure:"auth"`
	WebhookSecret string     `json:"webhook_secret" mapstructure:"webhook_secret"`
	SyncItems     []SyncItem `json:"sync_items" mapstructure:"sync_items"`

	Binance struct {
		APIKey    string `json:"api_key" mapstructure:"api_key"`
		APISecret string `json:"api_secret" mapstructure:"api_secret"`
		Testnet   bool   `json:"testnet" mapstructure:"testnet"`
	} `json:"binance" mapstructure:"binance"`

	OKX struct {
		APIKey     string `json:"api_key" mapstructure:"api_key"`
		APISecret  string `json:"api_secret" mapstructure:"api_secret"`
		Passphrase string `json:"passphrase" mapstructure:"passphrase"`
	} `json:"okx" mapstructure:"okx"`

	Bybit struct {
		APIKey    string `json:"api_key" mapstructure:"api_key"`
		APISecret string `json:"api_secret" mapstructure:"api_secret"`
	} `json:"bybit" mapstructure:"bybit"`

	Sync struct {
		Symbol        string  `json:"symbol" mapstructure:"symbol"`
		PositionRatio float64 `json:"position_ratio" mapstructure:"position_ratio"`
		MaxPosition   float64 `json:"max_position" mapstructure:"max_position"`
		StopLossRatio float64 `json:"stop_loss_ratio" mapstructure:"stop_loss_ratio"`
		OrderTimeout  int     `json:"order_timeout" mapstructure:"order_timeout"`
		MaxRetries    int     `json:"max_retries" mapstructure:"max_retries"`
	} `json:"sync" mapstructure:"sync"`

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
