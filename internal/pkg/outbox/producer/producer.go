package producer

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
	ttopic "github.com/ereminiu/pvz/internal/pkg/outbox/topic"
)

const (
	timeout      = 50 * time.Millisecond
	flushTimeout = 5000
)

type repository interface {
	GetTasks(ctx context.Context) ([]*models.Task, error)
}

type Producer struct {
	r        repository
	producer *kafka.Producer

	stop chan struct{}
}

func New(r repository, producer *kafka.Producer) *Producer {
	return &Producer{
		r:        r,
		producer: producer,
		stop:     make(chan struct{}),
	}
}

func (p *Producer) Produce(ctx context.Context) error {
	ticker := time.NewTicker(timeout)

	for {
		select {
		case <-ticker.C:
			if err := p.produce(ctx); err != nil {
				return err
			}

		case <-p.stop:
			return nil
		}
	}
}

func (p *Producer) produce(ctx context.Context) error {
	tasks, err := p.r.GetTasks(ctx)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		deliverCh := make(chan kafka.Event)
		topic := ttopic.Task
		message, err := sonic.Marshal(task)
		if err != nil {
			return err
		}

		if err = p.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: message,
			Key:   nil,
		}, deliverCh); err != nil {
			return err
		}

		e := <-deliverCh
		switch event := e.(type) {
		case *kafka.Message:
			fmt.Printf("write to kafka: %s\n", event.String())
			continue

		case *kafka.Error:
			return event

		default:
			return fmt.Errorf("unknown message type")
		}
	}

	return nil
}

func (p *Producer) Stop() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()

	p.stop <- struct{}{}
}
