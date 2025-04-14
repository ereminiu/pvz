package user

import (
	"context"

	"github.com/ereminiu/pvz/internal/pb/api"
	"github.com/ereminiu/pvz/internal/usecases"
)

type Server struct {
	api.UnimplementedUserServer

	uc usecases.UserUsecases
}

func New(uc usecases.UserUsecases) *Server {
	return &Server{
		uc: uc,
	}
}

func (s *Server) Refund(ctx context.Context, req *api.RefundOrderRequest) (*api.RefundOrderResponse, error) {
	var (
		userID  = int(req.GetUserId())
		orderID = make([]int, 0, len(req.GetOrderId()))
	)

	for _, id := range req.GetOrderId() {
		orderID = append(orderID, int(id))
	}

	if err := s.uc.RefundOrders(ctx, userID, orderID); err != nil {
		return &api.RefundOrderResponse{
			Message: err.Error(),
		}, err
	}

	return &api.RefundOrderResponse{
		Message: "orders have been refund",
	}, nil
}

func (s *Server) Return(ctx context.Context, req *api.ReturnOrderRequest) (*api.ReturnOrderResponse, error) {
	var (
		userID  = int(req.GetUserId())
		orderID = make([]int, 0, len(req.GetOrderId()))
	)

	for _, id := range req.GetOrderId() {
		orderID = append(orderID, int(id))
	}

	if err := s.uc.ReturnOrders(ctx, userID, orderID); err != nil {
		return &api.ReturnOrderResponse{
			Message: err.Error(),
		}, err
	}

	return &api.ReturnOrderResponse{
		Message: "orders have been refund",
	}, nil
}

func (s *Server) List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {
	var (
		userID  = int(req.GetUserId())
		lastN   = int(req.GetLastN())
		located = req.GetLocated()
		pattern = req.GetPattern()
	)

	orders, err := s.uc.GetList(ctx, userID, lastN, located, pattern)

	return &api.ListResponse{
		Orders: api.ToOrderList(orders),
	}, err
}
