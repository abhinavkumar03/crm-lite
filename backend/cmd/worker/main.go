package main

import (
	"context"
	"time"

	"github.com/hibiken/asynq"

	activityrepo "github.com/abhinavkumar03/crm-lite/backend/internal/activity/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter"
	exportprocessor "github.com/abhinavkumar03/crm-lite/backend/internal/exporter/processor"
	fieldrepo "github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
	importprocessor "github.com/abhinavkumar03/crm-lite/backend/internal/importer/processor"
	importrepo "github.com/abhinavkumar03/crm-lite/backend/internal/importer/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	notificationprocessor "github.com/abhinavkumar03/crm-lite/backend/internal/notification/processor"
	notificationrepo "github.com/abhinavkumar03/crm-lite/backend/internal/notification/repository"
	notificationservice "github.com/abhinavkumar03/crm-lite/backend/internal/notification/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notify"
	recordrepo "github.com/abhinavkumar03/crm-lite/backend/internal/record/repository"
	recordservice "github.com/abhinavkumar03/crm-lite/backend/internal/record/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/database"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/logger"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/secrets"
	vrepo "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/repository"
	vservice "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow"
	workspacerepo "github.com/abhinavkumar03/crm-lite/backend/internal/workspace/repository"
	workspaceservice "github.com/abhinavkumar03/crm-lite/backend/internal/workspace/service"
)

func main() {
	cfg := config.Load()

	log := logger.New()
	defer log.Sync()

	dsn := database.BuildDSN(
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)
	db, err := database.New(dsn)
	if err != nil {
		log.Sugar().Fatalf("worker: postgres connection failed: %v", err)
	}
	defer db.Close()

	box, err := secrets.NewBox(cfg.CommunicationSecretsKey)
	if err != nil {
		log.Sugar().Warnf("worker: secrets box unavailable: %v", err)
	}

	dispatcher := notify.NewDispatcher(log)
	dispatcher.Register(notify.BuildEmailProvider(notify.EmailConfig{
		Provider:     cfg.EmailProvider,
		SMTPHost:     cfg.SMTPHost,
		SMTPPort:     cfg.SMTPPort,
		SMTPUsername: cfg.SMTPUsername,
		SMTPPassword: cfg.SMTPPassword,
		SMTPFrom:     firstNonEmpty(cfg.SMTPFrom, cfg.EmailFrom),
		Encryption:   cfg.SMTPEncryption,
		APIKey:       cfg.ResendAPIKey,
		From:         firstNonEmpty(cfg.EmailFrom, cfg.SMTPFrom),
		ReplyTo:      cfg.EmailReplyTo,
	}, log))
	dispatcher.Register(notify.BuildWhatsAppProvider(notify.WhatsAppConfig{
		Provider:   cfg.WhatsAppProvider,
		APIURL:     cfg.WhatsAppAPIURL,
		Token:      cfg.WhatsAppToken,
		PhoneID:    cfg.WhatsAppPhoneID,
		AccountSID: cfg.TwilioAccountSID,
		AuthToken:  cfg.TwilioAuthToken,
		FromNumber: cfg.TwilioFromNumber,
	}, log))

	if box != nil {
		providerSvc := notificationservice.NewProviderService(
			notificationrepo.NewProviderRepository(db), box, log,
		)
		dispatcher.SetResolver(providerSvc)
	}

	processor := notificationprocessor.New(
		notificationrepo.New(db),
		dispatcher,
		activityrepo.New(db),
		log,
	)
	processor.SetPublicBaseURL(cfg.PublicBaseURL)

	fieldRepo := fieldrepo.New(db)
	validator := vservice.New(vrepo.New(db), fieldRepo)
	recordSvc := recordservice.New(recordrepo.New(db), fieldRepo, validator, nil, nil)
	workspaceSvc := workspaceservice.New(workspacerepo.New(db))
	recordSvc.SetActivityLogger(workspaceSvc)
	recordSvc.SetListLayoutReader(workspaceSvc)

	redisOpt := jobs.RedisOpt(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB)
	producer := jobs.NewProducer(redisOpt)
	defer producer.Close()

	notifySvc := notificationservice.New(notificationrepo.New(db), producer)

	workflowModule := workflow.NewModule(workflow.ModuleDeps{
		DB: db, Producer: producer, Enabled: cfg.Features.Automation, Logger: log,
		Records: recordSvc, Notify: notifySvc, Notes: workspaceSvc, Activities: workspaceSvc,
	})
	recordSvc.SetMutationHook(workflowModule.Publisher)

	importProcessor := importprocessor.New(
		importrepo.New(db),
		recordrepo.New(db),
		fieldRepo,
		validator,
		log,
	)
	exportProcessor := exportprocessor.New(exporter.NewService(db, nil), log)

	// Periodic sweep for due scheduled + retrying notifications.
	scheduler := asynq.NewScheduler(redisOpt, nil)
	sweepTask := asynq.NewTask(string(jobs.JobProcessScheduledNotifications), []byte(`{"type":"notification.process_scheduled","user_id":"","payload":{}}`))
	_, err = scheduler.Register("@every 1m", sweepTask, jobs.DefaultOpts(jobs.JobProcessScheduledNotifications)...)
	if err != nil {
		log.Sugar().Warnf("worker: schedule sweep register failed: %v", err)
	}

	wfSweep := asynq.NewTask(string(jobs.JobWorkflowScheduledSweep), []byte(`{"type":"workflow.scheduled_sweep","user_id":"","payload":{}}`))
	_, err = scheduler.Register("@every 1m", wfSweep, jobs.DefaultOpts(jobs.JobWorkflowScheduledSweep)...)
	if err != nil {
		log.Sugar().Warnf("worker: workflow sweep register failed: %v", err)
	} else {
		go func() {
			if err := scheduler.Run(); err != nil {
				log.Sugar().Errorf("worker: scheduler stopped: %v", err)
			}
		}()
		log.Info("jobs: scheduled sweeps registered (@every 1m)")
	}

	server := jobs.NewServer(redisOpt, log, dispatcher, processor, importProcessor, exportProcessor, workflowModule.Processor)

	go func() {
		time.Sleep(2 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = processor.ProcessDueScheduled(ctx)
	}()

	if err := server.Run(); err != nil {
		log.Sugar().Fatalf("worker failed: %v", err)
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
