package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const ttl = 5 * time.Minute

type Cache struct {
	client *redis.Client
}

func New(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Set(ctx context.Context, key string, value any) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Remove(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
