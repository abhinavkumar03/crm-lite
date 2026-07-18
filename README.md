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

- Auth (JWT), leads / contacts / tasks, dashboard & search
- Metadata engine: modules, fields, validation, views, dynamic records
- Import / export (CSV & XLSX) via background workers
- Notifications (email / WhatsApp providers)
- Settings center, guided tour, RBAC (permissions + module/field ACL)
- OpenAPI 3 + Swagger UI at `/api/v1/docs`
- Redis caching, composite indexes, bulk import optimizations

## Documentation

Start here: **[docs/README.md](./docs/README.md)**

| Doc | Topic |
| --- | --- |
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

Demo admin: `admin@crmlite.com` / `Admin@12345` (after seed).

Hybrid (local API + Docker Postgres/Redis): see
[developer onboarding](./docs/developer-onboarding.md).

## Project structure

```text
crm-lite/
├── backend/          # API, worker, migrate, seed
├── frontend/         # Next.js App Router
├── docs/             # Architecture & guides (Phase 19)
├── deploy/           # K8s / Render
├── docker-compose.yml
└── COMMANDS.md
```

## Roadmap (high level)

- [x] Repository, auth, native CRM, dashboard
- [x] Metadata engine, validation, views, records
- [x] Notifications, import, export
- [x] Guided tour, settings, RBAC
- [x] Performance & caching
- [x] OpenAPI / Swagger
- [x] Documentation & architecture
- [ ] Hardening / production polish as needed

## License

MIT
