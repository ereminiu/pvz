package middleware

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/ereminiu/pvz/internal/monitoring"
	"github.com/ereminiu/pvz/internal/pkg/role_token"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type userRole struct{}

func Logging(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	slog.InfoContext(ctx, "grpc call", slog.String("method", info.FullMethod))

	monitoring.SetTotalRequestCounter()

	res, err := handler(ctx, req)
	if err != nil {
		monitoring.SetTotalErrorCounter()
	} else {
		monitoring.SetTotalOkCounter()
	}

	return res, err
}

func Auth(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	if len(meta["authorization"]) != 1 {
		return nil, status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	bearer := meta["authorization"][0]
	parts := strings.Split(bearer, " ")
	if len(parts) != 2 {
		return nil, status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	role, err := role_token.ParseRole(parts[1])
	if err != nil || !role_token.CheckRole(role) {
		return nil, status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	ctx = context.WithValue(ctx, userRole{}, role)

	return handler(ctx, req)
}

func ErrorRate(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	res, err := handler(ctx, req)

	if err != nil {
		monitoring.SetErrorCounter(info.FullMethod)
	} else {
		monitoring.SetOkCounter(info.FullMethod)
	}

	return res, err
}

func RPS(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	monitoring.SetRPSCounter(info.FullMethod)

	return handler(ctx, req)
}

func Tracing(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	defer span.Finish()

	startTime := time.Now()
	res, err := handler(ctx, req)

	if err != nil {
		span.SetTag("error", err)
	}

	span.SetTag("duration_ms", time.Since(startTime))

	return res, err
}

func Duration(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	now := time.Now()

	res, err := handler(ctx, req)

	monitoring.SetRequestDuration(info.FullMethod, float64(time.Since(now)))
	monitoring.SetResponseTimeSummary(float64(time.Since(now)))

	return res, err
}
