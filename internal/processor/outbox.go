package processor

import (
	"crypto-sync-bot/internal/models"
	"fmt"
	"time"
)

// SaveToOutbox saves a message to the outbox table (mock implementation)
func SaveToOutbox(payload string) error {
	outbox := &models.Outbox{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Payload:   payload,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	fmt.Printf("Mock: Saving to outbox: %+v\n", outbox)
	return nil
}
