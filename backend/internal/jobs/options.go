package jobs

import (
	"time"

	"github.com/hibiken/asynq"
)

// Queue names used by the worker. Weights are configured in NewServer.
const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueBulk     = "bulk"
)

// DefaultOpts returns MaxRetry, Timeout and Queue for a job type. Callers can
// still append/override options when publishing.
func DefaultOpts(t JobType) []asynq.Option {
	switch t {
	case JobSendEmail, JobSendWhatsApp, JobSendNotification, JobProcessScheduledNotifications:
		return []asynq.Option{
			asynq.Queue(QueueCritical),
			asynq.MaxRetry(5),
			asynq.Timeout(30 * time.Second),
		}
	case JobImportProcess, JobExportProcess:
		return []asynq.Option{
			asynq.Queue(QueueBulk),
			asynq.MaxRetry(3),
			asynq.Timeout(10 * time.Minute),
		}
	default:
		return []asynq.Option{
			asynq.Queue(QueueDefault),
			asynq.MaxRetry(3),
			asynq.Timeout(60 * time.Second),
		}
	}
}
