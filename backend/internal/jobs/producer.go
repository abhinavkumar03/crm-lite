package jobs

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Producer struct {
	redis *redis.Client
}

func NewProducer(
	client *redis.Client,
) *Producer {

	return &Producer{
		redis: client,
	}
}

func (p *Producer) Publish(
	ctx context.Context,
	job Job,
) error {

	bytes, err := json.Marshal(job)

	if err != nil {
		return err
	}

	return p.redis.RPush(
		ctx,
		"crm:jobs",
		bytes,
	).Err()
}
