package cache

import (
	"context"
	"time"
)

// provides an interface for creating cached functions
type CacheBuilder struct {
	prefix string
	ttl    time.Duration
}

func NewCache(prefix string) *CacheBuilder {
	return &CacheBuilder{prefix: prefix, ttl: 1 * time.Hour}
}

func (cb *CacheBuilder) TTL(duration time.Duration) *CacheBuilder {
	cb.ttl = duration
	return cb
}

func (cb *CacheBuilder) Wrap(fn func() (any, error)) func(context.Context, ...any) (any, error) {
	return CachedFunc(cb.prefix, cb.ttl, fn)
}
