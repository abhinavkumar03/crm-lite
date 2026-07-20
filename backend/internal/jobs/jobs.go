// Package jobs is the asynchronous work boundary for the application. It is
// backed by asynq (a Redis-backed task queue) which provides retries, a
// dead-letter mechanism (archived tasks), scheduling and visibility timeouts
// out of the box. Producers enqueue Jobs from request handlers; a separate
// worker process (cmd/worker) consumes and executes them.
package jobs

// JobType is the logical name of an asynchronous task. It doubles as the asynq
// task type string.
type JobType string

const (
	JobLeadCreated       JobType = "lead.created"
	JobLeadStatusChanged JobType = "lead.status_changed"
	JobSendEmail         JobType = "email.send"
	JobSendWhatsApp      JobType = "whatsapp.send"
	// JobSendNotification delivers a persisted notification (looked up by id) via
	// the notification pipeline: dispatch -> update status -> log activity.
	JobSendNotification JobType = "notification.send"
	// JobProcessScheduledNotifications promotes due scheduled notifications and
	// delivers them through the same pipeline.
	JobProcessScheduledNotifications JobType = "notification.process_scheduled"
	// JobImportProcess processes a staged import job (looked up by id): map,
	// validate and insert each row, then record progress and the error report.
	JobImportProcess JobType = "import.process"
	// JobExportProcess builds a persisted export (looked up by id): query the
	// records, serialize the file and store it for download.
	JobExportProcess JobType = "export.process"
)

// Job is the transport-agnostic payload enqueued for asynchronous processing.
// It is JSON-encoded into the asynq task payload.
type Job struct {
	Type    JobType                `json:"type"`
	UserID  string                 `json:"user_id"`
	Payload map[string]interface{} `json:"payload"`
}
