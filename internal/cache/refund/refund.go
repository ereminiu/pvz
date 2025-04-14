package refund

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/ereminiu/pvz/internal/entities"
)

const (
	refundKey = "refunds"
)

type RedisCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any) error
	Remove(ctx context.Context, key string) error
}

type Cache struct {
	r RedisCache
}

func New(r RedisCache) *Cache {
	return &Cache{
		r: r,
	}
}

func (c *Cache) GetRefunds(ctx context.Context) ([]*entities.Order, error) {
	data, err := c.r.Get(ctx, refundKey)
	if err != nil {
		return nil, err
	}

	orders := make([]*entities.Order, 0)
	if err = sonic.Unmarshal([]byte(data), &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (cache *Cache) SetRefunds(ctx context.Context, refunds []*entities.Order) error {
	data, err := sonic.Marshal(refunds)
	if err != nil {
		return err
	}

	return cache.r.Set(ctx, refundKey, data)
}

func (c *Cache) RemoveRefunds(ctx context.Context) error {
	return c.r.Remove(
		ctx,
		refundKey,
	)
}
