package history

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/ereminiu/pvz/internal/entities"
)

const (
	historyKey = "history"
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

func (c *Cache) GetHistory(ctx context.Context) ([]*entities.Order, error) {
	data, err := c.r.Get(ctx, historyKey)
	if err != nil {
		return nil, err
	}

	orders := make([]*entities.Order, 0)
	if err = sonic.Unmarshal([]byte(data), &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (cache *Cache) SetHistory(ctx context.Context, history []*entities.Order) error {
	data, err := sonic.Marshal(history)
	if err != nil {
		return err
	}

	return cache.r.Set(ctx, historyKey, data)
}

func (c *Cache) RemoveHistory(ctx context.Context) error {
	return c.r.Remove(
		ctx,
		historyKey,
	)
}
