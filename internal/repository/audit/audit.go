package audit

import (
	"context"
	"time"

	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
	txmanager "github.com/ereminiu/pvz/internal/pkg/tx_manager"
	"github.com/georgysavva/scany/pgxscan"
)

const (
	statusCreated    = "CREATED"
	statusProcessing = "PROCESSING"
	statusFailed     = "FAILED"
	statusNoAttempts = "NO_ATTEMPTS_LEFT"
)

type Repository struct {
	manager *txmanager.TxManager
}

func New(manager *txmanager.TxManager) *Repository {
	return &Repository{
		manager: manager,
	}
}

func (r *Repository) AddTask(ctx context.Context, task models.Task) error {
	query := `INSERT INTO tasks (user_id, action, description, timestamp, error, attempts, status, processing_from)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.manager.GetQueryEngine(ctx).Exec(
		ctx,
		query,
		task.UserID,
		task.Action,
		task.Description,
		task.Timestamp,
		task.Error,
		task.Attempts,
		task.Status,
		time.Now(),
	)

	return err
}

func (r *Repository) GetTasks(ctx context.Context) ([]*models.Task, error) {
	query := `DELETE FROM tasks 
			WHERE status=$1 AND attempts>$2
			RETURNING 
			user_id, action, description, timestamp, error, attempts, status, created_at, updated_at, processing_from`

	tasks := make([]*models.Task, 0)
	if err := pgxscan.Select(ctx, r.manager.GetQueryEngine(ctx), &tasks, query, statusCreated, 0); err != nil {
		return nil, err
	}

	return tasks, nil
}
