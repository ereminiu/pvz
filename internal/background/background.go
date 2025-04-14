package background

import (
	"context"
	"log/slog"
	"time"

	"github.com/ereminiu/pvz/internal/entities"
)

type Usecases interface {
	GetRefunds(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error)
	GetHistory(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error)
}

type Cache interface {
	SetRefunds(ctx context.Context, refunds []*entities.Order) error
	SetHistory(ctx context.Context, history []*entities.Order) error
}

type Backgrounder struct {
	u Usecases
	c Cache
}

func New(usecases Usecases, cache Cache) *Backgrounder {
	return &Backgrounder{
		u: usecases,
		c: cache,
	}
}

func (b *Backgrounder) Run(ctx context.Context, timeout time.Duration, task func(ctx context.Context) error) {
	ticker := time.NewTicker(timeout)

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if err := task(ctx); err != nil {
				slog.Error("error during running background task", slog.Any("err", err))
				return
			}
		}
	}
}

func (b *Backgrounder) FillHistory(ctx context.Context) error {
	orders, err := b.u.GetHistory(ctx, -1, -1, "", nil)
	if err != nil {
		return err
	}

	return b.c.SetHistory(ctx, orders)
}

func (b *Backgrounder) FillRefunds(ctx context.Context) error {
	orders, err := b.u.GetRefunds(ctx, -1, -1, "", nil)
	if err != nil {
		return err
	}

	return b.c.SetRefunds(ctx, orders)
}
