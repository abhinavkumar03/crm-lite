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
	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard"
	"github.com/abhinavkumar03/crm-lite/backend/internal/demo"
	"github.com/abhinavkumar03/crm-lite/backend/internal/docs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter"
	"github.com/abhinavkumar03/crm-lite/backend/internal/field"
	"github.com/abhinavkumar03/crm-lite/backend/internal/health"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer"
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/media"
	moduleengine "github.com/abhinavkumar03/crm-lite/backend/internal/module"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification"
	"github.com/abhinavkumar03/crm-lite/backend/internal/organization"
	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/record"
	"github.com/abhinavkumar03/crm-lite/backend/internal/roles"
	"github.com/abhinavkumar03/crm-lite/backend/internal/search"
	"github.com/abhinavkumar03/crm-lite/backend/internal/settings"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/cache"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/database"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/logger"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/redis"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tour"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine"
	"github.com/abhinavkumar03/crm-lite/backend/internal/view"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workspace"
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

	appCache := cache.New(redisClient)

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
	// organization-scoped (metadata-driven) modules. rbac.Load attaches the
	// role's permission keys so Require()/RequireModule() can enforce them.
	// Membership + RBAC grants are cached briefly in Redis (Phase 17).
	tenantResolver := tenant.NewResolver(db, appCache)
	orgMiddleware := tenant.Middleware(tenantResolver)
	guard := rbac.New(db, appCache)
	rbacLoad := guard.Load()

	healthModule := health.NewModule()
	docsModule := docs.NewModule()
	authModule := auth.NewModule(db, cfg.JWTSecret, cfg.JWTExpiration)
	authMW := authModule.Middleware()

	organizationModule := organization.NewModule(db, tenantResolver, authMW, orgMiddleware, rbacLoad, guard)
	moduleEngine := moduleengine.NewModule(db, authMW, orgMiddleware, rbacLoad, guard)
	fieldEngine := field.NewModule(db, authMW, orgMiddleware, rbacLoad, guard)
	validationEngine := validationengine.NewModule(db, authMW, orgMiddleware, rbacLoad, guard)
	viewEngine := view.NewModule(db, authMW, orgMiddleware, rbacLoad, guard)
	recordEngine := record.NewModule(db, authMW, orgMiddleware, rbacLoad, guard, appCache) // cache: invalidate org dashboard on CUD
	workspaceModule := workspace.NewModule(db, authMW, orgMiddleware, rbacLoad, guard)
	recordEngine.Service.SetActivityLogger(workspaceModule.Service)
	importEngine := importer.NewModule(db, authMW, orgMiddleware, rbacLoad, guard, producer)
	exportEngine := exporter.NewModule(db, authMW, orgMiddleware, rbacLoad, guard, producer)
	notificationModule := notification.NewModule(db, authMW, orgMiddleware, rbacLoad, guard, producer)
	tourModule := tour.NewModule(db, authMW, orgMiddleware)
	demoModule := demo.NewModule(db, tenantResolver, authMW)
	settingsModule := settings.NewModule(db, authMW, orgMiddleware, rbacLoad, guard)
	rolesModule := roles.NewModule(db, appCache, authMW, orgMiddleware, rbacLoad, guard)
	dashboardModule := dashboard.NewModule(db, appCache, authMW, orgMiddleware)
	searchModule := search.NewModule(db, authMW, orgMiddleware)
	// Legacy native CRM packages (lead/contact/task/note/calllog) remain on disk
	// but are intentionally not registered — product surface is dynamic-only.
	attachmentModule := attachment.NewModule(db, authMW)
	activityModule := activity.NewModule(db, authMW)
	mediaModule, err := media.NewModule(cfg, authMW)
	if err != nil {
		log.Sugar().Fatalf("failed to initialize media module: %v", err)
	}

	router := app.NewRouter(
		log,
		cfg,
		healthModule,
		docsModule,
		authModule,
		organizationModule,
		moduleEngine,
		fieldEngine,
		validationEngine,
		viewEngine,
		recordEngine,
		workspaceModule,
		importEngine,
		exportEngine,
		notificationModule,
		tourModule,
		demoModule,
		settingsModule,
		rolesModule,
		dashboardModule,
		searchModule,
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
