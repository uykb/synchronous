package models

import "time"

type Outbox struct {
	ID        string    `json:"id"`
	Payload   string    `json:"payload"`
	Status    string    `json:"status"` // "pending", "processed", "failed"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
