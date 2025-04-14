package grpc

import (
	"fmt"
	"net"

	"github.com/ereminiu/pvz/internal/pb/api"
	mid "github.com/ereminiu/pvz/internal/transport/grpc/handler/middleware"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

type Deps struct {
	OrderHandler api.OrderServer
	UserHandler  api.UserServer
	PVZHandler   api.PVZServer
}

type Server struct {
	Deps
	srv *grpc.Server
}

func New(deps Deps) *Server {
	return &Server{
		srv: grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				mid.Auth,
				mid.Logging,
				mid.ErrorRate,
				mid.RPS,
				mid.Tracing,
				mid.Duration,
			)),
		),
		Deps: deps,
	}
}

func (s *Server) ListenAndServe(port int) error {
	lis, err := net.Listen("tcp",
		fmt.Sprintf(":%d", port),
	)
	if err != nil {
		return err
	}

	api.RegisterOrderServer(s.srv, s.Deps.OrderHandler)
	api.RegisterUserServer(s.srv, s.Deps.UserHandler)
	api.RegisterPVZServer(s.srv, s.Deps.PVZHandler)

	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
