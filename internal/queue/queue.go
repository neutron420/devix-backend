package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Job struct {
	Type    string
	Payload interface{}
}

type JobData struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}

type Handler func(ctx context.Context, payload interface{}) error

type Queue struct {
	redis      *redis.Client
	streamName string
	groupName  string
	handlers   map[string]Handler
	mu         sync.RWMutex
	log        zerolog.Logger
	stopChan   chan struct{}
}

func New(rdb *redis.Client, streamName string, log zerolog.Logger) *Queue {
	return &Queue{
		redis:      rdb,
		streamName: streamName,
		groupName:  "devix-workers",
		handlers:   make(map[string]Handler),
		log:        log.With().Str("component", "queue").Logger(),
		stopChan:   make(chan struct{}),
	}
}

func (q *Queue) Register(jobType string, handler Handler) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.handlers[jobType] = handler
}

func (q *Queue) Enqueue(job Job) {
	if q.redis == nil {
		q.log.Warn().Msg("redis is nil, cannot enqueue job (skipping)")
		return
	}

	payloadBytes, err := json.Marshal(job.Payload)
	if err != nil {
		q.log.Error().Err(err).Str("job_type", job.Type).Msg("failed to marshal job payload")
		return
	}

	jobData, err := json.Marshal(JobData{
		Type:    job.Type,
		Payload: payloadBytes,
	})
	if err != nil {
		q.log.Error().Err(err).Str("job_type", job.Type).Msg("failed to marshal job data")
		return
	}

	err = q.redis.XAdd(context.Background(), &redis.XAddArgs{
		Stream: q.streamName,
		Values: map[string]interface{}{
			"data": jobData,
		},
	}).Err()

	if err != nil {
		q.log.Error().Err(err).Str("job_type", job.Type).Msg("failed to enqueue job to redis")
	} else {
		q.log.Debug().Str("job_type", job.Type).Msg("job enqueued")
	}
}

func (q *Queue) Start(ctx context.Context, numWorkers int) {
	if q.redis == nil {
		q.log.Warn().Msg("redis is nil, background queue processing disabled")
		return
	}

	err := q.redis.XGroupCreateMkStream(ctx, q.streamName, q.groupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		q.log.Error().Err(err).Msg("failed to create redis consumer group")
		return
	}

	for i := 0; i < numWorkers; i++ {
		go q.worker(ctx, i)
	}
	q.log.Info().Int("workers", numWorkers).Msg("redis job queue started")
}

func (q *Queue) worker(ctx context.Context, id int) {
	consumerName := fmt.Sprintf("worker-%d", id)
	q.log.Debug().Int("worker_id", id).Msg("worker started")
	for {
		select {
		case <-ctx.Done():
			q.log.Debug().Int("worker_id", id).Msg("worker stopped (context done)")
			return
		case <-q.stopChan:
			q.log.Debug().Int("worker_id", id).Msg("worker stopped (queue closed)")
			return
		default:
			streams, err := q.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    q.groupName,
				Consumer: consumerName,
				Streams:  []string{q.streamName, ">"},
				Count:    1,
				Block:    2 * time.Second,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue
				}
				if err == context.Canceled {
					return
				}
				q.log.Error().Err(err).Msg("error reading from redis stream")
				time.Sleep(1 * time.Second)
				continue
			}

			for _, stream := range streams {
				for _, msg := range stream.Messages {
					dataStr, ok := msg.Values["data"].(string)
					if !ok {
						q.log.Error().Msg("invalid message data format in stream")
						q.redis.XAck(ctx, q.streamName, q.groupName, msg.ID)
						continue
					}

					var jobData JobData
					if err := json.Unmarshal([]byte(dataStr), &jobData); err != nil {
						q.log.Error().Err(err).Msg("failed to unmarshal job container")
						q.redis.XAck(ctx, q.streamName, q.groupName, msg.ID)
						continue
					}

					q.mu.RLock()
					handler, ok := q.handlers[jobData.Type]
					q.mu.RUnlock()

					if !ok {
						q.log.Warn().Str("job_type", jobData.Type).Msg("no handler registered")
						q.redis.XAck(ctx, q.streamName, q.groupName, msg.ID)
						continue
					}

					var genericPayload interface{}
					if err := json.Unmarshal(jobData.Payload, &genericPayload); err != nil {
						q.log.Error().Err(err).Str("job_type", jobData.Type).Msg("failed to unmarshal generic payload")
						q.redis.XAck(ctx, q.streamName, q.groupName, msg.ID)
						continue
					}

					if err := handler(ctx, genericPayload); err != nil {
						q.log.Error().Err(err).Str("job_type", jobData.Type).Msg("job failed")
						q.redis.XAck(ctx, q.streamName, q.groupName, msg.ID)
					} else {
						q.log.Debug().Str("job_type", jobData.Type).Msg("job completed")
						q.redis.XAck(ctx, q.streamName, q.groupName, msg.ID)
					}
				}
			}
		}
	}
}

func (q *Queue) Close() {
	close(q.stopChan)
}
