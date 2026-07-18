# CRM Lite v2 — Command Reference

A single reference for every command in the project: what it does (feature) and
when to reach for it (use case). Commands are grouped by area. This file grows as
new phases add capabilities.

> Backend commands run from the `backend/` directory (they read `backend/.env`).
> Docker Compose commands run from the repository root.

---

## 1. Make targets (backend)

Run `make help` from `backend/` to see this list in your terminal. Each target
wraps a Go command so the everyday workflow is short and memorable.

| Command | Feature | Use case |
|---|---|---|
| `make help` | Lists all self-documented targets | You forgot a target name |
| `make build` | `go build ./...` — compiles every package/binary | Verify the whole backend compiles before committing |
| `make tidy` | `go mod tidy` — syncs `go.mod`/`go.sum` | After adding/removing a dependency |
| `make vet` | `go vet ./...` — static analysis | Catch suspicious constructs before CI |
| `make run-api` | Starts the HTTP API (`cmd/api`) | Local development of the API |
| `make run-worker` | Starts the async worker (`cmd/worker`) | Local development of background jobs/notifications |
| `make migrate-up` | Applies all pending migrations | Bring a database schema up to date |
| `make migrate-down n=1` | Rolls back the last `n` migrations | Undo a migration during development |
| `make migrate-version` | Prints current version + dirty flag | Check what schema version a DB is on |
| `make migrate-force v=2` | Force-sets the version (no SQL run) | Recover from a "dirty" state or baseline a brownfield DB |
| `make migrate-create name=add_widgets` | Scaffolds the next `up`/`down` file pair | Start writing a new migration |
| `make migrate-drop` | Drops all tables (dev only; blocked in prod) | Wipe a local DB before a clean rebuild |
| `make seed` | Runs pending seeders | Populate a fresh DB with baseline/demo data |
| `make seed-fresh` | Clears seed history and re-runs all seeders | Force re-seed after changing seeders |
| `make db-reset` | `drop` → `up` → `seed` (dev only) | One command to rebuild a local DB from scratch |

---

## 2. Migration runner (`cmd/migrate`)

The runner is `go run ./cmd/migrate <command>` (or the built `./crm-migrate`).
It uses golang-migrate with **embedded** SQL files and the **pgx v5** driver, so
it is self-contained and safe to run as a Kubernetes init container or a Render
pre-deploy hook.

| Command | Feature | Use case |
|---|---|---|
| `migrate up` | Apply every pending migration in order | Deploys, CI, first-time setup |
| `migrate down [n]` | Roll back the last `n` migrations (default 1) | Reverting a change locally |
| `migrate version` | Show current version and dirty state | Debugging schema drift |
| `migrate force <version>` | Set the version without running SQL | Fix a dirty state; baseline an existing DB |
| `migrate drop` | Drop all tables | Local reset (refuses in production) |
| `migrate create <name>` | Generate `NNNNNN_<name>.{up,down}.sql` | Author a new migration |

**Guard:** `migrate drop` refuses to run when `APP_ENV=production` unless
`MIGRATE_ALLOW_DROP=true` is set.

**Recover a dirty database:** fix the offending SQL, then
`migrate force <last-good-version>` and re-run `migrate up`.

---

## 3. Seed runner (`cmd/seed`)

`go run ./cmd/seed [flags]` (or the built `./crm-seed`). Seeders are ordered,
recorded in `schema_seed_history`, and idempotent.

| Command | Feature | Use case |
|---|---|---|
| `seed` | Run seeders not yet in the history | Normal seeding on a fresh/updated DB |
| `seed -fresh` | Clear history, then re-run all seeders | After editing seeders, to re-apply everything |

**Demo login credentials** (created by the seeders):

| Email | Password | Role |
|---|---|---|
| `admin@crmlite.com` | `Admin@12345` | Administrator |
| `priya@crmlite.com` | `Password@123` | Sales Manager |
| `vikram@crmlite.com` | `Password@123` | Sales Representative |
| `sneha@crmlite.com` | `Password@123` | Sales Representative |
| `arjun@crmlite.com` | `Password@123` | Viewer |

The demo dataset also creates 50 leads, 30 contacts, 60 tasks, ~100 activities,
notes, and 35 dynamic `company`/`deal` records (JSONB engine) owned by the admin.

---

## 4. Application binaries

| Command | Feature | Use case |
|---|---|---|
| `go run ./cmd/api` | HTTP API server (Gin) | Serve REST endpoints |
| `go run ./cmd/worker` | Async worker (asynq) | Process jobs + notifications (email/WhatsApp pipeline) |
| `go run ./cmd/migrate ...` | Migration runner | See §2 |
| `go run ./cmd/seed ...` | Seed runner | See §3 |

