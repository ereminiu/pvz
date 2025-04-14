package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	cluster *pgxpool.Pool
}

func New(ctx context.Context, connString string) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &Database{
		cluster: pool,
	}, nil
}

func (conn *Database) GetPool() *pgxpool.Pool {
	return conn.cluster
}

func (conn *Database) Get(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Get(ctx, conn.cluster, dest, query, args...)
}

func (conn *Database) Select(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Select(ctx, conn.cluster, dest, query, args...)
}

func (conn *Database) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return conn.cluster.Exec(ctx, query, args...)
}

func (conn *Database) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return conn.cluster.QueryRow(ctx, query, args...)
}
