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

---

## 10. Dynamic Field Engine API

Fields are a sub-resource of a module and are organization-scoped (same auth as
above). Each field response includes a `storage` descriptor — the metadata-driven
persistence strategy telling the client where the value lives (a real column for
native modules, or `data.<api_name>` in the JSONB record store for dynamic ones).
Base URL: `http://localhost:8080/api/v1`.

| Method & path | Feature | Use case |
| --- | --- | --- |
| `GET /modules/:id/fields` | List a module's fields (ordered), with full rendering metadata. | Drive dynamic forms and tables. |
| `POST /modules/:id/fields` | Add a field to a module. | Let admins add a custom attribute (e.g. "Budget"). |
| `GET /modules/:id/fields/:fieldId` | Fetch a single field. | Load a field's edit form. |
| `PUT /modules/:id/fields/:fieldId` | Update a field (labels, validation, options, flags). | Tune a field without recreating it. |
| `POST /modules/:id/fields/reorder` | Batch-update field `sort_order` (atomic). | Persist drag-and-drop field ordering on a form. |
| `DELETE /modules/:id/fields/:fieldId` | Delete a field (system fields are protected → `409`). | Remove a custom field. |

Supported `field_type` values: `text, textarea, email, phone, number, currency,
date, datetime, boolean, dropdown, multiselect, radio, checkbox, url, file,
image, user, lookup, formula, json, richtext`. Choice types (`dropdown`,
`multiselect`, `radio`) require a non-empty `options` list; `lookup` requires a
valid `lookup_module_id`.

**Add a dropdown field:**
```bash
curl -X POST http://localhost:8080/api/v1/modules/<moduleId>/fields \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"api_name":"priority","label":"Priority","field_type":"dropdown","is_required":true,
       "options":[{"label":"Low","value":"low"},{"label":"High","value":"high"}]}'
```

**Add a lookup field (relate to another module):**
```bash
curl -X POST http://localhost:8080/api/v1/modules/<dealModuleId>/fields \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"api_name":"company","label":"Company","field_type":"lookup","lookup_module_id":"<companyModuleId>"}'
```

**List a module's fields (form/table metadata):**
```bash
curl http://localhost:8080/api/v1/modules/<moduleId>/fields -H "Authorization: Bearer $TOKEN"
```

**Reorder fields on a form:**
```bash
curl -X POST http://localhost:8080/api/v1/modules/<moduleId>/fields/reorder \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"items":[{"id":"<field-a>","sort_order":1},{"id":"<field-b>","sort_order":2}]}'
```

---

## 11. Validation Engine API

Database-driven validation rules layered on top of each field's built-in
constraints. The engine merges field metadata (required, length, regex, type
format, options) with rows in `validation_rules` and evaluates a payload,
returning field-keyed errors with custom messages. It also compiles a schema the
frontend can use to mirror validation client-side. Organization-scoped (same
auth as above). Base URL: `http://localhost:8080/api/v1`.

| Method & path | Feature | Use case |
| --- | --- | --- |
| `GET /modules/:id/validation-rules` | List a module's validation rules. | Manage rules in the settings UI. |
| `POST /modules/:id/validation-rules` | Add a rule (field-level or cross-field). | Enforce e.g. "amount ≥ 0" or "email format". |
| `GET /modules/:id/validation-rules/:ruleId` | Fetch a single rule. | Load a rule's edit form. |
| `PUT /modules/:id/validation-rules/:ruleId` | Update a rule (params, message, active, order). | Tweak a rule or toggle it off. |
| `DELETE /modules/:id/validation-rules/:ruleId` | Delete a rule. | Remove a rule. |
| `GET /modules/:id/validation-schema` | Compiled per-field constraints for the client. | Drive frontend (client-side) validation. |
| `POST /modules/:id/validate` | Dry-run validation of a data payload. | Validate a form/record before saving. |

Supported `rule_type` values and their `params`:

| rule_type | Scope | params | Notes |
| --- | --- | --- | --- |
| `required` | field | — | Field must be non-empty. |
| `min_length` / `max_length` | field | `{"value": n}` | String length bounds. |
| `min` / `max` | field | `{"value": n}` | Numeric bounds. |
| `pattern` | field | `{"pattern": "regex"}` | Must match (regex is compiled/validated on save). |
| `email` / `url` | field | — | Format checks. |
| `in` / `not_in` | field | `{"values": [...]}` | Allow/deny list. |
| `required_if` | module | `{"field": "...", "equals": <v>, "target": "..."}` | Cross-field conditional requirement. |

**Add a field-level rule with a custom message:**
```bash
curl -X POST http://localhost:8080/api/v1/modules/<moduleId>/validation-rules \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"field_id":"<amountFieldId>","rule_type":"min","params":{"value":0},"error_message":"Amount cannot be negative"}'
```

