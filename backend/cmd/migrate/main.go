// Command migrate is the database migration runner. It wraps golang-migrate
// with the embedded migration files (see internal package "migrations") and the
// pgx v5 driver, providing an ordered, versioned, idempotent, rollback-capable
// workflow.
//
// Usage:
//
//	migrate up                 apply all pending migrations
//	migrate down [n]           roll back the last n migrations (default 1)
//	migrate version            print current version and dirty state
//	migrate force <version>    force set version (recover from a dirty state)
//	migrate drop               DROP everything (development only)
//	migrate create <name>      scaffold the next migration file pair (dev tool)
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/config"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/database"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/logger"
	"github.com/abhinavkumar03/crm-lite/backend/migrations"
)

// migrationsDir is the on-disk location (relative to the backend module root)
// used only by the "create" scaffolding command. Applying migrations uses the
// embedded filesystem instead.
const migrationsDir = "migrations"

func main() {
	log := logger.New()
	defer log.Sync()

	args := os.Args[1:]
	if len(args) == 0 {
		log.Sugar().Fatalf("missing command; expected one of: up, down, version, force, drop, create")
	}

	command := args[0]

	// "create" is a local scaffolding tool and needs no database connection.
	if command == "create" {
		if len(args) < 2 {
			log.Sugar().Fatalf("usage: migrate create <name>")
		}
		if err := createMigration(strings.Join(args[1:], "_")); err != nil {
			log.Sugar().Fatalf("create failed: %v", err)
		}
		return
	}

	cfg := config.Load()

	m, err := newMigrator(cfg)
	if err != nil {
		log.Sugar().Fatalf("failed to initialize migrator: %v", err)
	}
	defer func() {
		if srcErr, dbErr := m.Close(); srcErr != nil || dbErr != nil {
			log.Sugar().Warnf("migrator close: source=%v db=%v", srcErr, dbErr)
		}
	}()

	switch command {
	case "up":
		run(log, "migrations applied", m.Up())

	case "down":
		steps := 1
		if len(args) > 1 {
			steps, err = strconv.Atoi(args[1])
			if err != nil || steps < 1 {
				log.Sugar().Fatalf("invalid steps: %q", args[1])
			}
		}
		run(log, "migrations rolled back", m.Steps(-steps))

	case "version":
		version, dirty, verErr := m.Version()
		if errors.Is(verErr, migrate.ErrNilVersion) {
			log.Info("no migrations applied yet")
			return
		}
		if verErr != nil {
			log.Sugar().Fatalf("version failed: %v", verErr)
		}
		log.Sugar().Infof("version=%d dirty=%t", version, dirty)

	case "force":
		if len(args) < 2 {
			log.Sugar().Fatalf("usage: migrate force <version>")
		}
		version, convErr := strconv.Atoi(args[1])
		if convErr != nil {
			log.Sugar().Fatalf("invalid version: %q", args[1])
		}
		run(log, "version forced", m.Force(version))

	case "drop":
		if cfg.IsProduction() && os.Getenv("MIGRATE_ALLOW_DROP") != "true" {
			log.Sugar().Fatalf("refusing to drop in production; set MIGRATE_ALLOW_DROP=true to override")
		}
		run(log, "database dropped", m.Drop())

	default:
		log.Sugar().Fatalf("unknown command: %q", command)
	}
}

func newMigrator(cfg *config.Config) (*migrate.Migrate, error) {
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("load embedded migrations: %w", err)
	}

	dbURL := database.MigrationURL(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	m, err := migrate.NewWithSourceInstance("iofs", src, dbURL)
	if err != nil {
		return nil, fmt.Errorf("init migrate: %w", err)
	}

	return m, nil
}

// run logs a friendly message, treating ErrNoChange as success.
func run(log *zap.Logger, successMsg string, err error) {
	if errors.Is(err, migrate.ErrNoChange) {
		log.Info("no change: schema already up to date")
		return
	}
	if err != nil {
		log.Sugar().Fatalf("migration failed: %v", err)
	}
	log.Info(successMsg)
}

// createMigration scaffolds the next sequential migration file pair using the
// existing 6-digit zero-padded naming convention.
func createMigration(name string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read %s: %w", migrationsDir, err)
	}

	next := nextVersion(entries)
	base := fmt.Sprintf("%06d_%s", next, sanitize(name))

	files := []string{
		filepath.Join(migrationsDir, base+".up.sql"),
		filepath.Join(migrationsDir, base+".down.sql"),
	}

	for _, f := range files {
		if err := os.WriteFile(f, []byte("-- write your migration here\n"), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", f, err)
		}
		fmt.Println("created", f)
	}

	return nil
}

func nextVersion(entries []os.DirEntry) int {
	versions := make([]int, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		nameParts := strings.SplitN(e.Name(), "_", 2)
		if len(nameParts) < 2 {
			continue
		}
		if v, err := strconv.Atoi(nameParts[0]); err == nil {
			versions = append(versions, v)
		}
	}
	sort.Ints(versions)
	if len(versions) == 0 {
		return 1
	}
	return versions[len(versions)-1] + 1
}

func sanitize(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	return name
}
