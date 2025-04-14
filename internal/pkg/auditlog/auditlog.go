package auditlog

import (
	"context"

	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
	wp "github.com/ereminiu/pvz/internal/pkg/auditlog/worker_pool"
	ob "github.com/ereminiu/pvz/internal/pkg/outbox"
)

type Adapter interface {
	Process(ctx context.Context, event models.Log) error
}

type Filter interface {
	Check(event models.Log) bool
}

type AuditLog struct {
	workerPool *wp.WorkerPool
	outbox     *ob.Outbox
	filter     Filter
}

func New(adapter Adapter, filter Filter, outbox *ob.Outbox) *AuditLog {
	return &AuditLog{
		workerPool: wp.NewWorkerPool(
			2,
			adapter,
		),
		outbox: outbox,
		filter: filter,
	}
}

func (audit *AuditLog) Run(ctx context.Context) {
	audit.outbox.Run(ctx)

	audit.workerPool.Run(ctx)
}

func (audit *AuditLog) Send(log models.Log) {
	if audit.filter.Check(log) {
		audit.workerPool.Add(log)
	}
}

func (audit *AuditLog) Stop() {
	audit.workerPool.Stop()
	audit.outbox.Stop()
}