**Add a cross-field rule (require company_name when type = company):**
```bash
curl -X POST http://localhost:8080/api/v1/modules/<moduleId>/validation-rules \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"rule_type":"required_if","params":{"field":"type","equals":"company","target":"company_name"}}'
```

**Fetch the compiled schema for the frontend:**
```bash
curl http://localhost:8080/api/v1/modules/<moduleId>/validation-schema -H "Authorization: Bearer $TOKEN"
```

**Dry-run validate a payload:**
```bash
curl -X POST http://localhost:8080/api/v1/modules/<moduleId>/validate \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"data":{"name":"Al","email":"nope","amount":-5}}'
# => {"valid":false,"errors":[{"field":"name","message":"Value is too short"}, ...]}
```

---

## 12. Dynamic Forms (frontend)

Metadata-driven forms rendered entirely from the backend field metadata + the
compiled validation schema — no form is hand-coded. Lives under
`frontend/features/metadata/`.

| Piece | Path | Responsibility |
| --- | --- | --- |
| `DynamicForm` | `features/metadata/components/DynamicForm.tsx` | Renders a full form from `ModuleField[]`, handles layout, visibility, validation, submit. |
| `DynamicField` | `features/metadata/components/DynamicField.tsx` | Maps a single `field_type` to the right shared input primitive. |
| `useDynamicForm` | `features/metadata/hooks/useDynamicForm.ts` | Form state, computed visibility, validation. |
| `lib/conditions.ts` | `features/metadata/lib/conditions.ts` | Conditional-rendering engine (`VisibilityRule[]`). |
| `lib/validation.ts` | `features/metadata/lib/validation.ts` | Client-side validation derived from the Phase 7 schema. |
| `api.ts` | `features/metadata/api.ts` | `getModules`, `getModuleFields`, `getValidationSchema`, `validateRecord`. |
| Playground page | `app/(dashboard)/forms/page.tsx` | Pick a module → generate its form → server-validate. Route: `/forms`. |

Run the frontend and open the page:
```bash
cd frontend
npm run dev
# then visit http://localhost:3000/forms  (also linked in the sidebar)
```

Reuse the renderer for any module:
```tsx
import DynamicForm from "@/features/metadata/components/DynamicForm";
import { getModuleFields, getValidationSchema } from "@/features/metadata/api";

// fields = await getModuleFields(moduleId); schema = await getValidationSchema(moduleId);
<DynamicForm
  fields={fields}
  schema={schema}
  submitText="Save"
  visibilityRules={[
    { when: { field: "type", operator: "equals", value: "company" },
      effect: "show", targets: ["company_name"] },
  ]}
  onSubmit={async (values) => { /* persist via record runtime (Phase 10) */ }}
/>
```

---

## 13. Dynamic Tables & Saved Views

Metadata-driven tables with client-side sorting, filtering and pagination, plus
**saved views** persisted per module in the backend. Columns are derived from
field metadata and each cell is rendered by its `field_type` (badges for choices,
links for urls/emails, formatted dates/currency, check marks for booleans).

### Saved Views API (backend)

Views store a table configuration (visible columns, filter clauses, sort) scoped
to an organization + module. A view is either **public** (shared with the org) or
**private** to its owner; one view per module can be the org **default**.
Organization-scoped (same auth as the other engines).
Base URL: `http://localhost:8080/api/v1`.

| Method & path | Feature | Use case |
| --- | --- | --- |
| `GET /modules/:id/views` | List views visible to the user (public + own), default first. | Populate the view switcher. |
| `POST /modules/:id/views` | Save the current table config as a new view. | "Save view" button. |
| `GET /modules/:id/views/:viewId` | Fetch one view. | Load a view's config. |
| `PUT /modules/:id/views/:viewId` | Update a view (owner only). | Rename / re-configure. |
| `DELETE /modules/:id/views/:viewId` | Delete a view (owner only). | Remove a saved view. |
| `POST /modules/:id/views/:viewId/default` | Make this the module default (clears others). | Star a default view. |

View payload shape (`columns`, `filters`, `sort` are stored as JSONB):

```jsonc
{
  "name": "Open high-value",
  "columns": ["name", "amount", "stage"],
  "filters": [{ "field": "stage", "operator": "equals", "value": "open" }],
  "sort": { "field": "amount", "order": "desc" },
  "is_public": true
}
```

Filter `operator` values: `contains`, `equals`, `not_equals`, `gt`, `lt`, `in`.

