package txmanager

import (
	"context"

	"github.com/ereminiu/pvz/internal/pkg/db"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type txManagerKey struct{}

type QueryEngine interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type TxManager struct {
	db *db.Database
}

func New(db *db.Database) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (m *TxManager) RunSerializable(ctx context.Context, fn func(ctx context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	return m.beginFunc(ctx, opts, fn)
}

func (m *TxManager) ReadUncommitted(ctx context.Context, fn func(ctx context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadUncommitted,
		AccessMode: pgx.ReadOnly,
	}

	return m.beginFunc(ctx, opts, fn)
}

func (m *TxManager) beginFunc(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context) error) error {
	tx, err := m.db.GetPool().BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ctx = context.WithValue(ctx, txManagerKey{}, tx)
	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (m *TxManager) GetQueryEngine(ctx context.Context) QueryEngine {
	v, ok := ctx.Value(txManagerKey{}).(QueryEngine)
	if ok && v != nil {
		return v
	}

	return m.db.GetPool()
}
