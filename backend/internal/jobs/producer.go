package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

// Producer enqueues Jobs onto the asynq queue. It is safe for concurrent use
// and should be shared for the lifetime of the process, then Close()d on
// shutdown.
type Producer struct {
	client *asynq.Client
}

func NewProducer(opt asynq.RedisClientOpt) *Producer {
	return &Producer{
		client: asynq.NewClient(opt),
	}
}

// Publish serializes the job and enqueues it as an asynq task. DefaultOpts
// (queue, MaxRetry, Timeout) are applied first; callers may pass extra asynq
// options to override or extend them (e.g. ProcessIn).
func (p *Producer) Publish(ctx context.Context, job Job, opts ...asynq.Option) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("jobs: marshal job: %w", err)
	}

	task := asynq.NewTask(string(job.Type), payload)
	merged := append(DefaultOpts(job.Type), opts...)

	if _, err := p.client.EnqueueContext(ctx, task, merged...); err != nil {
		return fmt.Errorf("jobs: enqueue %q: %w", job.Type, err)
	}

	return nil
}

// Close releases the underlying Redis connection.
func (p *Producer) Close() error {
	return p.client.Close()
}
