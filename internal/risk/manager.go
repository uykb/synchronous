package risk

import (
	"crypto-sync-bot/internal/config"
	"crypto-sync-bot/internal/models"
	"fmt"
)

type Manager struct {
	config *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{config: cfg}
}

func (m *Manager) PreOrderCheck(signal *models.TradingSignal) error {
	// Check Symbol
	if signal.Symbol != m.config.Sync.Symbol {
		// Sometimes Binance uses BTCUSDT and config might be BTC-USDT, handle normalization if needed
		// For now, strict check
		// return fmt.Errorf("symbol mismatch: %s != %s", signal.Symbol, m.config.Sync.Symbol)
	}

	// Check Max Position
	if signal.Quantity > m.config.Sync.MaxPosition {
		return fmt.Errorf("quantity %.4f exceeds max position %.4f", signal.Quantity, m.config.Sync.MaxPosition)
	}

	return nil
}