**Save a view:**
```bash
curl -X POST http://localhost:8080/api/v1/modules/<moduleId>/views \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"Open high-value","columns":["name","amount","stage"],"filters":[{"field":"stage","operator":"equals","value":"open"}],"sort":{"field":"amount","order":"desc"},"is_public":true}'
```

**List / set default:**
```bash
curl http://localhost:8080/api/v1/modules/<moduleId>/views -H "Authorization: Bearer $TOKEN"
curl -X POST http://localhost:8080/api/v1/modules/<moduleId>/views/<viewId>/default -H "Authorization: Bearer $TOKEN"
```

### Dynamic table (frontend)

Lives under `frontend/features/metadata/`, reusing the Phase 8 form primitives.

| Piece | Path | Responsibility |
| --- | --- | --- |
| `DynamicTable` | `features/metadata/components/DynamicTable.tsx` | Presentational table: sortable headers, per-column filters, pagination. |
| `TableCell` | `features/metadata/components/TableCell.tsx` | Renders one value by `field_type`. |
| `ViewBar` | `features/metadata/components/ViewBar.tsx` | Saved-view switcher: apply / save / default / delete. |
| `useDynamicTable` | `features/metadata/hooks/useDynamicTable.ts` | Owns columns/filters/sort/pagination state and derives visible rows. |
| `lib/table.ts` | `features/metadata/lib/table.ts` | Pure `sortRows` / `filterRows` / `paginate` (type-aware). |
| `api.ts` | `features/metadata/api.ts` | `getViews`, `createView`, `updateView`, `deleteView`, `setDefaultView`. |
| Playground page | `app/(dashboard)/tables/page.tsx` | Pick a dynamic module → table + saved views; create/delete records via the record runtime. Route: `/tables`. |

```bash
cd frontend
npm run dev
# then visit http://localhost:3000/tables  (also linked in the sidebar)
```

> The playground lists dynamic modules only and is now backed by the Phase 10
> record runtime: rows are real records fetched from `GET /modules/:id/records`,
> "Add record" persists via `POST`, and the row trash icon calls `DELETE`. The
> `DynamicTable` component itself was untouched by that swap.

---

## 14. Module Runtime (generic record engine)

The runtime turns any **dynamic** module into a working object: create, read,
update, delete and query records that live entirely in the `records.data` JSONB
column — no per-module tables or migrations. Every write is validated by the
Phase 7 engine, and lookup/user references can be expanded to display labels.
Native modules (leads/contacts/tasks) keep their first-class endpoints and are
rejected here. Organization-scoped. Base URL: `http://localhost:8080/api/v1`.

| Method & path | Feature | Use case |
| --- | --- | --- |
| `GET /modules/:id/records` | Paginated, sortable, filterable, searchable list. | Power a table/list view. |
| `POST /modules/:id/records` | Create a record (validated). | Add a row. |
| `GET /modules/:id/records/:recordId` | Fetch one record (`?expand=true` for relations). | Detail view. |
| `PUT /modules/:id/records/:recordId` | Replace a record's data (validated). | Edit a row. |
| `DELETE /modules/:id/records/:recordId` | Delete a record. | Remove a row. |

**List query parameters:**

| Param | Meaning |
| --- | --- |
| `page`, `page_size` | Pagination (default 20, max 100). |
| `search` | Case-insensitive match across the module's **searchable** fields. |
| `sort`, `order` | Sort by a field api_name (or `created_at`/`updated_at`) + `asc`/`desc`. |
| `expand` | `true` resolves lookup/user fields into `{id, label}` under `relations`. |
| `filter.<field>=<v>` | Shorthand equality filter on a **filterable** field. |
| `filters` | JSON array for richer filters, e.g. `[{"field":"amount","operator":"gt","value":100}]`. |

Filter `operator` values: `eq`, `ne`, `contains`, `gt`, `lt`, `gte`, `lte`, `in`.
Only filterable/searchable fields are honoured and field names are whitelisted
against the module schema, so query input can never widen or inject SQL.

**Create a record:**
```bash
curl -X POST http://localhost:8080/api/v1/modules/<moduleId>/records \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"data":{"name":"Acme Corp","amount":5000,"stage":"open"}}'
# 400 with {"errors":[{"field":"...","message":"..."}]} if validation fails
```

**Query (search + filter + sort + paginate + expand):**
```bash
curl "http://localhost:8080/api/v1/modules/<moduleId>/records?search=acme&sort=amount&order=desc&page=1&page_size=20&expand=true&filters=%5B%7B%22field%22%3A%22amount%22%2C%22operator%22%3A%22gte%22%2C%22value%22%3A1000%7D%5D" \
  -H "Authorization: Bearer $TOKEN"
```

