package repository

import (
	"context"

	"github.com/ereminiu/pvz/internal/entities"
	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
	txmanager "github.com/ereminiu/pvz/internal/pkg/tx_manager"
	"github.com/ereminiu/pvz/internal/repository/audit"
	"github.com/ereminiu/pvz/internal/repository/auth"
	"github.com/ereminiu/pvz/internal/repository/order"
	"github.com/ereminiu/pvz/internal/repository/pvz"
	"github.com/ereminiu/pvz/internal/repository/user"
)

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

type AuditRepository interface {
	AddTask(ctx context.Context, task models.Task) error
	GetTasks(ctx context.Context) ([]*models.Task, error)
}

type Repository struct {
	OrderRepository
	UserRepository
	PVZRepository
	AuthRepository
	AuditRepository
}

func New(manager *txmanager.TxManager) *Repository {
	return &Repository{
		OrderRepository: order.New(manager),
		UserRepository:  user.New(manager),
		PVZRepository:   pvz.New(manager),
		AuthRepository:  auth.New(manager),
		AuditRepository: audit.New(manager),
	}
}
