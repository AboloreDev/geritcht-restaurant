package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
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

func (c *NopCache) FlushByPattern(ctx context.Context, pattern string) error {
	return nil
}

func (c *NopCache) Hold(ctx context.Context, keys string, value interface{}, args redis.SetArgs) error {
	return nil
}

func (c *NopCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return nil
}

func (c *NopCache) Increment(ctx context.Context, key string) (int64, error) {
	return 0, nil
}
