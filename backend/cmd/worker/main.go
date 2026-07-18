package main

import (
	activityrepo "github.com/abhinavkumar03/crm-lite/backend/internal/activity/repository"
	fieldrepo "github.com/abhinavkumar03/crm-lite/backend/internal/field/repository"
	importprocessor "github.com/abhinavkumar03/crm-lite/backend/internal/importer/processor"
	importrepo "github.com/abhinavkumar03/crm-lite/backend/internal/importer/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	notificationprocessor "github.com/abhinavkumar03/crm-lite/backend/internal/notification/processor"
	notificationrepo "github.com/abhinavkumar03/crm-lite/backend/internal/notification/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notify"
	recordrepo "github.com/abhinavkumar03/crm-lite/backend/internal/record/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/database"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/logger"
	vrepo "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/repository"
	vservice "github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/service"
)

// The worker is a separate process from the API so async processing scales
// independently. It consumes tasks enqueued by the API onto the asynq queue and
// routes notification work through the shared notify pipeline, persisting
// delivery outcomes to the database.
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

	// Wire notification providers. Email always uses the simulation provider for
	// now; WhatsApp is config-driven (simulation by default, Meta Cloud API when
	// credentials are supplied) — a one-line vendor swap with no handler changes.
	dispatcher := notify.NewDispatcher(log)
	dispatcher.Register(notify.NewSimulationProvider("simulation", notify.ChannelEmail, log))
	dispatcher.Register(notify.BuildWhatsAppProvider(notify.WhatsAppConfig{
		Provider: cfg.WhatsAppProvider,
		APIURL:   cfg.WhatsAppAPIURL,
		Token:    cfg.WhatsAppToken,
		PhoneID:  cfg.WhatsAppPhoneID,
	}, log))

	// The processor delivers persisted notifications and logs the outcome.
	processor := notificationprocessor.New(
		notificationrepo.New(db),
		dispatcher,
		activityrepo.New(db),
		log,
	)

	// The import processor maps, validates (Phase 7 engine) and inserts (Phase 10
	// record repository) each staged row — an import obeys the same rules as an
	// API-created record.
	fieldRepo := fieldrepo.New(db)
	importProcessor := importprocessor.New(
		importrepo.New(db),
		recordrepo.New(db),
		fieldRepo,
		vservice.New(vrepo.New(db), fieldRepo),
		log,
	)

	server := jobs.NewServer(
		jobs.RedisOpt(
			cfg.RedisHost,
			cfg.RedisPort,
			cfg.RedisPassword,
			cfg.RedisDB,
		),
		log,
		dispatcher,
		processor,
		importProcessor,
	)

	if err := server.Run(); err != nil {
		log.Sugar().Fatalf("worker failed: %v", err)
	}
}
