package consumer

import (
	"context"
	"log/slog"

	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
)

type StdinHandler struct {
}

func (h *StdinHandler) Process(ctx context.Context, task models.Task) error {
	slog.InfoContext(ctx, "log", slog.Any("value", task))

	return nil
}
