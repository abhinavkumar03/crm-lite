# CRM Lite documentation

Phase 19 — architecture, data model, and operator/developer guides.

| Doc | Audience | Contents |
| --- | --- | --- |
| [Roadmap](./roadmap.md) | Everyone | Tenancy-first product phases + historical mapping |
| [Architecture](./architecture.md) | Engineers | Processes, layers, tenancy, caching, queues |
| [ERD](./erd.md) | Engineers / DBAs | Entity-relationship diagrams by domain |
| [Sequence diagrams](./sequences.md) | Engineers | Auth, record CRUD, import, export, notifications |
| [Migration guide](./migration-guide.md) | Operators | golang-migrate workflow, versions, recovery |
| [Metadata guide](./metadata-guide.md) | Product / engineers | Modules, fields, validation, views, records |
| [Import / export guide](./import-export-guide.md) | Operators / engineers | File pipelines, APIs, worker behaviour |
| [Automation guide](./automation-guide.md) | Operators / engineers | Notifications, channels, settings, rules table |
| [Developer onboarding](./developer-onboarding.md) | New contributors | Local setup, demo login, first PR |

Related:

- [`COMMANDS.md`](../COMMANDS.md) — make targets, curl recipes, phase notes
- [`CONTRIBUTING.md`](../CONTRIBUTING.md) — branch / commit conventions
- OpenAPI / Swagger UI — `/api/v1/docs` on your API host (Phase 18)
- **How it works (UI)** — interactive guides at [`/help`](../frontend/app/help/page.tsx)
  (`frontend/components/help/`, content in `frontend/features/docs/content.ts`)
