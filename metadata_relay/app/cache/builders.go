package cache

import (
	"context"
	"time"
)

// CacheBuilder provides a fluent interface for creating cached functions
// with customizable TTL settings.
type CacheBuilder struct {
	prefix string
	ttl    time.Duration
}

// NewCache creates a new cache builder with the specified prefix and default 1-hour TTL.
func NewCache(prefix string) *CacheBuilder {
	return &CacheBuilder{prefix: prefix, ttl: 1 * time.Hour}
}

// TTL sets the time-to-live duration for cached data and returns the builder for chaining.
func (cb *CacheBuilder) TTL(duration time.Duration) *CacheBuilder {
	cb.ttl = duration
	return cb
}

// Wrap wraps a function with caching functionality using the configured settings.
func (cb *CacheBuilder) Wrap(fn func() (any, error)) func(context.Context, ...any) (any, error) {
	return CachedFunc(cb.prefix, cb.ttl, fn)
}
