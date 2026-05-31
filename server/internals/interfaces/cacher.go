package interfaces

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cacher interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Flush(ctx context.Context) error
	FlushByPattern(ctx context.Context, pattern string) error
	Hold(ctx context.Context, keys string, value interface{}, args redis.SetArgs) error
}
