package cache

import (
	"context"
	"time"
)

// CacheBuilder provides a fluent interface for creating cached functions
type CacheBuilder struct {
	prefix string
	ttl    time.Duration
}

func NewCache(prefix string) *CacheBuilder {
	return &CacheBuilder{prefix: prefix, ttl: 1 * time.Hour} // default TTL
}

func (cb *CacheBuilder) TTL(duration time.Duration) *CacheBuilder {
	cb.ttl = duration
	return cb
}

func (cb *CacheBuilder) Wrap(fn func() (interface{}, error)) func(context.Context, ...interface{}) (interface{}, error) {
	return CachedFunc(cb.prefix, cb.ttl, fn)
}
