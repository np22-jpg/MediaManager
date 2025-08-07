package cache

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host string `env:"VALKEY_HOST" envDefault:"localhost"`
	Port int    `env:"VALKEY_PORT" envDefault:"6379"`
	DB   int    `env:"VALKEY_DB" envDefault:"0"`
}

var redisClient *redis.Client

func init() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		slog.Error("failed to parse cache configuration", "error", err)
		// Use defaults if parsing fails
		cfg = Config{
			Host: "localhost",
			Port: 6379,
			DB:   0,
		}
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:   cfg.DB,
	})

	// Test the connection and error if not available
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		fmt.Printf("ERROR: Valkey/Redis is not available at %s:%d - %v\n", cfg.Host, cfg.Port, err)
		panic("Cache service unavailable")
	}

	// Get server info to display version and confirmation
	info, err := redisClient.Info(ctx, "server").Result()
	if err != nil {
		fmt.Printf("WARNING: Connected to cache at %s:%d but could not get server info: %v\n", cfg.Host, cfg.Port, err)
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
				// Confirm this is Valkey
				serverType = "Valkey"
			} else if strings.HasPrefix(line, "redis_version:") && serverType == "Redis" {
				// Only use redis_version if we haven't found Valkey info
				version = strings.TrimSpace(strings.Split(line, ":")[1])
			}
		}

		fmt.Printf("Connected to %s v%s at %s:%d\n", serverType, version, cfg.Host, cfg.Port)
	}
}

// GenerateCacheKey creates a cache key from prefix and parameters
func GenerateCacheKey(prefix string, params ...interface{}) string {
	keyData := fmt.Sprintf("%s:%v", prefix, params)
	hash := md5.Sum([]byte(keyData))
	return fmt.Sprintf("%x", hash)
}

// GetCachedResponse retrieves cached data
func GetCachedResponse(ctx context.Context, cacheKey string) (interface{}, error) {
	cachedData, err := redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		slog.Error("error getting cached response", "error", err)
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal([]byte(cachedData), &result); err != nil {
		slog.Error("error unmarshaling cached data", "error", err)
		return nil, err
	}

	return result, nil
}

// SetCachedResponse stores data in cache
func SetCachedResponse(ctx context.Context, cacheKey string, data interface{}, ttl time.Duration) error {
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
