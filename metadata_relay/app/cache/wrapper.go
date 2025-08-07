package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"
)

// CachedFunc wraps a function with caching functionality
func CachedFunc[T any](prefix string, ttl time.Duration, fn func() (T, error)) func(context.Context, ...interface{}) (T, error) {
	return func(ctx context.Context, params ...interface{}) (T, error) {
		var zero T
		cacheKey := GenerateCacheKey(prefix, params...)

		// Extract meaningful info from prefix and params for logging
		operation := extractOperation(prefix)
		paramDetails := formatParams(params)

		// Try cache first
		start := time.Now()
		if cached, err := GetCachedResponse(ctx, cacheKey); err == nil && cached != nil {
			duration := time.Since(start)

			// Get cache TTL info
			ttlRemaining := getCacheRemainingTTL(ctx, cacheKey)

			// Try to extract content info from cached data
			contentInfo := extractContentInfo(cached)

			slog.Info("cache hit",
				"operation", operation,
				"params", paramDetails,
				"key", cacheKey,
				"ttl_remaining", ttlRemaining,
				"retrieval_time", duration,
				"content", contentInfo,
			)

			if result, ok := cached.(T); ok {
				return result, nil
			}
			// If type assertion fails, fall through to actual call
			slog.Warn("cache type mismatch, falling back to API call",
				"operation", operation,
				"key", cacheKey,
				"expected_type", reflect.TypeOf(zero).String(),
				"actual_type", reflect.TypeOf(cached).String(),
			)
		}

		cacheMissStart := time.Now()
		slog.Info("cache miss",
			"operation", operation,
			"params", paramDetails,
			"key", cacheKey,
			"ttl", ttl,
			"cache_check_time", time.Since(start),
		)

		// Call the actual function
		apiStart := time.Now()
		result, err := fn()
		if err != nil {
			apiDuration := time.Since(apiStart)
			slog.Error("API call failed",
				"operation", operation,
				"params", paramDetails,
				"error", err,
				"api_call_time", apiDuration,
			)
			return zero, err
		}

		apiDuration := time.Since(apiStart)

		// Extract info about what we're caching
		resultInfo := extractContentInfo(result)

		// Cache the result
		cacheStart := time.Now()
		if err := SetCachedResponse(ctx, cacheKey, result, ttl); err != nil {
			slog.Error("failed to cache response",
				"operation", operation,
				"params", paramDetails,
				"key", cacheKey,
				"error", err,
				"api_call_time", apiDuration,
			)
		} else {
			cacheDuration := time.Since(cacheStart)
			totalDuration := time.Since(cacheMissStart)

			slog.Info("cached new response",
				"operation", operation,
				"params", paramDetails,
				"key", cacheKey,
				"ttl", ttl,
				"content", resultInfo,
				"api_call_time", apiDuration,
				"cache_store_time", cacheDuration,
				"total_time", totalDuration,
			)
		}

		return result, nil
	}
}

// extractOperation converts cache prefix to human-readable operation
func extractOperation(prefix string) string {
	parts := strings.Split(prefix, "_")
	if len(parts) >= 2 {
		provider := strings.ToUpper(parts[0]) // tmdb, tvdb
		category := parts[1]                  // tv, movies
		action := ""
		if len(parts) > 2 {
			action = parts[2] // trending, search, show, etc.
		}
		return fmt.Sprintf("%s %s %s", provider, category, action)
	}
	return prefix
}

// formatParams creates a readable string from parameters
func formatParams(params []interface{}) string {
	if len(params) == 0 {
		return "none"
	}

	var parts []string
	for i, param := range params {
		switch v := param.(type) {
		case string:
			if v != "" {
				parts = append(parts, fmt.Sprintf("query=%q", v))
			}
		case int:
			if i == 0 {
				parts = append(parts, fmt.Sprintf("id=%d", v))
			} else {
				parts = append(parts, fmt.Sprintf("param%d=%d", i+1, v))
			}
		default:
			parts = append(parts, fmt.Sprintf("param%d=%v", i+1, v))
		}
	}
	return strings.Join(parts, ", ")
}

// extractContentInfo tries to extract meaningful info from cached content
func extractContentInfo(data interface{}) string {
	// Try to convert to JSON and extract relevant fields
	if jsonData, err := json.Marshal(data); err == nil {
		var parsed map[string]interface{}
		if err := json.Unmarshal(jsonData, &parsed); err == nil {
			var info []string

			// Check for common fields in TMDB/TVDB responses
			if results, ok := parsed["results"].([]interface{}); ok {
				info = append(info, fmt.Sprintf("results=%d", len(results)))

				// Try to get title/name from first result
				if len(results) > 0 {
					if first, ok := results[0].(map[string]interface{}); ok {
						if title, ok := first["title"].(string); ok && title != "" {
							info = append(info, fmt.Sprintf("first_title=%q", title))
						} else if name, ok := first["name"].(string); ok && name != "" {
							info = append(info, fmt.Sprintf("first_name=%q", name))
						}
					}
				}
			} else {
				// Single item response
				if title, ok := parsed["title"].(string); ok && title != "" {
					info = append(info, fmt.Sprintf("title=%q", title))
				} else if name, ok := parsed["name"].(string); ok && name != "" {
					info = append(info, fmt.Sprintf("name=%q", name))
				}

				if id, ok := parsed["id"].(float64); ok {
					info = append(info, fmt.Sprintf("id=%.0f", id))
				}
			}

			if len(info) > 0 {
				return strings.Join(info, ", ")
			}
		}
	}

	return fmt.Sprintf("type=%T", data)
}

// getCacheRemainingTTL gets the remaining TTL for a cache key
func getCacheRemainingTTL(ctx context.Context, cacheKey string) time.Duration {
	if ttl, err := redisClient.TTL(ctx, cacheKey).Result(); err == nil {
		return ttl
	}
	return 0
}

// WithCache is a simpler wrapper that takes the function and returns a cached version
func WithCache[T any](prefix string, ttl time.Duration) func(func() (T, error)) func(context.Context, ...interface{}) (T, error) {
	return func(fn func() (T, error)) func(context.Context, ...interface{}) (T, error) {
		return CachedFunc(prefix, ttl, fn)
	}
}
