package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Cache metrics
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "metadata_relay_cache_hits_total",
			Help: "Total number of cache hits by operation and provider",
		},
		[]string{"operation", "provider"},
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "metadata_relay_cache_misses_total",
			Help: "Total number of cache misses by operation and provider",
		},
		[]string{"operation", "provider"},
	)

	CacheEntries = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "metadata_relay_cache_entries",
			Help: "Current number of entries in cache by provider",
		},
		[]string{"provider"},
	)

	CacheRetrievalDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "metadata_relay_cache_retrieval_duration_seconds",
			Help:    "Time spent retrieving data from cache",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "provider", "status"},
	)

	CacheStoreDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "metadata_relay_cache_store_duration_seconds",
			Help:    "Time spent storing data to cache",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "provider"},
	)

	// API metrics
	APIRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "metadata_relay_api_request_duration_seconds",
			Help:    "Time spent making API requests to external providers",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"provider", "endpoint", "status"},
	)

	APIRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "metadata_relay_api_requests_total",
			Help: "Total number of API requests to external providers",
		},
		[]string{"provider", "endpoint", "status"},
	)

	APIErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "metadata_relay_api_errors_total",
			Help: "Total number of API errors by provider and error type",
		},
		[]string{"provider", "endpoint", "error_type"},
	)

	// Authentication metrics (for TVDB)
	AuthenticationAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "metadata_relay_auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"provider", "status"},
	)

	AuthTokenExpiry = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "metadata_relay_auth_token_expiry_timestamp",
			Help: "Timestamp when the authentication token expires",
		},
		[]string{"provider"},
	)

	// HTTP metrics
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "metadata_relay_http_request_duration_seconds",
			Help:    "Time spent processing HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)

	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "metadata_relay_http_requests_total",
			Help: "Total number of HTTP requests received",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Content metrics
	ContentItemsReturned = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "metadata_relay_content_items_returned",
			Help:    "Number of content items returned in responses",
			Buckets: []float64{0, 1, 5, 10, 20, 50, 100, 200, 500},
		},
		[]string{"provider", "content_type"},
	)
)

// records a cache hit with timing information
func RecordCacheHit(operation, provider string, duration time.Duration) {
	CacheHits.WithLabelValues(operation, provider).Inc()
	CacheRetrievalDuration.WithLabelValues(operation, provider, "hit").Observe(duration.Seconds())
}

// records a cache miss with timing information
func RecordCacheMiss(operation, provider string, duration time.Duration) {
	CacheMisses.WithLabelValues(operation, provider).Inc()
	CacheRetrievalDuration.WithLabelValues(operation, provider, "miss").Observe(duration.Seconds())
}

// records cache storage timing
func RecordCacheStore(operation, provider string, duration time.Duration) {
	CacheStoreDuration.WithLabelValues(operation, provider).Observe(duration.Seconds())
}

// records an API request with timing and status
func RecordAPIRequest(provider, endpoint, status string, duration time.Duration) {
	APIRequestsTotal.WithLabelValues(provider, endpoint, status).Inc()
	APIRequestDuration.WithLabelValues(provider, endpoint, status).Observe(duration.Seconds())
}

// records an API error
func RecordAPIError(provider, endpoint, errorType string) {
	APIErrors.WithLabelValues(provider, endpoint, errorType).Inc()
}

// records an authentication attempt
func RecordAuthAttempt(provider, status string) {
	AuthenticationAttempts.WithLabelValues(provider, status).Inc()
}

// updates the token expiry timestamp
func UpdateAuthTokenExpiry(provider string, expiry time.Time) {
	AuthTokenExpiry.WithLabelValues(provider).Set(float64(expiry.Unix()))
}

// records an HTTP request to the relay
func RecordHTTPRequest(method, endpoint, status string, duration time.Duration) {
	HTTPRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	HTTPRequestDuration.WithLabelValues(method, endpoint, status).Observe(duration.Seconds())
}

// records the number of content items returned
func RecordContentItems(provider, contentType string, count int) {
	ContentItemsReturned.WithLabelValues(provider, contentType).Observe(float64(count))
}
