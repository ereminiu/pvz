package order

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/ereminiu/pvz/internal/entities"
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

func (cache *Cache) GetOrders(ctx context.Context, userID int) ([]*entities.Order, error) {
	data, err := cache.r.Get(ctx, fmt.Sprintf("%d", userID))
	if err != nil {
		return nil, err
	}

	orders := make([]*entities.Order, 0)
	if err = sonic.Unmarshal([]byte(data), &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (cache *Cache) SetOrders(ctx context.Context, userID int, orders []*entities.Order) error {
	data, err := sonic.Marshal(orders)
	if err != nil {
		return err
	}

	return cache.r.Set(ctx, fmt.Sprintf("%d", userID), data)
}

func (cache *Cache) RemoveOrders(ctx context.Context, userID int) error {
	return cache.r.Remove(
		ctx,
		fmt.Sprintf("%d", userID),
	)
}
