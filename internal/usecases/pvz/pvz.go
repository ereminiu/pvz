package pvz

import (
	"context"

	"github.com/ereminiu/pvz/internal/entities"
	"github.com/ereminiu/pvz/internal/monitoring"
)

type Repository interface {
	GetRefunds(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error)
	GetHistory(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error)
}

type Usecases struct {
	r Repository
	c Cache
}

type Cache interface {
	GetRefunds(ctx context.Context) ([]*entities.Order, error)
	SetRefunds(ctx context.Context, refunds []*entities.Order) error
	RemoveRefunds(ctx context.Context) error
	GetHistory(ctx context.Context) ([]*entities.Order, error)
	SetHistory(ctx context.Context, history []*entities.Order) error
	RemoveHistory(ctx context.Context) error
}

func New(repos Repository, cache Cache) *Usecases {
	return &Usecases{
		r: repos,
		c: cache,
	}
}

func (uc *Usecases) GetRefunds(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error) {
	if pattern == nil {
		pattern = make(map[string]string, 0)
	}

	cachedOrders, err := uc.c.GetRefunds(ctx)
	if err == nil {
		monitoring.SetCacheCounter("getRefunds")

		return cachedOrders, nil
	}

	orders, err := uc.r.GetRefunds(ctx, page, limit, orderBy, pattern)
	if err != nil {
		return nil, err
	}

	if err = uc.c.SetRefunds(ctx, orders); err != nil {
		return nil, err
	}

	return orders, uc.c.RemoveRefunds(ctx)
}

func (uc *Usecases) GetHistory(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error) {
	if pattern == nil {
		pattern = make(map[string]string, 0)
	}

	cachedOrders, err := uc.c.GetHistory(ctx)
	if err == nil {
		monitoring.SetCacheCounter("getHistory")

		return cachedOrders, nil
	}

	orders, err := uc.r.GetHistory(ctx, page, limit, orderBy, pattern)
	if err != nil {
		return nil, err
	}

	if err = uc.c.SetHistory(ctx, orders); err != nil {
		return nil, err
	}

	return orders, uc.c.RemoveHistory(ctx)
}
