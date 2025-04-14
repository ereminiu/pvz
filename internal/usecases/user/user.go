package user

import (
	"context"

	"github.com/ereminiu/pvz/internal/entities"
	"github.com/ereminiu/pvz/internal/monitoring"
)

type Repository interface {
	RefundOrders(ctx context.Context, userID int, orderIDs []int) error
	ReturnOrders(ctx context.Context, userID int, orderIDs []int) error
	GetList(ctx context.Context, userID, lastN int, located bool, pattern map[string]string) ([]*entities.Order, error)
}

type Cache interface {
	GetOrders(ctx context.Context, userID int) ([]*entities.Order, error)
	SetOrders(ctx context.Context, userID int, orders []*entities.Order) error
	RemoveOrders(ctx context.Context, userID int) error
	RemoveHistory(ctx context.Context) error
	RemoveRefunds(ctx context.Context) error
}

type Usecases struct {
	r Repository
	c Cache
}

func New(repository Repository, cache Cache) *Usecases {
	return &Usecases{
		r: repository,
		c: cache,
	}
}

func (uc *Usecases) RefundOrders(ctx context.Context, userID int, orderIDs []int) error {
	if err := uc.r.RefundOrders(ctx, userID, orderIDs); err != nil {
		return err
	}

	if err := uc.c.RemoveOrders(ctx, userID); err != nil {
		return err
	}

	if err := uc.c.RemoveHistory(ctx); err != nil {
		return err
	}

	return uc.c.RemoveRefunds(ctx)
}

func (uc *Usecases) ReturnOrders(ctx context.Context, userID int, orderIDs []int) error {
	if err := uc.r.ReturnOrders(ctx, userID, orderIDs); err != nil {
		return err
	}

	if err := uc.c.RemoveOrders(ctx, userID); err != nil {
		return err
	}

	if err := uc.c.RemoveHistory(ctx); err != nil {
		return err
	}

	return uc.c.RemoveRefunds(ctx)
}

func (uc *Usecases) GetList(ctx context.Context, userID, lastN int, located bool, pattern map[string]string) ([]*entities.Order, error) {
	if pattern == nil {
		pattern = make(map[string]string)
	}

	cachedOrders, err := uc.c.GetOrders(ctx, userID)
	if err == nil {
		monitoring.SetCacheCounter("getList")

		return cachedOrders, nil
	}

	orders, err := uc.r.GetList(ctx, userID, lastN, located, pattern)
	if err != nil {
		return nil, err
	}

	if err := uc.c.SetOrders(ctx, userID, orders); err != nil {
		return nil, err
	}

	return orders, err
}
