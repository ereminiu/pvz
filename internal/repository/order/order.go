package order

import (
	"context"

	"github.com/ereminiu/pvz/internal/entities"
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	txmanager "github.com/ereminiu/pvz/internal/pkg/tx_manager"
	"github.com/georgysavva/scany/pgxscan"
)

type Repository struct {
	manager *txmanager.TxManager
}

func New(manager *txmanager.TxManager) *Repository {
	return &Repository{
		manager: manager,
	}
}

func (r *Repository) AddOrder(ctx context.Context, order *entities.Order) error {
	query := `INSERT INTO 
			orders (user_id, id, expire_at, weight, price, packing, extra, status)
			values ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.manager.GetQueryEngine(ctx).Exec(ctx, query,
		order.UserID,
		order.OrderID,
		order.ExpireAt,
		order.Weight,
		order.Price,
		order.Packing,
		order.Extra,
		order.Status,
	)

	if err != nil {
		return myerrors.ErrOrderAlreadyCreated
	}

	return nil
}

func (r *Repository) RemoveOrder(ctx context.Context, id int) (int, error) {
	query := `DELETE FROM orders
			WHERE id=$1
			RETURNING user_id`

	var userID int

	if err := pgxscan.Get(ctx, r.manager.GetQueryEngine(ctx), &userID, query, id); err != nil {
		return -1, myerrors.ErrOrderAlreadyRemoved
	}

	return userID, nil
}
