package processor

import (
	"context"
	"crypto-sync-bot/internal/database"
	"crypto-sync-bot/internal/models"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ProduceSignal(ctx context.Context, signal *models.TradingSignal) error {
	data, err := json.Marshal(signal)
	if err != nil {
		return fmt.Errorf("failed to marshal signal: %w", err)
	}

	err = database.RDB.XAdd(ctx, &redis.XAddArgs{
		Stream: "signals:trading",
		Values: map[string]interface{}{
			"payload": string(data),
			"retry_count": 0,
		},
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to add signal to stream: %w", err)
	}

	return nil
}
