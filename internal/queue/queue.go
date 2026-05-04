package queue

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
)

// Job represents a background job to be processed.
type Job struct {
	Type    string
	Payload interface{}
}

// Handler processes a specific job type.
type Handler func(ctx context.Context, payload interface{}) error

// Queue is a simple in-memory job queue.
// For production at scale, replace with Redis-based queue (e.g., asynq).
type Queue struct {
	jobs     chan Job
	handlers map[string]Handler
	mu       sync.RWMutex
	log      zerolog.Logger
}

// New creates a new job queue.
func New(bufferSize int, log zerolog.Logger) *Queue {
	return &Queue{
		jobs:     make(chan Job, bufferSize),
		handlers: make(map[string]Handler),
		log:      log.With().Str("component", "queue").Logger(),
	}
}

// Register adds a handler for a job type.
func (q *Queue) Register(jobType string, handler Handler) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.handlers[jobType] = handler
}

// Enqueue adds a job to the queue.
func (q *Queue) Enqueue(job Job) {
	select {
	case q.jobs <- job:
		q.log.Debug().Str("job_type", job.Type).Msg("job enqueued")
	default:
		q.log.Warn().Str("job_type", job.Type).Msg("job queue full, dropping job")
	}
}

// Start begins processing jobs with the given number of workers.
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

// Close shuts down the queue.
func (q *Queue) Close() {
	close(q.jobs)
}
