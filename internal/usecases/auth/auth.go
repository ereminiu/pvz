package auth

import (
	"context"

	"github.com/ereminiu/pvz/internal/pkg/role_token"
	"github.com/ereminiu/pvz/internal/usecases/auth/hashgen"
)

type Repository interface {
	CreateUser(ctx context.Context, username, password string) error
}

type usecases struct {
	r Repository
}

func New(repos Repository) *usecases {
	return &usecases{
		r: repos,
	}
}

func (uc *usecases) SignIn(ctx context.Context, username, password string) (string, error) {
	hash := hashgen.GenerateHash(password)

	if err := uc.r.CreateUser(ctx, username, hash); err != nil {
		return "", err
	}

	return role_token.GenerateToken()
}
