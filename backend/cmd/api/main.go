package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/app"
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth"
	"github.com/abhinavkumar03/crm-lite/backend/internal/contact"
	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard"
	"github.com/abhinavkumar03/crm-lite/backend/internal/health"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead"
	"github.com/abhinavkumar03/crm-lite/backend/internal/search"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/database"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/logger"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/redis"
	"github.com/abhinavkumar03/crm-lite/backend/internal/task"
)

func main() {

	cfg := config.Load()

	log := logger.New()

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

	log.Info("PostgreSQL connected")

	log.Info("Connecting to Redis...")

	redisClient, err := redis.New(cfg)
	if err != nil {
		log.Sugar().Fatalf("Redis connection failed: %v", err)
	}

	log.Info("Redis connected")

	producer := jobs.NewProducer(
		redisClient,
	)

	worker := jobs.NewWorker(
		redisClient,
	)

	go worker.Start(
		context.Background(),
	)

	healthModule := health.NewModule()
	authModule := auth.NewModule(db, cfg.JWTSecret)
	leadModule := lead.NewModule(db, authModule.Middleware(), producer)
	contactModule := contact.NewModule(db, authModule.Middleware())
	taskModule := task.NewModule(db, authModule.Middleware())
	dashboardModule := dashboard.NewModule(db, redisClient, authModule.Middleware())
	searchModule := search.NewModule(db, authModule.Middleware())
	router := app.NewRouter(
		log,
		healthModule,
		authModule,
		leadModule,
		contactModule,
		taskModule,
		dashboardModule,
		searchModule,
	)

	application := &app.Application{
		Config: cfg,
		Logger: log,
		Router: router,
	}

	defer log.Sync()

	server := &http.Server{
		Addr:    ":" + application.Config.AppPort,
		Handler: application.Router,
	}

	go func() {

		log.Info("server started")

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