> The API and worker are **separate processes** so they can scale
> independently. Both need Postgres; the worker (and API producer) need Redis.

---

## 5. Docker & Docker Compose (repository root)

| Command | Feature | Use case |
|---|---|---|
| `docker compose up --build` | Builds and starts postgres, redis, migrate (one-shot), backend, worker, frontend | Full local stack with migrations auto-applied |
| `docker compose up -d` | Start in the background | Run the stack without holding the terminal |
| `docker compose logs -f backend` | Follow a service's logs | Debug a specific service |
| `docker compose ps` | Show service status | Check what's running/healthy |
| `docker compose down` | Stop and remove containers | Tear down the stack |
| `docker compose down -v` | Also remove volumes (DB data) | Wipe local database state entirely |
| `docker compose run --rm migrate ./crm-migrate version` | One-off command in the image | Check schema version inside the container network |

**Compose services:** `postgres`, `redis`, `migrate` (runs `crm-migrate up`
once, then exits), `backend` (depends on `migrate`), `worker`, `frontend`.

---

## 6. Go quality & dependencies (backend)

| Command | Feature | Use case |
|---|---|---|
| `go build ./...` | Compile everything | Sanity check after edits |
| `go vet ./...` | Static analysis | Pre-commit checks |
| `go test ./...` | Run tests | Verify behaviour (tests added in later phases) |
| `go mod tidy` | Prune/sync modules | After dependency changes |
| `go get <pkg>@<ver>` | Add/upgrade a dependency | Introduce a new library |

---

## 7. Frontend (from `frontend/`)

| Command | Feature | Use case |
|---|---|---|
| `npm install` | Install dependencies | First checkout / lockfile change |
| `npm run dev` | Start the Next.js dev server | Local UI development |
| `npm run build` | Production build | Verify the app builds; prepare deploy |
| `npm run start` | Serve the production build | Run the built app locally |
| `npm run lint` | ESLint | Catch lint issues |
| `npx tsc --noEmit` | Type-check without emitting | Verify TypeScript types |

---

## 8. Typical workflows

**First-time local setup (Docker):**
```bash
docker compose up --build           # migrations auto-run before the API starts
```

**First-time local setup (without Docker):**
```bash
cd backend
make migrate-up
make seed
make run-api          # terminal 1
make run-worker       # terminal 2
```

**Rebuild a local database from scratch:**
```bash
cd backend && make db-reset
```

**Add a schema change:**
```bash
cd backend
make migrate-create name=add_deals_module
# edit the generated up/down SQL
make migrate-up
```

**Point at an existing (brownfield) database that already has tables:**
```bash
cd backend
make migrate-force v=2   # baseline to the last version its schema matches
make migrate-up          # apply the rest
```

---

## 9. Dynamic Module Engine API

All endpoints are organization-scoped: they require a valid JWT (`Authorization: Bearer <token>`)
and resolve the caller's organization automatically. Base URL: `http://localhost:8080/api/v1`.

| Method & path | Feature | Use case |
| --- | --- | --- |
| `GET /modules` | List every module for the org (enabled or not), ordered by `sort_order`. | Render the module-management/settings screen. |
| `GET /navigation` | List only enabled, sidebar-visible modules. | Build the app sidebar dynamically. |
| `POST /modules` | Create a new dynamic module. | Let admins add a custom object type (e.g. "Projects"). |
| `GET /modules/:id` | Fetch a single module. | Load a module's detail/edit form. |
| `PUT /modules/:id` | Update labels, icon, color, sidebar visibility, default sort. | Rename/re-style a module. |
| `PATCH /modules/:id/status` | Enable or disable a module. | Toggle a module on/off without deleting it. |
| `POST /modules/reorder` | Batch-update `sort_order` (atomic). | Persist drag-and-drop navigation ordering. |
| `DELETE /modules/:id` | Delete a module (system modules are protected → `409`). | Remove a custom module and its records. |

**Create a dynamic module:**
```bash
curl -X POST http://localhost:8080/api/v1/modules \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"api_name":"project","singular_label":"Project","plural_label":"Projects","icon":"folder","color":"#2563eb"}'
```

**Fetch the sidebar navigation:**
```bash
curl http://localhost:8080/api/v1/navigation -H "Authorization: Bearer $TOKEN"
```

**Disable a module:**
```bash
curl -X PATCH http://localhost:8080/api/v1/modules/<id>/status \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"enabled":false}'
```

**Reorder navigation (drag-and-drop persistence):**
```bash
curl -X POST http://localhost:8080/api/v1/modules/reorder \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"items":[{"id":"<id-a>","sort_order":1},{"id":"<id-b>","sort_order":2}]}'
```
