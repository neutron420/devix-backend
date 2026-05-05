package queue

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
)

type Job struct {
	Type    string
	Payload interface{}
}

type Handler func(ctx context.Context, payload interface{}) error

type Queue struct {
	jobs     chan Job
	handlers map[string]Handler
	mu       sync.RWMutex
	log      zerolog.Logger
}

func New(bufferSize int, log zerolog.Logger) *Queue {
	return &Queue{
		jobs:     make(chan Job, bufferSize),
		handlers: make(map[string]Handler),
		log:      log.With().Str("component", "queue").Logger(),
	}
}

func (q *Queue) Register(jobType string, handler Handler) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.handlers[jobType] = handler
}

func (q *Queue) Enqueue(job Job) {
	select {
	case q.jobs <- job:
		q.log.Debug().Str("job_type", job.Type).Msg("job enqueued")
	default:
		q.log.Warn().Str("job_type", job.Type).Msg("job queue full, dropping job")
	}
}

func (q *Queue) Start(ctx context.Context, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go q.worker(ctx, i)
	}
	q.log.Info().Int("workers", numWorkers).Msg("job queue started")
}

func (q *Queue) worker(ctx context.Context, id int) {
	q.log.Debug().Int("worker_id", id).Msg("worker started")
	for {
		select {
		case <-ctx.Done():
			q.log.Debug().Int("worker_id", id).Msg("worker stopped")
			return
		case job := <-q.jobs:
			q.mu.RLock()
			handler, ok := q.handlers[job.Type]
			q.mu.RUnlock()

			if !ok {
				q.log.Warn().Str("job_type", job.Type).Msg("no handler registered")
				continue
			}

			if err := handler(ctx, job.Payload); err != nil {
				q.log.Error().Err(err).Str("job_type", job.Type).Msg("job failed")
			} else {
				q.log.Debug().Str("job_type", job.Type).Msg("job completed")
			}
		}
	}
}

func (q *Queue) Close() {
	close(q.jobs)
}
