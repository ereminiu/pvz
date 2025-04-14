package pvz

import (
	"context"

	"github.com/ereminiu/pvz/internal/pb/api"
	"github.com/ereminiu/pvz/internal/usecases"
)

type Server struct {
	api.UnimplementedPVZServer

	uc usecases.PVZUsecases
}

func New(uc usecases.PVZUsecases) *Server {
	return &Server{
		uc: uc,
	}
}

func (s *Server) RefundList(ctx context.Context, req *api.RefundListRequest) (*api.RefundListResponse, error) {
	var (
		page    = int(req.Page)
		limit   = int(req.Limit)
		orderBy = req.OrderBy
		pattern = req.Pattern
	)

	orders, err := s.uc.GetRefunds(ctx, page, limit, orderBy, pattern)

	return &api.RefundListResponse{
		Orders: api.ToOrderList(orders),
	}, err
}

func (s *Server) HistoryList(ctx context.Context, req *api.HistoryListRequest) (*api.HistoryListResponse, error) {
	var (
		page    = int(req.Page)
		limit   = int(req.Limit)
		orderBy = req.OrderBy
		pattern = req.Pattern
	)

	orders, err := s.uc.GetHistory(ctx, page, limit, orderBy, pattern)

	return &api.HistoryListResponse{
		Orders: api.ToOrderList(orders),
	}, err
}
