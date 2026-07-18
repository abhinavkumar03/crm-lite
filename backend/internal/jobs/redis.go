package jobs

import (
	"fmt"

	"github.com/hibiken/asynq"
)

// RedisOpt builds the asynq Redis connection options from discrete config
// values. Both the producer (cmd/api) and the worker (cmd/worker) use this so
// they always agree on the queue backend.
func RedisOpt(host, port, password string, db int) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	}
}
