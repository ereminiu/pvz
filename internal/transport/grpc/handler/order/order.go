package order

import (
	"context"
	"log/slog"

	"github.com/ereminiu/pvz/internal/entities"
	"github.com/ereminiu/pvz/internal/pb/api"
	"github.com/ereminiu/pvz/internal/usecases"
)

type Server struct {
	api.UnimplementedOrderServer

	uc usecases.OrderUsecases
}

func New(uc usecases.OrderUsecases) *Server {
	return &Server{
		uc: uc,
	}
}

func (s *Server) Create(ctx context.Context, req *api.AddOrderRequest) (*api.AddOrderResponse, error) {
	order := toOrder(req)

	if err := s.uc.AddOrder(ctx, order); err != nil {
		slog.Error("error during adding order", slog.Any("err", err))

		return &api.AddOrderResponse{
			Message: err.Error(),
		}, err
	}

	return &api.AddOrderResponse{
		Message: "order created",
	}, nil
}

func (s *Server) Remove(ctx context.Context, req *api.RemoveOrderRequest) (*api.RemoveOrderResponse, error) {
	orderID := int(req.GetOrderId())

	if err := s.uc.RemoveOrder(ctx, orderID); err != nil {
		return &api.RemoveOrderResponse{
			Message: err.Error(),
		}, err
	}

	return &api.RemoveOrderResponse{
		Message: "order has been removed",
	}, nil
}

func toOrder(req *api.AddOrderRequest) *entities.Order {
	return &entities.Order{
		OrderID:     int(req.GetOrderId()),
		Price:       int(req.GetPrice()),
		UserID:      int(req.GetUserId()),
		ExpireAfter: int(req.GetExpireAfter()),
		Weight:      int(req.GetWeight()),
		Packing:     req.GetPacking(),
		Extra:       req.GetExtra(),
	}
}
