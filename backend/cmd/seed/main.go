// Command seed runs the database seeders. Seeders are ordered, history-tracked
// and idempotent (see internal/seed). Pass -fresh to re-run all seeders.
//
// Usage:
//
//	seed            apply pending seeders
//	seed -fresh     clear history and re-run every seeder
package main

import (
	"context"
	"flag"

	"github.com/abhinavkumar03/crm-lite/backend/internal/seed"
	"github.com/abhinavkumar03/crm-lite/backend/internal/seed/seeders"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/database"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/logger"
)

func main() {
	fresh := flag.Bool("fresh", false, "clear seed history and re-run all seeders")
	flag.Parse()

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

	db, err := database.New(dsn)
	if err != nil {
		log.Sugar().Fatalf("Postgres connection failed: %v", err)
	}
	defer db.Close()

	runner := seed.NewRunner(db, log)

	// Register seeders in dependency order: catalog -> tenant -> metadata ->
	// business data. Each later seeder looks up what it needs by natural key so
	// they remain decoupled and idempotent.
	runner.Register(
		seeders.PermissionsSeeder{},
		seeders.UsersSeeder{},
		seeders.OrganizationSeeder{},
		seeders.RolesSeeder{},
		seeders.MembershipsSeeder{},
		seeders.OrgStructureSeeder{},
		seeders.ModulesSeeder{},
		seeders.FieldsSeeder{},
		seeders.LayoutsSeeder{},
		seeders.TourStepsSeeder{},
		seeders.BusinessDataSeeder{},
		seeders.WorkspaceSideDataSeeder{},
		seeders.NotificationDemoSeeder{},
	)

	if err := runner.Run(context.Background(), *fresh); err != nil {
		log.Sugar().Fatalf("seeding failed: %v", err)
	}

	log.Info("seeding complete")
}
