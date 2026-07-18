package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/activity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/app"
	"github.com/abhinavkumar03/crm-lite/backend/internal/attachment"
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth"
	"github.com/abhinavkumar03/crm-lite/backend/internal/calllog"
	"github.com/abhinavkumar03/crm-lite/backend/internal/contact"
	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard"
	"github.com/abhinavkumar03/crm-lite/backend/internal/field"
	"github.com/abhinavkumar03/crm-lite/backend/internal/health"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead"
	"github.com/abhinavkumar03/crm-lite/backend/internal/media"
	moduleengine "github.com/abhinavkumar03/crm-lite/backend/internal/module"
	"github.com/abhinavkumar03/crm-lite/backend/internal/note"
	"github.com/abhinavkumar03/crm-lite/backend/internal/search"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/database"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/logger"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/redis"
	"github.com/abhinavkumar03/crm-lite/backend/internal/task"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine"
	"github.com/abhinavkumar03/crm-lite/backend/internal/view"
)

func main() {

	cfg := config.Load()

	log := logger.New()
	defer log.Sync()

	dsn := database.BuildDSN(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	log.Info("Connecting to PostgreSQL...")

	db, err := database.New(dsn)
	if err != nil {
		log.Sugar().Fatalf("Postgres connection failed: %v", err)
	}
	defer db.Close()

	log.Info("PostgreSQL connected")

	log.Info("Connecting to Redis...")

	redisClient, err := redis.New(cfg)
	if err != nil {
		log.Sugar().Fatalf("Redis connection failed: %v", err)
	}

	log.Info("Redis connected")

	// The producer enqueues async work onto the asynq queue; the worker that
	// consumes it runs as a separate process (cmd/worker).
	producer := jobs.NewProducer(jobs.RedisOpt(
		cfg.RedisHost,
		cfg.RedisPort,
		cfg.RedisPassword,
		cfg.RedisDB,
	))
	defer producer.Close()

	// Resolves the authenticated user's organization; shared by all
	// organization-scoped (metadata-driven) modules.
	orgMiddleware := tenant.Middleware(tenant.NewResolver(db))

	healthModule := health.NewModule()
	authModule := auth.NewModule(db, cfg.JWTSecret, cfg.JWTExpiration)
	moduleEngine := moduleengine.NewModule(db, authModule.Middleware(), orgMiddleware)
	fieldEngine := field.NewModule(db, authModule.Middleware(), orgMiddleware)
	validationEngine := validationengine.NewModule(db, authModule.Middleware(), orgMiddleware)
	viewEngine := view.NewModule(db, authModule.Middleware(), orgMiddleware)
	recordEngine := record.NewModule(db, authModule.Middleware(), orgMiddleware)
	leadModule := lead.NewModule(db, authModule.Middleware(), producer)
	contactModule := contact.NewModule(db, authModule.Middleware())
	taskModule := task.NewModule(db, authModule.Middleware())
	dashboardModule := dashboard.NewModule(db, redisClient, authModule.Middleware())
	searchModule := search.NewModule(db, authModule.Middleware())
	noteModule := note.NewModule(db, authModule.Middleware())
	calllogModule := calllog.NewModule(db, authModule.Middleware())
	attachmentModule := attachment.NewModule(db, authModule.Middleware())
	activityModule := activity.NewModule(db, authModule.Middleware())
	mediaModule, err := media.NewModule(cfg, authModule.Middleware())
	if err != nil {
		log.Sugar().Fatalf("failed to initialize media module: %v", err)
	}

	router := app.NewRouter(
		log,
		cfg,
		healthModule,
		authModule,
		moduleEngine,
		fieldEngine,
		validationEngine,
		viewEngine,
		recordEngine,
		leadModule,
		contactModule,
		taskModule,
		dashboardModule,
		searchModule,
		noteModule,
		calllogModule,
		attachmentModule,
		activityModule,
		mediaModule,
	)

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		log.Info("server started on port " + cfg.AppPort)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err.Error())
	}

	log.Info("server stopped")
}
