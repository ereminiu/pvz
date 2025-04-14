package consumer

import (
	"context"
	"log/slog"
	"time"

	"github.com/bytedance/sonic"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
	"github.com/ereminiu/pvz/internal/pkg/outbox/status"
)

type repository interface {
	AddTask(ctx context.Context, task models.Task) error
}

type MessageHandler interface {
	Process(ctx context.Context, task models.Task) error
}

type Consumer struct {
	r        repository
	handler  MessageHandler
	consumer *kafka.Consumer

	stop chan struct{}
}

func New(r repository, handler MessageHandler, consumer *kafka.Consumer) *Consumer {
	return &Consumer{
		r:        r,
		handler:  handler,
		consumer: consumer,
		stop:     make(chan struct{}),
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	for {
		select {
		case <-c.stop:
			return nil

		default:
			kafkaMsg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				slog.Error("error during reading message", slog.Any("error", err))
				continue
			}

			if kafkaMsg == nil {
				continue
			}

			var task models.Task
			if err = sonic.Unmarshal(kafkaMsg.Value, &task); err != nil {
				slog.Error("err", slog.Any("error", err))
				continue
			}

			if err = c.handler.Process(ctx, task); err != nil {
				slog.Error("error duing", slog.Any("error", err))

				now := time.Now()

				task.Updated_at = now
				task.Processing_from = now.Add(2 * time.Second)
				task.Attempts--
				if task.Attempts == 0 {
					task.Status = status.Failed
					task.Complited_at = now
				}

				if err = c.r.AddTask(ctx, task); err != nil {
					slog.Error("error duing sending retry", slog.Any("error", err))
				}
			}
		}
	}
}

func (c *Consumer) Stop() {
	c.stop <- struct{}{}

	_, err := c.consumer.Commit()
	if err != nil {
		slog.Error("error during commit", slog.Any("err", err))
	}

	if err := c.consumer.Close(); err != nil {
		slog.Error("error during closing consumer", slog.Any("err", err))
	}
}
