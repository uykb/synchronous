package processor

import (
	"context"
	"fmt"
	"time"

	"crypto-sync-bot/internal/database"
)

func IsDuplicate(ctx context.Context, signalID, exchange string, quantity, price float64) (bool, error) {
	key := fmt.Sprintf("signal:%s:%s:%.8f:%.8f", exchange, signalID, quantity, price)
	
	n, err := database.RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	
	return n > 0, nil
}

func MarkProcessed(ctx context.Context, signalID, exchange string, quantity, price float64) error {
	key := fmt.Sprintf("signal:%s:%s:%.8f:%.8f", exchange, signalID, quantity, price)
	return database.RDB.Set(ctx, key, "1", 24*time.Hour).Err()
}
