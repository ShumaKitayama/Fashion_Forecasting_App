package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Global Redis client
var RedisClient *redis.Client

// InitRedis initializes the Redis connection
func InitRedis() error {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// Create Redis client
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return nil
}

// CloseRedis closes the Redis connection
func CloseRedis() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		} else {
			log.Println("Redis connection closed")
		}
	}
}

// SetWithTTL sets a key-value pair with expiration time
func SetWithTTL(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func Get(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// Del deletes a key
func Del(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// KeyExists checks if a key exists
func KeyExists(ctx context.Context, key string) (bool, error) {
	result, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
} 