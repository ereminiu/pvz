package logger

import (
	"errors"
	"log/slog"
	"os"
)

const (
	local = "local"
	prod  = "prod"
)

func NewLogger(env string) (*slog.Logger, error) {
	var handler slog.Handler

	switch env {
	case local:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})

	case prod:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})

	default:
		return nil, errors.New("invalid environment")
	}

	return slog.New(handler), nil
}
