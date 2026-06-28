package jobs

import (
	"context"
	"encoding/json"
	"log"

	"github.com/redis/go-redis/v9"
)

type Worker struct {
	redis *redis.Client
}

func NewWorker(
	client *redis.Client,
) *Worker {

	return &Worker{
		redis: client,
	}
}

func (w *Worker) Start(
	ctx context.Context,
) {

	log.Println("Worker started...")

	for {

		result, err := w.redis.BLPop(
			ctx,
			0,
			"crm:jobs",
		).Result()

		if err != nil {
			continue
		}

		if len(result) < 2 {
			continue
		}

		var job Job

		if err := json.Unmarshal(
			[]byte(result[1]),
			&job,
		); err != nil {
			continue
		}

		w.handle(job)
	}
}

func (w *Worker) handle(job Job) {

	switch job.Type {

	case JobLeadCreated:

		log.Println(
			"[Activity]",
			job.Payload,
		)

	case JobLeadStatusChanged:

		log.Println(
			"[Status Changed]",
			job.Payload,
		)

	case JobSendEmail:

		log.Println(
			"[Email]",
			job.Payload,
		)

	default:

		log.Println(
			"Unknown Job",
			job.Type,
		)
	}
}
