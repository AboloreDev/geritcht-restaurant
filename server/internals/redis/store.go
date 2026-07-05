package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{
		client: client,
	}
}

func (r *Store) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", nil
	}

	return val, err
}

func (r *Store) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Store) Delete(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Store) Exists(ctx context.Context, keys string) (bool, error) {
	count, err := r.client.Exists(ctx, keys).Result()
	if err != nil {
		return count > 0, err
	}
	return false, nil
}

func (r *Store) Flush(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

func (r *Store) FlushByPattern(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}

	return r.client.Del(ctx, keys...).Err()
}

func (r *Store) Hold(ctx context.Context, keys string, value interface{}, args redis.SetArgs) error {
	err := r.client.SetArgs(ctx, keys, value, args).Err()

	if err != nil {
		return err
	}

	return nil
}

func (r *Store) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}

func (r *Store) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}
