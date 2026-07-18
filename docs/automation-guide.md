# Automation guide

Automation in CRM Lite today means **notification delivery** plus **org-level
automation settings**. A richer `automation_rules` table exists for trigger /
condition / action JSON and is intended for rule-engine expansion; the shipped
path is the notification pipeline.

Feature flag: `FEATURE_AUTOMATION`. Permission: `automation.manage`.

## Organization settings

Stored on `organizations.settings` JSONB (no separate settings table). API:

| Method | Path | Notes |
| --- | --- | --- |
| `GET` | `/api/v1/settings` | Any org member |
| `PUT` | `/api/v1/settings` | Requires `settings.manage` |

Shape (relevant slice):

```json
{
  "automation": {
    "notifications_enabled": true,
    "default_channel": "email",
    "daily_digest": false
  },
  "general": {
    "timezone": "Asia/Kolkata",
    "currency": "INR",
    "locale": "en-IN"
  }
}
```

UI: `/settings` (General) and `/settings/automation`.

## Notification pipeline

```text
API Send → INSERT notifications (queued) → asynq notification.send (critical)
        → worker Dispatcher → MarkSent | MarkFailed → activity log
```

### API

| Method | Path |
| --- | --- |
| `POST` | `/api/v1/notifications` |
| `GET` | `/api/v1/notifications` |
| `GET` | `/api/v1/notifications/:id` |

Create body example:

```json
{
  "channel": "email",
  "to": "ada@example.com",
  "subject": "Welcome",
  "body": "Thanks for joining",
  "template": "lead_welcome",
  "entity_type": "LEAD",
  "entity_id": "…",
  "data": { "name": "Ada" }
}
```

Channels: `email` | `whatsapp`. Statuses: `queued` → `sent` | `failed`.

### Providers (`internal/notify`)

| Channel | Default | Production switch |
| --- | --- | --- |
| Email | Simulation logger | Swap provider registration in worker |
| WhatsApp | Simulation | `WHATSAPP_PROVIDER=meta` + `WHATSAPP_TOKEN` + `WHATSAPP_PHONE_ID` (+ optional `WHATSAPP_API_URL`) |

The dispatcher is strategy-based: handlers never talk to Meta/SMTP directly.

### Lead side-effects

Creating a lead may also enqueue:

- `lead.created` (default queue) — logged by worker
- `email.send` (critical) — if the lead has an email

These share the same notify dispatcher.

### Queue profile

`notification.send` / `email.send` / `whatsapp.send` → **critical** queue,
MaxRetry 5, Timeout 30s.

## `automation_rules` table

Schema (from core metadata migration):

| Column | Role |
| --- | --- |
| `trigger_event` | e.g. record created / status changed |
| `conditions` | JSONB predicate tree |
| `actions` | JSONB action list (notify, assign, …) |
| `module_id` | Optional scope |
| `is_active` | Toggle |

Seed / Settings may expose toggles; a full rules engine evaluator is the natural
next step. Prefer storing durable rule definitions here rather than hardcoding
workflows in handlers.

## Operational checklist

1. Ensure worker is running (`make run-worker`).
2. Confirm Redis connectivity (jobs sit in asynq otherwise).
3. `notifications_enabled` true in settings when UI should offer send.
4. For real WhatsApp: set Meta env vars and restart worker.
5. Inspect `notifications` rows + worker logs for `failed` + `error`.
6. Use Swagger (`/api/v1/docs`) tag **Notifications** to Try it out.

## Related

- [Sequences](./sequences.md) — notification sequence
- [Architecture](./architecture.md) — queues
- Settings Center notes in [`COMMANDS.md`](../COMMANDS.md)
