# CRM Lite

Production-style, multi-tenant CRM built to demonstrate modern full-stack
engineering: Go/Gin API, asynq worker, PostgreSQL, Redis, Next.js, and a
metadata-driven module engine.

## Tech stack

| Layer | Choices |
| --- | --- |
| Frontend | Next.js, TypeScript, Tailwind CSS |
| Backend | Go, Gin, pgx, validator |
| Data | PostgreSQL 17 |
| Cache / queue | Redis 7, asynq |
| Media | Cloudinary |
| Ship | Docker Compose, Kubernetes, Render |

## Capabilities

- Multi-tenant organizations (isolation, membership, active-org switch)
- Roles, permissions, module/field ACL, reporting hierarchy, departments/teams
- Record visibility engine (owner / manager / hierarchy / org / team / …)
- Metadata engine: modules, fields, validation, views, dynamic records
- Import / export (CSV & XLSX) via background workers
- Notifications, settings center, guided tour + interactive sandbox tutorial (Explore CRM)
- OpenAPI 3 + Swagger UI at `/api/v1/docs`
- Redis caching, composite indexes, bulk import optimizations

## Documentation

Start here: **[docs/README.md](./docs/README.md)** · **[Product roadmap](./docs/roadmap.md)** (tenancy-first)

| Doc | Topic |
| --- | --- |
| [Roadmap](./docs/roadmap.md) | Tenancy → hierarchy → visibility → metadata |
| [Developer onboarding](./docs/developer-onboarding.md) | Local setup & first PR |
| [Architecture](./docs/architecture.md) | Processes, layers, tenancy |
| [ERD](./docs/erd.md) | Data model diagrams |
| [Sequence diagrams](./docs/sequences.md) | Auth, import, export, notify |
| [Migration guide](./docs/migration-guide.md) | Schema workflow |
| [Metadata guide](./docs/metadata-guide.md) | Modules / fields / records |
| [Import / export guide](./docs/import-export-guide.md) | Bulk pipelines |
| [Automation guide](./docs/automation-guide.md) | Notifications & settings |

Also: [`COMMANDS.md`](./COMMANDS.md) (make targets & curl), [`CONTRIBUTING.md`](./CONTRIBUTING.md).

## Quick start

```bash
cp backend/.env.example backend/.env
docker compose up --build
```

- App: http://localhost:3000  
- API: http://localhost:8080/api/v1/health  
- Swagger: http://localhost:8080/api/v1/docs  

Demo login: `demo@crmlite.com` / `Password@123` (owner of 5 workspaces after seed).  
Also available: `admin@crm.com` / `Admin@123` on the primary workspace.

Hybrid (local API + Docker Postgres/Redis): see
[developer onboarding](./docs/developer-onboarding.md).

## Project structure

```text
crm-lite/
├── backend/          # API, worker, migrate, seed
├── frontend/         # Next.js App Router
├── docs/             # Architecture, roadmap & guides
├── deploy/           # K8s / Render
├── docker-compose.yml
└── COMMANDS.md
```

## Roadmap (high level)

Canonical order (see [`docs/roadmap.md`](./docs/roadmap.md)):

- [x] Phase 1 — Multi-tenant foundation (harden org profile)
- [x] Phase 2 — Organization, membership, invites, bootstrap
- [x] Phase 3 — Roles, hierarchy, departments / teams
- [x] Phase 4 — Record visibility engine
- [x] Phase 5 — Tenant-aware seed / demos
- [x] Phase 6+ — Metadata CRM, import/export, notify, OpenAPI, docs
- [ ] Hardening / production polish as needed

## License

MIT
