package jobs

type JobType string

const (
	JobLeadCreated JobType = "lead.created"

	JobLeadStatusChanged JobType = "lead.status_changed"

	JobSendEmail JobType = "email.send"
)

type Job struct {
	Type JobType `json:"type"`

	UserID string `json:"user_id"`

	Payload map[string]interface{} `json:"payload"`
}
