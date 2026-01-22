package processor

import (
	"context"
	"fmt"
	"time"

	"crypto-sync-bot/internal/database"
)

func IsDuplicate(ctx context.Context, signalID, exchange string, quantity, price float64) (bool, error) {
	key := fmt.Sprintf("signal:%s:%s:%.8f:%.8f", exchange, signalID, quantity, price)
	
	// Try to set the key with an expiration of 24 hours. 
	// If it already exists, it's a duplicate.
	success, err := database.RDB.SetNX(ctx, key, "1", 24*time.Hour).Result()
	if err != nil {
		return false, err
	}
	
	return !success, nil
}
