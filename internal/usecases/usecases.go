package usecases

import (
	"context"

	"github.com/ereminiu/pvz/internal/entities"
	"github.com/ereminiu/pvz/internal/usecases/auth"
	"github.com/ereminiu/pvz/internal/usecases/order"
	"github.com/ereminiu/pvz/internal/usecases/pvz"
	"github.com/ereminiu/pvz/internal/usecases/user"
)

//go:generate mockgen -source=usecases.go -destination=mocks/mock.go

type OrderRepository interface {
	AddOrder(ctx context.Context, order *entities.Order) error
	RemoveOrder(ctx context.Context, id int) (int, error)
}

type UserRepository interface {
	RefundOrders(ctx context.Context, userID int, orderIDs []int) error
	ReturnOrders(ctx context.Context, userID int, orderIDs []int) error
	GetList(ctx context.Context, userID, lastN int, located bool, pattern map[string]string) ([]*entities.Order, error)
}

type PVZRepository interface {
	GetRefunds(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error)
	GetHistory(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error)
}

type AuthRepository interface {
	CreateUser(ctx context.Context, username, password string) error
}

type Repository interface {
	OrderRepository
	UserRepository
	PVZRepository
	AuthRepository
}

type RedisCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any) error
	Remove(ctx context.Context, key string) error
}

type OrderCache interface {
	GetOrders(ctx context.Context, userID int) ([]*entities.Order, error)
	SetOrders(ctx context.Context, userID int, orders []*entities.Order) error
	RemoveOrders(ctx context.Context, userID int) error
}

type RefundCache interface {
	GetRefunds(ctx context.Context) ([]*entities.Order, error)
	SetRefunds(ctx context.Context, refunds []*entities.Order) error
	RemoveRefunds(ctx context.Context) error
}

type HistoryCache interface {
	GetHistory(ctx context.Context) ([]*entities.Order, error)
	SetHistory(ctx context.Context, history []*entities.Order) error
	RemoveHistory(ctx context.Context) error
}

type Cache interface {
	OrderCache
	RefundCache
	HistoryCache
}

type OrderUsecases interface {
	AddOrder(ctx context.Context, order *entities.Order) error
	RemoveOrder(ctx context.Context, id int) error
}

type UserUsecases interface {
	RefundOrders(ctx context.Context, userID int, orderIDs []int) error
	ReturnOrders(ctx context.Context, userID int, orderIDs []int) error
	GetList(ctx context.Context, userID, lastN int, located bool, pattern map[string]string) ([]*entities.Order, error)
}

type PVZUsecases interface {
	GetRefunds(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error)
	GetHistory(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error)
}

type AuthUsecases interface {
	SignIn(ctx context.Context, username, password string) (string, error)
}

type Usecases struct {
	OrderUsecases
	UserUsecases
	PVZUsecases
	AuthUsecases
}

func New(repos Repository, cache Cache) *Usecases {
	return &Usecases{
		OrderUsecases: order.New(repos, cache),
		UserUsecases:  user.New(repos, cache),
		PVZUsecases:   pvz.New(repos, cache),
		AuthUsecases:  auth.New(repos),
	}
}
