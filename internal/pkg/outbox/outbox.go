package outbox

import (
	"context"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
	"github.com/ereminiu/pvz/internal/pkg/outbox/consumer"
	"github.com/ereminiu/pvz/internal/pkg/outbox/producer"
	"github.com/ereminiu/pvz/internal/pkg/outbox/topic"
)

type repository interface {
	AddTask(ctx context.Context, task models.Task) error
	GetTasks(ctx context.Context) ([]*models.Task, error)
}

type Outbox struct {
	p *producer.Producer
	c *consumer.Consumer
}

func New(r repository, adress string) (*Outbox, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": adress,
	})
	if err != nil {
		return nil, err
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        adress,
		"group.id":                 topic.TaskGroup,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.offset.reset":        "earliest",
	})
	if err != nil {
		return nil, err
	}

	if err := c.Subscribe("task", nil); err != nil {
		return nil, err
	}

	return &Outbox{
		p: producer.New(r, p),
		c: consumer.New(r, &consumer.StdinHandler{}, c),
	}, nil
}

func (outbox *Outbox) Run(ctx context.Context) {
	go func() {
		if err := outbox.produce(ctx); err != nil {
			slog.ErrorContext(ctx, "error during producing in outbox", slog.Any("err", err))
		}
	}()

	go func() {
		if err := outbox.consume(ctx); err != nil {
			slog.ErrorContext(ctx, "error during consuming in outbox", slog.Any("err", err))
		}
	}()
}

func (outbox *Outbox) produce(ctx context.Context) error {
	return outbox.p.Produce(ctx)
}

func (outbox *Outbox) consume(ctx context.Context) error {
	return outbox.c.Consume(ctx)
}

func (outbox *Outbox) Stop() {
	outbox.p.Stop()
	outbox.c.Stop()
}
