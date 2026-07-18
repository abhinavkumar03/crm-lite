# CRM Lite backend — developer workflow.
# Configuration is read from the environment / backend/.env by each command.

MIGRATE := go run ./cmd/migrate
SEED    := go run ./cmd/seed

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

## ---- Build & quality ----

.PHONY: build
build: ## Build all binaries
	go build ./...

.PHONY: tidy
tidy: ## Sync go.mod/go.sum
	go mod tidy

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: run-api
run-api: ## Run the API server
	go run ./cmd/api

.PHONY: run-worker
run-worker: ## Run the async worker
	go run ./cmd/worker

## ---- Migrations ----

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## Roll back N migrations (default 1): make migrate-down n=2
	$(MIGRATE) down $(n)

.PHONY: migrate-version
migrate-version: ## Print current migration version
	$(MIGRATE) version

.PHONY: migrate-force
migrate-force: ## Force set version (recover dirty state): make migrate-force v=2
	$(MIGRATE) force $(v)

.PHONY: migrate-create
migrate-create: ## Scaffold a new migration: make migrate-create name=add_widgets
	$(MIGRATE) create $(name)

.PHONY: migrate-drop
migrate-drop: ## DROP all tables (development only)
	$(MIGRATE) drop

## ---- Seeders ----

.PHONY: seed
seed: ## Apply pending seeders
	$(SEED)

.PHONY: seed-fresh
seed-fresh: ## Clear seed history and re-run all seeders
	$(SEED) -fresh

## ---- Composite (development) ----

.PHONY: db-reset
db-reset: ## DANGER (dev): drop, re-migrate, and re-seed
	$(MIGRATE) drop
	$(MIGRATE) up
	$(SEED)
