# Migration guide

Schema changes ship as numbered SQL pairs under `backend/migrations/`, applied
by `cmd/migrate` (golang-migrate + embedded files + pgx v5).

## Layout

```text
backend/migrations/
  000001_init_schema.up.sql
  000001_init_schema.down.sql
  000002_add_indexes.up.sql
  …
  000009_perf_indexes.up.sql
```

Rules:

- One concern per migration.
- `up` is forward-only and idempotent where practical (`IF NOT EXISTS` for indexes).
- `down` reverses `up` (drops in reverse FK order).
- Never edit an already-applied migration in shared environments — add a new one.

## Commands

Run from `backend/` (reads `backend/.env`).

| Command | Purpose |
| --- | --- |
| `make migrate-up` | Apply all pending |
| `make migrate-down n=1` | Roll back last `n` |
| `make migrate-version` | Show version + dirty flag |
| `make migrate-force v=N` | Set version without running SQL |
| `make migrate-create name=add_widgets` | Scaffold next pair |
| `make migrate-drop` | Drop all tables (dev only) |
| `make db-reset` | drop → up → seed |

Equivalent: `go run ./cmd/migrate up`.

## Version history (summary)

| Version | Name | Adds |
| --- | --- | --- |
| 1 | `init_schema` | users, native CRM, orgs/roles seed tables |
| 2 | `add_indexes` | Hot-path indexes on native tables |
| 3 | `core_metadata` | modules, fields, views, rules, records, templates |
| 4 | `notifications` | notifications table |
| 5 | `import_jobs` | import_jobs |
| 6 | `exports` | export_jobs (+ template shape) |
| 7 | `tour_progress` | richer tour_progress (replaces weaker 000003 shape) |
| 8 | `rbac_access` | role_module_access, role_field_access |
| 9 | `perf_indexes` | composite indexes for records / notifications |

## Local / CI workflow

```bash
# Fresh machine
cp .env.example .env          # fill DB_* 
make migrate-up
make seed

# After pulling a new migration
make migrate-up
make seed                     # only if seeders changed
```

Docker Compose runs `crm-migrate up` before API/worker start.

## Recovering a dirty database

If a migration fails mid-way, golang-migrate sets **dirty = true** and refuses
further ups.

1. Inspect the failed SQL and fix data/schema by hand if needed.
2. `make migrate-force v=<last-good-version>`
3. `make migrate-up`

Example: failed while applying 8 → force to 7 → re-run up.

## Production guardrails

- `migrate drop` refuses when `APP_ENV=production` unless `MIGRATE_ALLOW_DROP=true`.
- Prefer expand/contract for breaking column changes (add column → backfill →
  switch reads → drop old).
- Take a DB snapshot before `migrate-up` on production.

## Seeding vs migrating

| Concern | Tool |
| --- | --- |
| Tables, indexes, constraints | `cmd/migrate` |
| Demo users, permission catalog, sample modules/records | `cmd/seed` |

Seeders are recorded in `schema_seed_history` and are idempotent.
`make seed-fresh` clears history and re-runs all seeders (dev only).

## Authoring a new migration

```bash
make migrate-create name=add_widget_flags
# edit the generated .up.sql / .down.sql
make migrate-up
make migrate-down n=1   # verify down
make migrate-up
```

Checklist:

- [ ] FKs reference existing tables
- [ ] Org-scoped tables include `organization_id`
- [ ] Down drops in safe order
- [ ] Document in this table when the migration lands
- [ ] Update [ERD](./erd.md) if entities change

## Related

- [ERD](./erd.md)
- [Developer onboarding](./developer-onboarding.md)
- [`COMMANDS.md`](../COMMANDS.md) §1–2
