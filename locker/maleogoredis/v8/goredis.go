package maleogoredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/tigorlazuardi/maleo/locker"
	"time"
)

// New create new redis locker.
func New(client *redis.Client) locker.Locker {
	return &goredis{client: client}
}

type goredis struct {
	client *redis.Client
}

// Set the Cache key and value.
func (g *goredis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return g.client.Set(ctx, key, value, ttl).Err()
}

// Get the Value by Key. Returns tower.ErrNilCache if not found or ttl has passed.
func (g *goredis) Get(ctx context.Context, key string) ([]byte, error) {
	v, err := g.client.Get(ctx, key).Result()
	return []byte(v), err
}

// Delete cache by key.
func (g *goredis) Delete(ctx context.Context, key string) {
	g.client.Del(ctx, key)
}

// Exist Checks if Key exist in cache.
func (g *goredis) Exist(ctx context.Context, key string) bool {
	return g.client.Exists(ctx, key).Val() > 0
}

// Separator Returns Accepted separator value for the Cacher implementor.
func (g *goredis) Separator() string {
	return ":"
}
