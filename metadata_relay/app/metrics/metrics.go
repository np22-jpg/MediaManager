package metrics

import (
	"fmt"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

// Helper function to format metric names with labels
func formatMetricName(base string, labels map[string]string) string {
	if len(labels) == 0 {
		return base
	}

	result := base + "{"
	first := true
	for k, v := range labels {
		if !first {
			result += ","
		}
		result += fmt.Sprintf("%s=%q", k, v)
		first = false
	}
	result += "}"
	return result
}

// Cache metrics for monitoring cache performance
func RecordCacheHit(operation, provider string, duration time.Duration) {
	metrics.GetOrCreateCounter(formatMetricName("metadata_relay_cache_hits_total", map[string]string{
		"operation": operation,
		"provider":  provider,
	})).Inc()

	metrics.GetOrCreateHistogram(formatMetricName("metadata_relay_cache_retrieval_duration_seconds", map[string]string{
		"operation": operation,
		"provider":  provider,
		"status":    "hit",
	})).Update(duration.Seconds())
}

func RecordCacheMiss(operation, provider string, duration time.Duration) {
	metrics.GetOrCreateCounter(formatMetricName("metadata_relay_cache_misses_total", map[string]string{
		"operation": operation,
		"provider":  provider,
	})).Inc()

	metrics.GetOrCreateHistogram(formatMetricName("metadata_relay_cache_retrieval_duration_seconds", map[string]string{
		"operation": operation,
		"provider":  provider,
		"status":    "miss",
	})).Update(duration.Seconds())
}

func RecordCacheStore(operation, provider string, duration time.Duration) {
	metrics.GetOrCreateHistogram(formatMetricName("metadata_relay_cache_store_duration_seconds", map[string]string{
		"operation": operation,
		"provider":  provider,
	})).Update(duration.Seconds())
}

func SetCacheEntries(provider string, count int) {
	metrics.GetOrCreateGauge(formatMetricName("metadata_relay_cache_entries", map[string]string{
		"provider": provider,
	}), func() float64 { return float64(count) })
}

// API metrics for external provider request monitoring
func RecordAPIRequest(provider, endpoint, status string, duration time.Duration) {
	metrics.GetOrCreateCounter(formatMetricName("metadata_relay_api_requests_total", map[string]string{
		"provider": provider,
		"endpoint": endpoint,
		"status":   status,
	})).Inc()

	metrics.GetOrCreateHistogram(formatMetricName("metadata_relay_api_request_duration_seconds", map[string]string{
		"provider": provider,
		"endpoint": endpoint,
		"status":   status,
	})).Update(duration.Seconds())
}

func RecordAPIError(provider, endpoint, errorType string) {
	metrics.GetOrCreateCounter(formatMetricName("metadata_relay_api_errors_total", map[string]string{
		"provider":   provider,
		"endpoint":   endpoint,
		"error_type": errorType,
	})).Inc()
}

// Authentication metrics for tracking auth attempts (e.g., TVDB token refresh)
func RecordAuthAttempt(provider, status string) {
	metrics.GetOrCreateCounter(formatMetricName("metadata_relay_auth_attempts_total", map[string]string{
		"provider": provider,
		"status":   status,
	})).Inc()
}

func UpdateAuthTokenExpiry(provider string, expiry time.Time) {
	metrics.GetOrCreateGauge(formatMetricName("metadata_relay_auth_token_expiry_timestamp", map[string]string{
		"provider": provider,
	}), func() float64 { return float64(expiry.Unix()) })
}

// HTTP metrics for monitoring incoming requests to the relay service
func RecordHTTPRequest(method, endpoint, status string, duration time.Duration) {
	metrics.GetOrCreateCounter(formatMetricName("metadata_relay_http_requests_total", map[string]string{
		"method":   method,
		"endpoint": endpoint,
		"status":   status,
	})).Inc()

	metrics.GetOrCreateHistogram(formatMetricName("metadata_relay_http_request_duration_seconds", map[string]string{
		"method":   method,
		"endpoint": endpoint,
		"status":   status,
	})).Update(duration.Seconds())
}

// Content metrics for tracking the volume of data returned
func RecordContentItems(provider, contentType string, count int) {
	metrics.GetOrCreateHistogram(formatMetricName("metadata_relay_content_items_returned", map[string]string{
		"provider":     provider,
		"content_type": contentType,
	})).Update(float64(count))
}
