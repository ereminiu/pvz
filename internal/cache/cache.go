package cache

import (
	"context"

	"github.com/ereminiu/pvz/internal/cache/history"
	"github.com/ereminiu/pvz/internal/cache/order"
	"github.com/ereminiu/pvz/internal/cache/refund"
	"github.com/ereminiu/pvz/internal/entities"
)

type RedisCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any) error
	Remove(ctx context.Context, key string) error
}

type OrderCache interface {
	GetOrders(ctx context.Context, userID int) ([]*entities.Order, error)
	SetOrders(ctx context.Context, userID int, orders []*entities.Order) error
	RemoveOrders(ctx context.Context, userID int) error
}

type RefundCache interface {
	GetRefunds(ctx context.Context) ([]*entities.Order, error)
	SetRefunds(ctx context.Context, refunds []*entities.Order) error
	RemoveRefunds(ctx context.Context) error
}

type HistoryCache interface {
	GetHistory(ctx context.Context) ([]*entities.Order, error)
	SetHistory(ctx context.Context, history []*entities.Order) error
	RemoveHistory(ctx context.Context) error
}

type Cache struct {
	OrderCache
	RefundCache
	HistoryCache
}

func New(redisCache RedisCache) *Cache {
	return &Cache{
		OrderCache:   order.New(redisCache),
		RefundCache:  refund.New(redisCache),
		HistoryCache: history.New(redisCache),
	}
}
