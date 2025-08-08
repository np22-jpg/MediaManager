package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitCache initializes the Redis/Valkey client with the provided configuration
// and tests the connection to ensure cache service availability.
func InitCache(host string, port, db int) {
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
		DB:   db,
	})

	// Test the connection and error if not available
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		fmt.Printf("ERROR: Valkey/Redis is not available at %s:%d - %v\n", host, port, err)
		panic("Cache service unavailable")
	}

	// Get server info to display version and confirmation
	info, err := redisClient.Info(ctx, "server").Result()
	if err != nil {
		fmt.Printf("WARNING: Connected to cache at %s:%d but could not get server info: %v\n", host, port, err)
	} else {
		// Extract version and server type from info string
		version := "unknown"
		serverType := "Redis"

		lines := strings.Split(info, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "valkey_version:") {
				// Valkey reports its actual version here
				version = strings.TrimSpace(strings.Split(line, ":")[1])
				serverType = "Valkey"
			} else if strings.HasPrefix(line, "server_name:") && strings.Contains(line, "valkey") {
				// Confirm this is Valkey (NOTE: Valkey reports a redis version for compatibility reasons)
				serverType = "Valkey"
			} else if strings.HasPrefix(line, "redis_version:") && serverType == "Redis" {
				// Only use redis_version if we haven't found Valkey info
				version = strings.TrimSpace(strings.Split(line, ":")[1])
			}
		}

		fmt.Printf("Connected to %s v%s at %s:%d\n", serverType, version, host, port)
	}
}

// GenerateCacheKey creates a cache key from prefix and parameters using xxHash
// to ensure consistent and collision-resistant key generation.
func GenerateCacheKey(prefix string, params ...any) string {
	keyData := fmt.Sprintf("%s:%v", prefix, params)
	hash := xxhash.Sum64String(keyData)
	return fmt.Sprintf("%x", hash)
}

// GetCachedResponse retrieves cached data from Redis and deserializes it.
// Returns nil for cache misses or unmarshaling errors.
func GetCachedResponse(ctx context.Context, cacheKey string) (any, error) {
	// Return cache miss if Redis client is not initialized
	if redisClient == nil {
		return nil, nil
	}

	cachedData, err := redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		slog.Error("error getting cached response", "error", err)
		return nil, err
	}

	var result any
	if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
		slog.Error("error unmarshaling cached data", "error", err)
		return nil, err
	}

	return result, nil
}

// SetCachedResponse stores data in cache with the specified TTL.
// Serializes data to JSON before storage.
func SetCachedResponse(ctx context.Context, cacheKey string, data any, ttl time.Duration) error {
	// Skip caching if Redis client is not initialized
	if redisClient == nil {
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("error marshaling data for cache", "error", err)
		return err
	}

	if err := redisClient.Set(ctx, cacheKey, jsonData, ttl).Err(); err != nil {
		slog.Error("error setting cached response", "error", err)
		return err
	}

	return nil
}

// GetString retrieves a raw string value. Returns (value, true, nil) on hit; ("", false, nil) on miss.
func GetString(ctx context.Context, key string) (string, bool, error) {
	if redisClient == nil {
		return "", false, nil
	}
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", false, nil
		}
		return "", false, err
	}
	return val, true, nil
}

// SetString stores a raw string value with an optional TTL (0 for no expiry).
func SetString(ctx context.Context, key, value string, ttl time.Duration) error {
	if redisClient == nil {
		return nil
	}
	return redisClient.Set(ctx, key, value, ttl).Err()
}
