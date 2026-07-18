package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	"github.com/abhinavkumar03/crm-lite/backend/internal/notify"
)

// Server consumes and processes Jobs from the asynq queue. It runs in the
// dedicated worker process so async work scales independently of the API.
type Server struct {
	srv    *asynq.Server
	mux    *asynq.ServeMux
	logger *zap.Logger
}

// NewServer wires the asynq server and routes each JobType to a handler. The
// notification-oriented jobs are delegated to the shared notify.Dispatcher so
// email and WhatsApp travel the same pipeline.
func NewServer(opt asynq.RedisClientOpt, logger *zap.Logger, dispatcher *notify.Dispatcher) *Server {
	srv := asynq.NewServer(opt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"default": 1,
		},
		Logger: newZapLogger(logger),
	})

	h := &handlers{logger: logger, dispatcher: dispatcher}

	mux := asynq.NewServeMux()
	mux.HandleFunc(string(JobLeadCreated), h.handleLeadEvent)
	mux.HandleFunc(string(JobLeadStatusChanged), h.handleLeadEvent)
	mux.HandleFunc(string(JobSendEmail), h.handleSendEmail)
	mux.HandleFunc(string(JobSendWhatsApp), h.handleSendWhatsApp)

	return &Server{srv: srv, mux: mux, logger: logger}
}

// Run starts processing and blocks until the process receives a termination
// signal (asynq handles graceful shutdown internally).
func (s *Server) Run() error {
	s.logger.Info("jobs: worker started")
	return s.srv.Run(s.mux)
}

// handlers holds dependencies shared by all job handlers.
type handlers struct {
	logger     *zap.Logger
	dispatcher *notify.Dispatcher
}

func decode(t *asynq.Task) (Job, error) {
	var job Job
	if err := json.Unmarshal(t.Payload(), &job); err != nil {
		// Returning asynq.SkipRetry avoids retrying tasks that can never
		// succeed because their payload is malformed.
		return Job{}, fmt.Errorf("jobs: decode payload: %v: %w", err, asynq.SkipRetry)
	}
	return job, nil
}

func (h *handlers) handleLeadEvent(_ context.Context, t *asynq.Task) error {
	job, err := decode(t)
	if err != nil {
		return err
	}

	h.logger.Info("jobs: lead event processed",
		zap.String("type", string(job.Type)),
		zap.String("user_id", job.UserID),
		zap.Any("payload", job.Payload),
	)
	return nil
}

func (h *handlers) handleSendEmail(ctx context.Context, t *asynq.Task) error {
	job, err := decode(t)
	if err != nil {
		return err
	}

	return h.dispatcher.Dispatch(ctx, notify.Message{
		Channel:  notify.ChannelEmail,
		To:       stringField(job.Payload, "email"),
		Subject:  "Welcome to CRM Lite",
		Template: "lead_welcome",
		Data:     job.Payload,
	})
}

func (h *handlers) handleSendWhatsApp(ctx context.Context, t *asynq.Task) error {
	job, err := decode(t)
	if err != nil {
		return err
	}

	return h.dispatcher.Dispatch(ctx, notify.Message{
		Channel:  notify.ChannelWhatsApp,
		To:       stringField(job.Payload, "phone"),
		Template: stringField(job.Payload, "template"),
		Data:     job.Payload,
	})
}

func stringField(payload map[string]interface{}, key string) string {
	if payload == nil {
		return ""
	}
	if v, ok := payload[key].(string); ok {
		return v
	}
	return ""
}
