package redis

import (
	"context"
	"time"
)

type NopCache struct{}

func NewNopCache() *NopCache {
	return &NopCache{}
}

func (c *NopCache) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (c *NopCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (c *NopCache) Delete(ctx context.Context, keys ...string) error {
	return nil
}

func (c *NopCache) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (c *NopCache) Flush(ctx context.Context) error {
	return nil
}
