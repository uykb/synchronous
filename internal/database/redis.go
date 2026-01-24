package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() error {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		// Redis not configured - this is optional
		return nil
	}

	password := os.Getenv("REDIS_PASSWORD")

	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	err := RDB.Ping(context.Background()).Err()
	if err != nil {
		RDB = nil // Clear client on failure
		return err
	}
	return nil
}
