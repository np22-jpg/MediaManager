package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log/slog"
	"strings"
	"time"

	"github.com/valkey-io/valkey-go"
)

var valkeyClient valkey.Client

// InitCache initializes the Redis/Valkey client with the provided configuration
// and tests the connection to ensure cache service availability.
func InitCache(host string, port, db int) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{fmt.Sprintf("%s:%d", host, port)},
		SelectDB:    db,
	})
	if err != nil {
		fmt.Printf("ERROR: Failed to create Valkey client: %v\n", err)
		panic("Cache service unavailable")
	}
	valkeyClient = client

	// Test the connection and error if not available
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := valkeyClient.Do(ctx, valkeyClient.B().Ping().Build()).Error(); err != nil {
		fmt.Printf("ERROR: Valkey/Redis is not available at %s:%d - %v\n", host, port, err)
		panic("Cache service unavailable")
	}

	// Get server info to display version and confirmation
	infoResp := valkeyClient.Do(ctx, valkeyClient.B().Info().Section("server").Build())
	if infoResp.Error() != nil {
		fmt.Printf("WARNING: Connected to cache at %s:%d but could not get server info: %v\n", host, port, infoResp.Error())
	} else {
		info, _ := infoResp.ToString()
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

// GenerateCacheKey creates a cache key from prefix and parameters using FNV-1a
// to ensure consistent and collision-resistant key generation.
func GenerateCacheKey(prefix string, params ...any) string {
	keyData := fmt.Sprintf("%s:%v", prefix, params)
	h := fnv.New64a()
	h.Write([]byte(keyData))
	hash := h.Sum64()
	return fmt.Sprintf("%x", hash)
}

// GetCachedResponse retrieves cached data from Redis and deserializes it.
// Returns nil for cache misses or unmarshaling errors.
func GetCachedResponse(ctx context.Context, cacheKey string) (any, error) {
	// Return cache miss if Valkey client is not initialized
	if valkeyClient == nil {
		return nil, nil
	}

	resp := valkeyClient.Do(ctx, valkeyClient.B().Get().Key(cacheKey).Build())
	if resp.Error() != nil {
		if valkey.IsValkeyNil(resp.Error()) {
			return nil, nil // Cache miss
		}
		slog.Error("error getting cached response", "error", resp.Error())
		return nil, resp.Error()
	}

	cachedData, err := resp.ToString()
	if err != nil {
		slog.Error("error converting cached data to string", "error", err)
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
	// Skip caching if Valkey client is not initialized
	if valkeyClient == nil {
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("error marshaling data for cache", "error", err)
		return err
	}

	var cmd valkey.Completed
	if ttl > 0 {
		cmd = valkeyClient.B().Setex().Key(cacheKey).Seconds(int64(ttl.Seconds())).Value(string(jsonData)).Build()
	} else {
		cmd = valkeyClient.B().Set().Key(cacheKey).Value(string(jsonData)).Build()
	}

	if err := valkeyClient.Do(ctx, cmd).Error(); err != nil {
		slog.Error("error setting cached response", "error", err)
		return err
	}

	return nil
}

// GetString retrieves a raw string value. Returns (value, true, nil) on hit; ("", false, nil) on miss.
func GetString(ctx context.Context, key string) (string, bool, error) {
	if valkeyClient == nil {
		return "", false, nil
	}
	resp := valkeyClient.Do(ctx, valkeyClient.B().Get().Key(key).Build())
	if resp.Error() != nil {
		if valkey.IsValkeyNil(resp.Error()) {
			return "", false, nil
		}
		return "", false, resp.Error()
	}
	val, err := resp.ToString()
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

// SetString stores a raw string value with an optional TTL (0 for no expiry).
func SetString(ctx context.Context, key, value string, ttl time.Duration) error {
	if valkeyClient == nil {
		return nil
	}
	var cmd valkey.Completed
	if ttl > 0 {
		cmd = valkeyClient.B().Setex().Key(key).Seconds(int64(ttl.Seconds())).Value(value).Build()
	} else {
		cmd = valkeyClient.B().Set().Key(key).Value(value).Build()
	}
	return valkeyClient.Do(ctx, cmd).Error()
}
