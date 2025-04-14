package auth

import (
	"context"

	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	txmanager "github.com/ereminiu/pvz/internal/pkg/tx_manager"
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

func (r *Repository) CreateUser(ctx context.Context, username, password string) error {
	query := `INSERT INTO admins (username, password) 
		VALUES ($1, $2)
		ON CONFLICT(username) DO UPDATE
		SET username=EXCLUDED.username
		RETURNING password`

	var actualPassword string
	if err := pgxscan.Get(ctx, r.manager.GetQueryEngine(ctx), &actualPassword, query, username, password); err != nil {
		return errors.Wrap(err, "error during inserting")
	}

	if password != actualPassword {
		return myerrors.ErrIncorrectPassword
	}

	return nil
}
