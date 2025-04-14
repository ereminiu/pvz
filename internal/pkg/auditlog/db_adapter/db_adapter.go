package dbadapter

import (
	"context"
	"time"

	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
	"github.com/ereminiu/pvz/internal/pkg/outbox/status"
)

const attemptsAmount = 3

type repository interface {
	AddTask(ctx context.Context, task models.Task) error
}

type Adapter struct {
	r repository
}

func New(r repository) *Adapter {
	return &Adapter{
		r: r,
	}
}

func (a *Adapter) Process(ctx context.Context, log models.Log) error {
	task := models.Task{
		UserID:          log.UserID,
		Action:          log.Action,
		Description:     log.Description,
		Timestamp:       log.Timestamp,
		Error:           log.Error,
		Attempts:        attemptsAmount,
		Status:          status.Created,
		Created_at:      time.Now(),
		Updated_at:      time.Now(),
		Complited_at:    time.Now(),
		Processing_from: time.Now(),
	}

	return a.r.AddTask(ctx, task)
}
