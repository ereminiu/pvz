package user

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ereminiu/pvz/internal/entities"
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	txmanager "github.com/ereminiu/pvz/internal/pkg/tx_manager"
	"github.com/ereminiu/pvz/internal/usecases/status"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
)

type Repository struct {
	manager *txmanager.TxManager
}

func New(manager *txmanager.TxManager) *Repository {
	return &Repository{
		manager: manager,
	}
}

func (r *Repository) RefundOrders(ctx context.Context, userID int, orderIDs []int) error {
	ok, err := r.checkOrders(ctx, userID, status.Given, orderIDs)
	if err != nil {
		return err
	}

	if !ok {
		return myerrors.ErrInvalidOrderInput
	}

	query := fmt.Sprintf(`UPDATE orders
		SET status='refund'
		WHERE user_id=$1 and id in (%s) and status='%s'`, r.prepareValues(2, orderIDs), status.Given)

	slog.InfoContext(ctx, "query: ", slog.String("query", query))

	args := make([]any, 0, len(orderIDs)+1)
	args = append(args, userID)
	for _, id := range orderIDs {
		args = append(args, id)
	}

	_, err = r.manager.GetQueryEngine(ctx).Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "error during updating")
	}

	return nil
}

func (r *Repository) ReturnOrders(ctx context.Context, userID int, orderIDs []int) error {
	ok, err := r.checkOrders(ctx, userID, status.Delivered, orderIDs)
	if err != nil {
		return err
	}

	if !ok {
		return myerrors.ErrInvalidOrderInput
	}

	query := fmt.Sprintf(`UPDATE orders 
		SET status='given'
		WHERE user_id=$1 AND id in (%s) AND status='%s'`, r.prepareValues(2, orderIDs), status.Delivered)

	args := r.prepareArgs(userID, orderIDs)
	_, err = r.manager.GetQueryEngine(ctx).Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "error during updating")
	}

	return nil
}

func (r *Repository) GetList(ctx context.Context, userID, lastN int, located bool, pattern map[string]string) ([]*entities.Order, error) {
	query := `SELECT id as "order_id", user_id, weight, price, packing, extra, status, expire_at
			FROM orders
			WHERE user_id=$1`

	args := make([]any, 0, 2)
	args = append(args, userID)
	cur := 2

	if len(pattern) > 0 {
		query = fmt.Sprintf("%s AND %s", query, r.preparePatternSearch(pattern))
	}

	if located {
		condition := fmt.Sprintf(`AND status=$%d`, cur)
		args = append(args, status.Delivered)
		query = fmt.Sprintf("%s %s", query, condition)
		cur++
	}

	if lastN != -1 {
		condition := fmt.Sprintf(`LIMIT $%d`, cur)
		args = append(args, lastN)
		query = fmt.Sprintf("%s %s", query, condition)
	}

	res := make([]*entities.Order, 0)
	if err := pgxscan.Select(ctx, r.manager.GetQueryEngine(ctx), &res, query, args...); err != nil {
		return nil, errors.Wrap(err, "error during selecting orders")
	}

	return res, nil
}

func (r *Repository) checkOrders(ctx context.Context, userID int, status string, orderIDs []int) (bool, error) {
	query := fmt.Sprintf(`SELECT COUNT(*)
			FROM orders 
			WHERE user_id=$1 AND status=$2 AND id IN (%s)`, r.prepareValues(3, orderIDs))

	slog.InfoContext(ctx, "query: ", slog.String("query", query))

	args := make([]any, 0, len(orderIDs)+2)
	args = append(args, userID)
	args = append(args, status)
	for _, id := range orderIDs {
		args = append(args, id)
	}

	countRow := make([]int, 1)
	err := pgxscan.Select(ctx, r.manager.GetQueryEngine(ctx), &countRow, query, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, myerrors.ErrInvalidOrderInput
	}
	count := countRow[0]

	return count == len(orderIDs), nil
}

func (r *Repository) prepareValues(startFrom int, values []int) string {
	valuesQuery := make([]string, 0, len(values))
	for i := range values {
		valuesQuery = append(valuesQuery, fmt.Sprintf("$%d", startFrom+i))
	}

	return strings.Join(valuesQuery, ",")
}

func (r *Repository) prepareArgs(userID int, orderIDs []int) []any {
	args := make([]any, 0, len(orderIDs)+1)
	args = append(args, userID)
	for _, id := range orderIDs {
		args = append(args, id)
	}

	return args
}

func (r *Repository) preparePatternSearch(pattern map[string]string) string {
	filters := make([]string, 0, len(pattern))
	for key, val := range pattern {
		filters = append(filters, fmt.Sprintf(` %s LIKE '%s' `, key, "%"+val+"%"))
	}

	return strings.Join(filters, " AND ")
}
