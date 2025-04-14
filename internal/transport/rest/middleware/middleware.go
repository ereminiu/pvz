package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ereminiu/pvz/internal/pkg/role_token"
)

type userRole struct{}

const (
	admin         = "admin"
	authorization = "Authorization"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorization)
		if header == "" {
			http.Error(
				w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized,
			)
			slog.Error("error during auth", slog.Any("err", http.StatusText(http.StatusUnauthorized)))
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 {
			http.Error(
				w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized,
			)
			return
		}

		role, err := role_token.ParseRole(parts[1])
		if err != nil || role != admin {
			http.Error(
				w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized,
			)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, userRole{}, role)

		next(w, r.WithContext(ctx))
	}
}

func Log(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("accept request: ",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		next(w, r)
	}
}

func Chain(f ...func(next http.HandlerFunc) http.HandlerFunc) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(f) - 1; i >= 0; i-- {
			next = f[i](next)
		}

		return next
	}
}
