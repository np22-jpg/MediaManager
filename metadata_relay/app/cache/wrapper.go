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

// CachedFunc wraps a function with caching functionality, providing automatic
// cache lookup, metrics recording, and structured logging for API calls.
func CachedFunc[T any](prefix string, ttl time.Duration, fn func() (T, error)) func(context.Context, ...any) (T, error) {
	return func(ctx context.Context, params ...any) (T, error) {
		var zero T
		cacheKey := GenerateCacheKey(prefix, params...)

		// Extract meaningful information from prefix and params for logging
		operation := extractOperation(prefix)
		provider := extractProvider(prefix)
		paramDetails := formatParams(params)

		// Try to retrieve from cache first
		start := time.Now()
		if cached, err := GetCachedResponse(ctx, cacheKey); err == nil && cached != nil {
			duration := time.Since(start)

			// Record cache hit metrics
			metrics.RecordCacheHit(operation, provider, duration)

			// Get cache TTL information
			ttlRemaining := getCacheRemainingTTL(ctx, cacheKey)

			// Try to extract content information from cached data
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
			// If type assertion fails, fall through to actual API call
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

		// Call the actual function and measure API response time
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

		// Extract information about what we're caching for logging
		resultInfo := extractContentInfo(result)

		// Cache the result for future requests
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

			// Record successful cache store operation
			metrics.RecordCacheStore(operation, provider, cacheDuration)

			// Record content metrics for monitoring purposes
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

// extractOperation extracts the operation type from the cache prefix.
// Example: "tmdb:search" -> "search"
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

// extractProvider extracts the provider name from the cache prefix and normalizes it.
// Example: "TMDB:search" -> "tmdb"
func extractProvider(prefix string) string {
	if prefix == "" {
		return "unknown"
	}

	// Get the first part before the colon
	parts := strings.Split(prefix, ":")
	if len(parts) > 0 {
		provider := parts[0]
		// Normalize common provider names to lowercase
		switch provider {
		case "tmdb", "TMDB":
			return "tmdb"
		case "tvdb", "TVDB":
			return "tvdb"
		case "musicbrainz", "MusicBrainz", "MUSICBRAINZ":
			return "musicbrainz"
		default:
			return provider
		}
	}
	return "unknown"
}

// formatParams creates a human-readable string representation of function parameters
// for logging purposes.
func formatParams(params []any) string {
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

// extractContentInfo extracts meaningful information from cached content for logging.
// Attempts to parse common API response fields like titles, names, and result counts.
func extractContentInfo(data any) string {
	// Try to convert to JSON and extract relevant fields
	if jsonData, err := json.Marshal(data); err == nil {
		var parsed map[string]any
		if err := json.Unmarshal(jsonData, &parsed); err == nil {
			var info []string

			// Check for common fields in TMDB/TVDB API responses
			if results, ok := parsed["results"].([]any); ok {
				info = append(info, fmt.Sprintf("results=%d", len(results)))

				// Try to get title/name from the first result for context
				if len(results) > 0 {
					if first, ok := results[0].(map[string]any); ok {
						if title, ok := first["title"].(string); ok && title != "" {
							info = append(info, fmt.Sprintf("first_title=%q", title))
						} else if name, ok := first["name"].(string); ok && name != "" {
							info = append(info, fmt.Sprintf("first_name=%q", name))
						}
					}
				}
			} else {
				// Handle single item responses
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

// extractContentMetrics extracts content count and type information for metrics recording.
// Analyzes API responses to determine content type (movie, TV show) and item count.
func extractContentMetrics(data any, operation string) (count int, contentType string) {
	contentType = "unknown"
	count = 0

	// Try to convert to JSON and extract metrics information
	if jsonData, err := json.Marshal(data); err == nil {
		var parsed map[string]any
		if err := json.Unmarshal(jsonData, &parsed); err == nil {
			// Check for results array (typical in search/trending responses)
			if results, ok := parsed["results"].([]any); ok {
				count = len(results)

				// Determine content type from operation name or first result
				if strings.Contains(operation, "movie") || strings.Contains(operation, "trending") {
					contentType = "movie"
				} else if strings.Contains(operation, "tv") || strings.Contains(operation, "series") {
					contentType = "tv"
				} else if len(results) > 0 {
					// Try to determine content type from the first result structure
					if first, ok := results[0].(map[string]any); ok {
						if _, hasTitle := first["title"]; hasTitle {
							contentType = "movie"
						} else if _, hasName := first["name"]; hasName {
							contentType = "tv"
						}
					}
				}
			} else {
				// Handle single item responses
				count = 1
				if strings.Contains(operation, "movie") {
					contentType = "movie"
				} else if strings.Contains(operation, "tv") || strings.Contains(operation, "series") {
					contentType = "tv"
				} else {
					// Try to determine content type from response structure
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

// getCacheRemainingTTL retrieves the remaining time-to-live for a cache key.
// Returns 0 if Redis client is not initialized or if TTL retrieval fails.
func getCacheRemainingTTL(ctx context.Context, cacheKey string) time.Duration {
	// Return 0 if Redis client is not initialized
	if redisClient == nil {
		return 0
	}

	if ttl, err := redisClient.TTL(ctx, cacheKey).Result(); err == nil {
		return ttl
	}
	return 0
}