**Update / delete:**
```bash
curl -X PUT http://localhost:8080/api/v1/modules/<moduleId>/records/<recordId> \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"data":{"name":"Acme Corp","amount":7500,"stage":"won"}}'

curl -X DELETE http://localhost:8080/api/v1/modules/<moduleId>/records/<recordId> \
  -H "Authorization: Bearer $TOKEN"
```

Implementation (`backend/internal/record/`, standard vertical slice):

| Piece | Path | Responsibility |
| --- | --- | --- |
| Query builder | `repository/query.go` | Pure, unit-tested SQL builder for search/filter/sort (whitelisted, parameterized). |
| Repository | `repository/record_repository.go` | CRUD, dynamic `List`, batched display-label lookups. |
| Service | `service/record_service.go` + `expand.go` | Validation, ownership stamping, relationship expansion. |
| Handler | `handler/record_handler.go` | Query-string parsing + error mapping. |
| Composition | `record.go` | Reuses the field repo (metadata) and validation service (Phase 7). |

---

## 15. WhatsApp Automation & Notification Pipeline

An end-to-end, provider-agnostic notification pipeline:
**API → `Producer.Publish` → Redis/asynq → worker → `Dispatcher` → `Provider.Send`
→ status + activity**. Messages are persisted in the `notifications` table with a
`queued → sent / failed` lifecycle, so delivery is auditable and retried
independently of the request. WhatsApp delivery is config-driven: the
**simulation** provider (no network) by default, or the **Meta Cloud API** when
credentials are supplied. Organization-scoped.

### Provider configuration (env)

| Variable | Default | Meaning |
| --- | --- | --- |
| `WHATSAPP_PROVIDER` | `simulation` | `simulation` (offline) or `meta` (Meta Cloud API). |
| `WHATSAPP_API_URL` | `https://graph.facebook.com/v20.0` | Graph API base URL. |
| `WHATSAPP_TOKEN` | — | Bearer token (required for `meta`). |
| `WHATSAPP_PHONE_ID` | — | Sender phone number id (required for `meta`). |

If `meta` is selected without full credentials, the pipeline safely falls back to
the simulation provider — the app is always functional out of the box.

### Notification API

Base URL: `http://localhost:8080/api/v1`.

| Method & path | Feature | Use case |
| --- | --- | --- |
| `POST /notifications` | Render, persist (queued) and enqueue a message. | Send WhatsApp/email. |
| `GET /notifications` | List with `status` / `channel` filters + pagination. | Delivery log. |
| `GET /notifications/:id` | Fetch one notification. | Inspect status/error. |

`Subject`/`Body` support `{{placeholder}}` tokens resolved from `data` at send
time. `channel` is `email` or `whatsapp`.

**Send a WhatsApp message (rendered from data):**
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"channel":"whatsapp","to":"+15551234567","template":"lead_welcome","body":"Hi {{name}}, thanks for your interest!","data":{"name":"Dana"}}'
# => { "id": "...", "status": "queued", "body": "Hi Dana, thanks for your interest!" }
```

**Read the delivery log:**
```bash
curl "http://localhost:8080/api/v1/notifications?status=sent&channel=whatsapp" \
  -H "Authorization: Bearer $TOKEN"
```

### Running the pipeline

```bash
# 1) start the worker so queued messages get delivered + status-updated
make run-worker            # go run ./cmd/worker (now also connects to Postgres)

# 2) start the API and POST /notifications; watch the worker log the dispatch,
#    then the notification flips queued -> sent and an activity is recorded.
```

Implementation:

| Piece | Path | Responsibility |
| --- | --- | --- |
| Provider interface + dispatcher | `internal/notify/{notify,dispatcher}.go` | Channel-agnostic Strategy pipeline + provider-name lookup. |
| Template rendering | `internal/notify/render.go` | Dependency-free `{{token}}` substitution. |
| Providers | `internal/notify/{simulation,whatsapp_meta,factory}.go` | Simulation + Meta Cloud API; config-driven selection. |
| Queue | `internal/jobs/{jobs,producer,server}.go` | `notification.send` job type + worker handler + `NotificationProcessor` interface. |
| Engine (API) | `internal/notification/` | Send/list/get vertical slice; renders, persists, enqueues. |
| Processor (worker) | `internal/notification/processor/` | Dispatches, transitions status, writes an activity. |
| Frontend | `frontend/features/notifications/` + `app/(dashboard)/notifications/page.tsx` | Compose form + live delivery log. Route: `/notifications`. |

> Activity logging: each delivery writes an `activities` row
> (`WHATSAPP_SENT` / `EMAIL_SENT` / `NOTIFICATION_FAILED`) referencing the
> notification, so sends appear in the audit trail alongside CRM actions.
