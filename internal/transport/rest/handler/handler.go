package handler

import (
	"context"
	"net/http"

	"github.com/ereminiu/pvz/internal/entities"
	"github.com/ereminiu/pvz/internal/pkg/auditlog"
	"github.com/ereminiu/pvz/internal/transport/rest/handler/auth"
	"github.com/ereminiu/pvz/internal/transport/rest/handler/order"
	"github.com/ereminiu/pvz/internal/transport/rest/handler/pvz"
	"github.com/ereminiu/pvz/internal/transport/rest/handler/user"
)

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

type Usecases interface {
	OrderUsecases
	UserUsecases
	PVZUsecases
	AuthUsecases
}

type OrderHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Remove(w http.ResponseWriter, r *http.Request)
}

type UserHandler interface {
	RefundOrders(w http.ResponseWriter, r *http.Request)
	ReturnOrders(w http.ResponseWriter, r *http.Request)
	GetList(w http.ResponseWriter, r *http.Request)
}

type PVZHandler interface {
	GetRefunds(w http.ResponseWriter, r *http.Request)
	GetHistory(w http.ResponseWriter, r *http.Request)
}

type AuthHandler interface {
	SignIn(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	OrderHandler
	UserHandler
	PVZHandler
	AuthHandler
}

func New(ctx context.Context, usecases Usecases, audit *auditlog.AuditLog) *Handler {
	return &Handler{
		OrderHandler: order.New(usecases, audit),
		UserHandler:  user.New(usecases, audit),
		PVZHandler:   pvz.New(usecases, audit),
		AuthHandler:  auth.New(usecases, audit),
	}
}
