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
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/database"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/logger"
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

	db, err := database.New(dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	authModule := auth.NewModule(
		db,
		cfg.JWTSecret,
	)
	router := app.NewRouter(
		log,
		authModule,
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
