package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"relay/app/metrics"
)

// CachedFunc wraps a function with caching functionality
func CachedFunc[T any](prefix string, ttl time.Duration, fn func() (T, error)) func(context.Context, ...interface{}) (T, error) {
	return func(ctx context.Context, params ...interface{}) (T, error) {
		var zero T
		cacheKey := GenerateCacheKey(prefix, params...)

		// Extract meaningful info from prefix and params for logging
		operation := extractOperation(prefix)
		provider := extractProvider(prefix)
		paramDetails := formatParams(params)

		// Try cache first
		start := time.Now()
		if cached, err := GetCachedResponse(ctx, cacheKey); err == nil && cached != nil {
			duration := time.Since(start)

			// Record cache hit metrics
			metrics.RecordCacheHit(operation, provider, duration)

			// Get cache TTL info
			ttlRemaining := getCacheRemainingTTL(ctx, cacheKey)

			// Try to extract content info from cached data
			contentInfo := extractContentInfo(cached)

			slog.Info("cache hit",
				"operation", operation,
				"provider", provider,
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
				"provider", provider,
				"key", cacheKey,
				"expected_type", reflect.TypeOf(zero).String(),
				"actual_type", reflect.TypeOf(cached).String(),
			)
		}

		cacheCheckDuration := time.Since(start)

		// Record cache miss metrics
		metrics.RecordCacheMiss(operation, provider, cacheCheckDuration)

		slog.Info("cache miss",
			"operation", operation,
			"provider", provider,
			"params", paramDetails,
			"key", cacheKey,
			"ttl", ttl,
			"cache_check_time", cacheCheckDuration,
		)

		// Call the actual function
		apiStart := time.Now()
		result, err := fn()
		apiDuration := time.Since(apiStart)
		totalDuration := time.Since(start)

		if err != nil {
			// Record API error metrics
			metrics.RecordAPIError(provider, operation, "error")

			slog.Error("API call failed",
				"operation", operation,
				"provider", provider,
				"params", paramDetails,
				"error", err,
				"api_time", apiDuration,
				"total_time", totalDuration,
			)
			return zero, err
		}

		// Record successful API call metrics
		metrics.RecordAPIRequest(provider, operation, "success", apiDuration)

		// Extract info about what we're caching
		resultInfo := extractContentInfo(result)

		// Cache the result
		cacheStart := time.Now()
		if err := SetCachedResponse(ctx, cacheKey, result, ttl); err != nil {
			slog.Error("failed to cache response",
				"operation", operation,
				"provider", provider,
				"params", paramDetails,
				"key", cacheKey,
				"error", err,
				"api_time", apiDuration,
				"cache_time", time.Since(cacheStart),
				"total_time", totalDuration,
			)
		} else {
			cacheDuration := time.Since(cacheStart)

			// Record successful cache store
			metrics.RecordCacheStore(operation, provider, cacheDuration)

			// Record content metrics
			count, contentType := extractContentMetrics(result, operation)
			if count > 0 {
				metrics.RecordContentItems(provider, contentType, count)
			}

			slog.Info("cached new response",
				"operation", operation,
				"provider", provider,
				"params", paramDetails,
				"key", cacheKey,
				"ttl", ttl,
				"content", resultInfo,
				"content_count", count,
				"content_type", contentType,
				"api_time", apiDuration,
				"cache_time", cacheDuration,
				"total_time", totalDuration,
			)
		}

		return result, nil
	}
}

// extractOperation extracts the operation type from the cache prefix
func extractOperation(prefix string) string {
	if prefix == "" {
		return "unknown"
	}

	// Split by colon and look for operation patterns
	parts := strings.Split(prefix, ":")
	if len(parts) >= 2 {
		return parts[1] // e.g., "tmdb:search" -> "search"
	}
	return prefix
}

// extractProvider extracts the provider from the cache prefix
func extractProvider(prefix string) string {
	if prefix == "" {
		return "unknown"
	}

	// Get the first part before the colon
	parts := strings.Split(prefix, ":")
	if len(parts) > 0 {
		provider := parts[0]
		// Normalize common provider names
		switch provider {
		case "tmdb", "TMDB":
			return "tmdb"
		case "tvdb", "TVDB":
			return "tvdb"
		default:
			return provider
		}
	}
	return "unknown"
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

// extractContentMetrics extracts content count and type for metrics
func extractContentMetrics(data interface{}, operation string) (count int, contentType string) {
	contentType = "unknown"
	count = 0

	// Try to convert to JSON and extract metrics info
	if jsonData, err := json.Marshal(data); err == nil {
		var parsed map[string]interface{}
		if err := json.Unmarshal(jsonData, &parsed); err == nil {
			// Check for results array (search/trending responses)
			if results, ok := parsed["results"].([]interface{}); ok {
				count = len(results)

				// Determine content type from operation or first result
				if strings.Contains(operation, "movie") || strings.Contains(operation, "trending") {
					contentType = "movie"
				} else if strings.Contains(operation, "tv") || strings.Contains(operation, "series") {
					contentType = "tv"
				} else if len(results) > 0 {
					// Try to determine from first result
					if first, ok := results[0].(map[string]interface{}); ok {
						if _, hasTitle := first["title"]; hasTitle {
							contentType = "movie"
						} else if _, hasName := first["name"]; hasName {
							contentType = "tv"
						}
					}
				}
			} else {
				// Single item response
				count = 1
				if strings.Contains(operation, "movie") {
					contentType = "movie"
				} else if strings.Contains(operation, "tv") || strings.Contains(operation, "series") {
					contentType = "tv"
				} else {
					// Try to determine from content
					if _, hasTitle := parsed["title"]; hasTitle {
						contentType = "movie"
					} else if _, hasName := parsed["name"]; hasName {
						contentType = "tv"
					}
				}
			}
		}
	}

	return count, contentType
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
