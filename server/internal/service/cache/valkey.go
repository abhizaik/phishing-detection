package cache

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New() (*Cache, error) {
	// Get configuration from environment variables with defaults
	addr := getEnv("CACHE_ADDR", "safesurf-valkey-dev:6379")
	password := getEnv("CACHE_PASSWORD", "")
	db := getEnvAsInt("CACHE_DB", 0)
	poolSize := getEnvAsInt("CACHE_POOL_SIZE", 50)
	minIdleConns := getEnvAsInt("CACHE_MIN_IDLE_CONNS", 10)

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	})

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Cache{client: rdb}, nil
}

func (c *Cache) Close() error {
	return c.client.Close()
}

func (c *Cache) Set(
	ctx context.Context,
	key string,
	value any,
	ttl time.Duration,
) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // key does not exist
	}
	return val, err
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// FlushAll removes all keys from the current database
func (c *Cache) FlushAll(ctx context.Context) error {
	return c.client.FlushAll(ctx).Err()
}

// Increment increments a key and returns the new value
func (c *Cache) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// Expire sets an expiration on a key
func (c *Cache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, key, ttl).Err()
}

// GetJSON retrieves a JSON value from cache and unmarshals it into the provided pointer
func (c *Cache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	if val == "" {
		return redis.Nil // Key doesn't exist
	}
	return json.Unmarshal([]byte(val), dest)
}

// SetJSON marshals the value to JSON and stores it in cache
func (c *Cache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, string(data), ttl)
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
