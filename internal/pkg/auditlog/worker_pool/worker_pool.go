package workerpool

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
)

const (
	timeout  = 500 * time.Millisecond
	maxBatch = 5
)

type Adapter interface {
	Process(ctx context.Context, event models.Log) error
}

type worker struct {
	ID    int
	batch chan models.Log
	stop  chan bool
	ready chan bool

	adapter Adapter

	batchSize atomic.Int32
}

func newWorker(ID int, adapter Adapter) *worker {
	return &worker{
		ID:      ID,
		batch:   make(chan models.Log, maxBatch),
		stop:    make(chan bool),
		ready:   make(chan bool),
		adapter: adapter,
	}
}

func (w *worker) send(log models.Log) {
	w.batch <- log

	if w.batchSize.Add(1) == maxBatch {
		w.ready <- true
	}
}

func (w *worker) process(ctx context.Context, amount int) {
	for i := 0; i < amount; i++ {
		job := <-w.batch

		if err := w.adapter.Process(ctx, job); err != nil {
			slog.Error("error during processing log", slog.Any("err", err))
		}
	}
}

func (w *worker) Run(ctx context.Context) {
	timer := time.NewTimer(timeout)

	for {
		select {
		case <-w.stop:
			w.process(ctx, len(w.batch))
			return

		case <-w.ready:
			timer.Reset(timeout)
			w.process(ctx, maxBatch)

		case <-timer.C:
			timer.Reset(timeout)
			w.process(ctx, len(w.batch))
		}
	}
}

type WorkerPool struct {
	jobs chan models.Log
	stop chan bool

	N       int
	workers []*worker
}

func NewWorkerPool(N int, adapter Adapter) *WorkerPool {
	workers := make([]*worker, N)
	for i := 0; i < N; i++ {
		workers[i] = newWorker(i+1, adapter)
	}

	return &WorkerPool{
		N:       N,
		jobs:    make(chan models.Log, 2*maxBatch),
		stop:    make(chan bool),
		workers: workers,
	}
}

func (wp *WorkerPool) Add(log models.Log) {
	wp.jobs <- log
}

func (wp *WorkerPool) Stop() {
	wp.stop <- true
}

func (wp *WorkerPool) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	for i := 0; i < wp.N; i++ {
		go wp.workers[i].Run(ctx)
	}

	for i := 0; i < wp.N; i++ {
		go func(i int) {
			wg.Done()
			for j := range wp.jobs {
				wp.workers[i].send(j)
			}
			wp.workers[i].stop <- true
		}(i)
	}

	<-wp.stop
	close(wp.jobs)
	wg.Wait()
}
