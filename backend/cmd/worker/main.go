package main

import (
	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notify"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/logger"
)

// The worker is a separate process from the API so async processing scales
// independently. It consumes tasks enqueued by the API onto the asynq queue and
// routes notification work through the shared notify pipeline.
func main() {
	cfg := config.Load()

	log := logger.New()
	defer log.Sync()

	// Wire notification providers. The simulation provider is the default for
	// every channel; real vendors (Meta Cloud API, Twilio, Gupshup, 360Dialog,
	// SMTP) are registered here in later phases without touching job handlers.
	dispatcher := notify.NewDispatcher(log)
	dispatcher.Register(notify.NewSimulationProvider("simulation", notify.ChannelEmail, log))
	dispatcher.Register(notify.NewSimulationProvider("simulation", notify.ChannelWhatsApp, log))

	server := jobs.NewServer(
		jobs.RedisOpt(
			cfg.RedisHost,
			cfg.RedisPort,
			cfg.RedisPassword,
			cfg.RedisDB,
		),
		log,
		dispatcher,
	)

	if err := server.Run(); err != nil {
		log.Sugar().Fatalf("worker failed: %v", err)
	}
}
