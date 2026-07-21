# Automation guide

CRM Lite ships a **Workflow Automation Engine** plus the **Notification Center**.
The legacy `automation_rules` table is deprecated (unused placeholder from early
metadata migrations). Prefer the `workflows*` tables.

Feature flag: `FEATURE_AUTOMATION`. Permissions: `workflow.*` (and
`automation.manage` as a transitional OR on routes).

## Architecture

```text
Record Create/Update/Delete
  → record.Service MutationHook
  → enqueue workflow.evaluate (asynq)
  → worker matches active published workflows
  → evaluate condition tree
  → run actions (notify / records / notes / activities / webhook / delay / invoke)
  → append-only workflow_executions + steps
```

Workspace scope = `organization_id`. Switching workspace isolates workflows and logs.

## Settings UI

| Route | Purpose |
| --- | --- |
| `/settings/automation` | Hub + notification prefs + metrics |
| `/settings/automation/workflows` | List / publish / disable |
| `/settings/automation/workflows/[id]` | Form builder (metadata-driven) |
| `/settings/automation/logs` | Execution history + retry |
| `/settings/automation/templates` | Clone built-in starters |

## API (auth + org + RBAC)

| Method | Path |
| --- | --- |
| GET/POST | `/api/v1/workflows` |
| GET/PATCH/DELETE | `/api/v1/workflows/:id` |
| POST | `/api/v1/workflows/:id/publish` |
| POST | `/api/v1/workflows/:id/disable` |
| GET | `/api/v1/workflows/:id/versions` |
| POST | `/api/v1/workflows/:id/versions/:versionId/rollback` |
| POST | `/api/v1/workflows/:id/run` |
| GET | `/api/v1/workflows/builder-metadata` |
| GET | `/api/v1/workflows/executions` |
| GET | `/api/v1/workflows/executions/:id` |
| POST | `/api/v1/workflows/executions/:id/retry` |
| GET | `/api/v1/workflows/templates` |
| POST | `/api/v1/workflows/templates/:id/clone` |
| GET | `/api/v1/workflows/metrics` |

## Triggers (MVP+)

`record_created`, `record_updated`, `field_updated`, `record_deleted`,
`manual`, `scheduled` (hour/minute UTC + optional `days_of_week` + `batch_size`),
`date_based` (`field_api_name` + `offset_days`).

### Scheduled / date-based sweep

Worker registers `workflow.scheduled_sweep` `@every 1m`:

- **scheduled** — when current UTC hour:minute matches, fan-out to up to
  `batch_size` (default 100) module records; skip if already ran today.
- **date_based** — records where `LEFT(data->>field, 10) = today+offset`;
  skip if already ran today.

### Manual run

`POST /workflows/:id/run` with `{ "record_id", "module_id?" }` from Settings
or the **Run workflow** panel on record detail (`/m/[apiName]/[recordId]`).

### Delay / webhook / resume

- `delay` action enqueues `workflow.resume` via Asynq `ProcessAt`.
- `webhook` POSTs JSON to `config.url` (optional headers/method/payload).

### Soft-delete readiness

`record_deleted` fires from `MutationHook.AfterDelete` with a **before**
snapshot. Records are hard-deleted today; the hook remains valid if soft-delete
is introduced later (fire on status transition or tombstone write).

## Actions (MVP+)

`update_record`, `create_record`, `assign_owner`, `send_email`, `send_whatsapp`,
`create_note`, `create_activity`, `invoke_workflow`, `webhook`, `delay`.

Email/WhatsApp reuse `notification.Compose` — never duplicate providers.

## Organization notification settings

Still stored on `organizations.settings` JSONB:

```json
{
  "automation": {
    "notifications_enabled": true,
    "default_channel": "email",
    "daily_digest": false
  }
}
```

## Versioning

- **Save draft** writes the editable draft version (`workflow_versions.state=draft`).
- **Publish** (with optional changelog) marks the draft published, sets
  `workflows.published_version_id`, and status `active`.
- **Rollback** clones a past published/rolled_back version into a **new**
  published version (history is never rewritten).

UI: workflow editor shows status badges, live version marker, and confirm
dialogs for publish/rollback.

## Template gallery

`GET /workflows/templates` calls `EnsureBuiltinTemplates` so the gallery is
always populated (even without re-seed). Categories: sales, nurture, tasks,
lifecycle — including Lead Qualification, Lost Lead, Birthday, Anniversary,
Task Overdue, Manual Follow-up, and more.

Clone creates a **draft** workflow you can edit and publish.

## Failed-run retry

`POST /workflows/executions/:id/retry` re-queues only `failed` / `partial`
runs when the workflow is still active. Logs UI supports Failed filter and
inline Retry on list rows + detail panel.

## Reentrancy

Workflow-caused record updates set context `source=workflow` and
`exclude_workflow_id`. Max depth is 3.

## Related

- [Architecture](./architecture.md) — queues
- Notification Center routes under `/api/v1/notifications`
- Swagger: `/api/v1/docs`
