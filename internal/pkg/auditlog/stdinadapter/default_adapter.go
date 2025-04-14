package stdinadapter

import (
	"context"
	"log/slog"

	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
)

type StdinAdapter struct {
}

func New() *StdinAdapter {
	return &StdinAdapter{}
}

func (adapter *StdinAdapter) Process(ctx context.Context, event models.Log) error {
	slog.InfoContext(ctx, "event",
		slog.String("action", event.Action),
		slog.String("description", event.Description),
		slog.String("error", event.Error),
	)

	return nil
}
