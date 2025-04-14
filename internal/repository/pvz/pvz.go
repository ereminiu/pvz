package pvz

import (
	"context"
	"fmt"
	"strings"

	"github.com/ereminiu/pvz/internal/entities"
	txmanager "github.com/ereminiu/pvz/internal/pkg/tx_manager"
	"github.com/ereminiu/pvz/internal/usecases/status"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
)

const perPage = 5

type Repository struct {
	manager *txmanager.TxManager
}

func New(manager *txmanager.TxManager) *Repository {
	return &Repository{
		manager: manager,
	}
}

func (r *Repository) GetRefunds(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error) {
	query := `SELECT id as "order_id", user_id, expire_at, weight, price, packing, status, extra, updated_at, created_at 
		FROM orders
		WHERE status=$1`

	if len(pattern) > 0 {
		query = fmt.Sprintf("%s AND %s", query, r.preparePatternSearch(pattern))
	}

	args := make([]any, 0, 3)
	args = append(args, status.Refund)
	pos := 2

	if page > 0 {
		condition, _ := r.preparePaging(orderBy, pos)
		query = fmt.Sprintf("%s %s", query, condition)

		args = append(args, (page-1)*perPage)
		if orderBy != "" {
			args = append(args, orderBy)
		}
		args = append(args, perPage)
	}

	res := make([]*entities.Order, 0)
	if err := pgxscan.Select(ctx, r.manager.GetQueryEngine(ctx), &res, query, args...); err != nil {
		return nil, errors.Wrap(err, "error during selecting refund")
	}

	return res, nil
}

func (r *Repository) GetHistory(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error) {
	query := `SELECT id as "order_id", user_id, expire_at, weight, price, packing, status, extra, updated_at, created_at 
		FROM orders
		WHERE status=$1`

	if len(pattern) > 0 {
		query = fmt.Sprintf("%s AND %s", query, r.preparePatternSearch(pattern))
	}

	args := make([]any, 0, 3)
	args = append(args, status.Delivered)
	pos := 2

	if page > 0 {
		condition, _ := r.preparePaging(orderBy, pos)
		query = fmt.Sprintf("%s %s", query, condition)

		args = append(args, (page-1)*perPage)
		if orderBy != "" {
			args = append(args, orderBy)
		}
		args = append(args, perPage)
	}

	res := make([]*entities.Order, 0)
	if err := pgxscan.Select(ctx, r.manager.GetQueryEngine(ctx), &res, query, args...); err != nil {
		return nil, errors.Wrap(err, "error during selecting history")
	}

	return res, nil
}

func (r *Repository) preparePaging(field string, pos int) (string, int) {
	query := fmt.Sprintf(`AND created_moment > $%d`, pos)
	pos++
	if field != "" {
		condition := fmt.Sprintf(`ORDER BY $%d`, pos)
		pos++
		query = fmt.Sprintf("%s %s", query, condition)
	}
	limitCondition := fmt.Sprintf(`LIMIT $%d`, pos)
	pos++
	query = fmt.Sprintf("%s %s", query, limitCondition)

	return query, pos
}

func (r *Repository) preparePatternSearch(pattern map[string]string) string {
	filters := make([]string, 0, len(pattern))
	for key, val := range pattern {
		filters = append(filters, fmt.Sprintf(` %s LIKE '%s' `, key, "%"+val+"%"))
	}

	return strings.Join(filters, " AND ")
}
