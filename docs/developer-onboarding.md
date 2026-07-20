# Developer onboarding

Get from zero to a running CRM Lite stack and a first meaningful change.

## What you are joining

Portfolio CRM: Go/Gin API + asynq worker + PostgreSQL + Redis + Next.js.
Metadata-driven dynamic modules (module listings + Form Designer preview) are the entire CRM surface.

Read next (order matters):

1. [Architecture](./architecture.md) — processes and request path
2. [ERD](./erd.md) — data model map
3. [Metadata guide](./metadata-guide.md) — how custom objects work
4. [`COMMANDS.md`](../COMMANDS.md) — everyday make/curl recipes
5. OpenAPI — `/api/v1/docs` on your API host

## Prerequisites

- Go **1.26+** (see `backend/go.mod`)
- Node **20+** / npm
- Docker Desktop (optional but easiest for Postgres + Redis)
- `make`, `git`

## Option A — Docker Compose (fastest)

From repo root:

```bash
cp backend/.env.example backend/.env
# edit JWT_SECRET and DB/Redis if needed

docker compose up --build
```

Services: Postgres 17, Redis 7, migrate, API `:8080`, worker, frontend `:3000`.

API health: `curl http://localhost:8080/api/v1/health`  
Swagger: `http://localhost:8080/api/v1/docs`  
App: `http://localhost:3000`

## Option B — Hybrid (recommended for backend work)

```bash
# Infra only
docker compose up -d postgres redis

cd backend
cp .env.example .env
make migrate-up
make seed

# Terminal 1
make run-api

# Terminal 2
make run-worker

# Terminal 3
cd ../frontend
cp .env.local.example .env.local 2>/dev/null || true
echo 'NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1' > .env.local
npm install
npm run dev
```

## Demo logins (after seed)

| Email | Password | Role |
| --- | --- | --- |
| `demo@crmlite.com` | `Password@123` | Owner of all 5 demo workspaces |
| `admin@crm.com` | `Admin@123` | Owner on CRM Lite Demo |
| `priya@crmlite.com` | `Password@123` | Sales Manager |
| `vikram@crmlite.com` | `Password@123` | Sales Executive |
| `arjun@crmlite.com` | `Password@123` | Viewer |

Switch workspaces from the sidebar switcher (no re-login required).

## First 30 minutes

1. Login as `demo@crmlite.com` → switch workspaces → open Dashboard and module pages.
2. Hit Swagger → Authorize with a token from `POST /auth/login`.
3. Walk Settings → Modules / Fields / Validation / Roles.
4. Run an import on a dynamic module (`/imports`) with worker running.
5. Skim [sequence diagrams](./sequences.md) for import + notifications.

## Repo map

```text
crm-lite/
├── backend/
│   ├── cmd/api|worker|migrate|seed
│   ├── internal/<feature>/   # vertical slices
│   ├── migrations/
│   └── Makefile
├── frontend/
│   ├── app/(dashboard)/      # routes
│   └── features/             # API clients + domain UI
├── docs/                     # you are here
├── deploy/                   # k8s / render assets
└── docker-compose.yml
```

Backend pattern: **handler → service → repository**, wired in `module.go`.
Do not put SQL in handlers or HTTP details in repositories.

## Making a change

Branch naming / commits: see [`CONTRIBUTING.md`](../CONTRIBUTING.md).

```bash
git checkout -b feature/my-change

# backend
cd backend && go test ./internal/<pkg>/...
gofmt -w .

# frontend
cd frontend && npm run lint
```

Suggested first tasks:

- Add a field validation rule type + test in `validationengine`
- Extend OpenAPI (`make openapi`) when you add a route
- Document non-obvious behaviour in the matching guide under `docs/`

## Environment cheat sheet

Required locally: `DB_*`, `REDIS_*`, `JWT_SECRET`, `FRONTEND_URLS`,
`NEXT_PUBLIC_API_URL`.

Optional: Cloudinary (`CLOUDINARY_*`), WhatsApp Meta (`WHATSAPP_*`), feature
flags `FEATURE_*`.

Never commit real secrets. `.env` / `.env.local` stay gitignored.

## Troubleshooting

| Symptom | Check |
| --- | --- |
| API cannot start | Postgres up? `DB_PASSWORD`? |
| Import stuck `pending` | Worker running? Redis URL? |
| 403 on modules | Role permissions / seed roles |
| Dirty migration | [Migration guide](./migration-guide.md) recovery |
| CORS errors | `FRONTEND_URLS` includes `http://localhost:3000` |
| Stale field list in UI | Settings mutate → metadata cache invalidate; or wait 60s |

## Definition of done (team norm)

- [ ] Tests for new service logic
- [ ] Migration pair if schema changed (+ ERD note)
- [ ] OpenAPI regenerated if routes/DTOs changed (`make openapi`)
- [ ] Guide updated if behaviour is operator-visible
- [ ] `go test ./...` (touched packages at minimum) green

Welcome aboard.
