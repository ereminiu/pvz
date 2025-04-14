package order

import (
	"context"
	"time"

	"github.com/ereminiu/pvz/internal/entities"
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	pck "github.com/ereminiu/pvz/internal/usecases/order/packing"
	"github.com/ereminiu/pvz/internal/usecases/status"
)

type Repository interface {
	AddOrder(ctx context.Context, order *entities.Order) error
	RemoveOrder(ctx context.Context, id int) (int, error)
}

type Cache interface {
	RemoveOrders(ctx context.Context, userID int) error
}

type OrderUsecases struct {
	r Repository
	c Cache
}

func New(repos Repository, cache Cache) *OrderUsecases {
	return &OrderUsecases{
		r: repos,
		c: cache,
	}
}

func (uc *OrderUsecases) AddOrder(ctx context.Context, order *entities.Order) error {
	packing, err := pck.New(order.Packing)
	if err != nil {
		return myerrors.ErrInvalidOrderPackingType
	}

	roller := &Roller{
		order: order,
	}

	if err = roller.Roll(packing); err != nil {
		return err
	}

	if order.Extra {
		packing, err = pck.New("film")
		if err != nil {
			return err
		}

		if err = roller.Roll(packing); err != nil {
			return err
		}
	}

	now := time.Now()
	order.ExpireAt = now.AddDate(0, 0, order.ExpireAfter)
	order.UpdatedAt = now
	order.Status = status.Delivered

	if err = uc.r.AddOrder(ctx, order); err != nil {
		return err
	}

	return uc.c.RemoveOrders(ctx, order.UserID)
}

func (uc *OrderUsecases) RemoveOrder(ctx context.Context, id int) error {
	userID, err := uc.r.RemoveOrder(ctx, id)
	if err != nil {
		return err
	}

	return uc.c.RemoveOrders(ctx, userID)
}
