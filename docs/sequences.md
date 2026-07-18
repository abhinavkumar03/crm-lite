# Sequence diagrams

Happy-path flows for the critical subsystems. Error branches follow the shared
API envelope (`success: false`, optional `errors`).

## Login

```mermaid
sequenceDiagram
  actor User
  participant UI as Next.js
  participant API as cmd/api
  participant DB as PostgreSQL

  User->>UI: email + password
  UI->>API: POST /api/v1/auth/login
  API->>DB: SELECT user by email
  API->>API: bcrypt compare + sign JWT
  API-->>UI: { access_token, user }
  UI->>UI: store token, redirect /dashboard
```

## Authenticated org-scoped request

```mermaid
sequenceDiagram
  participant UI as Next.js
  participant API as Gin
  participant Cache as Redis
  participant DB as PostgreSQL

  UI->>API: Bearer JWT + request
  API->>API: validate JWT → userID
  API->>Cache: GET tenant:membership:{userID}
  alt cache miss
    API->>DB: organization_members + roles
    API->>Cache: SET membership TTL 2m
  end
  API->>Cache: GET rbac:perms:{roleID}
  alt cache miss
    API->>DB: role_permissions
    API->>Cache: SET perms TTL 2m
  end
  API->>API: Require(permission) / RequireModule(action)
  API->>DB: handler → service → repository
  API-->>UI: APIResponse
```

## Dynamic record create (with validation)

```mermaid
sequenceDiagram
  participant UI as Forms / Tables
  participant API as record handler
  participant VE as validationengine
  participant DB as PostgreSQL

  UI->>API: POST /modules/{id}/records { data }
  API->>API: rbac record.create + module ACL + field ACL
  API->>VE: Validate(org, module, data)
  VE->>DB: load fields + active rules
  VE-->>API: Valid / FieldError[]
  alt invalid
    API-->>UI: 400 Validation failed
  else valid
    API->>DB: INSERT records
    API-->>UI: 201 RecordResponse
  end
```

## Import (analyze → create → worker)

```mermaid
sequenceDiagram
  actor User
  participant UI as Imports page
  participant API as cmd/api
  participant Q as Redis asynq
  participant W as cmd/worker
  participant DB as PostgreSQL

  User->>UI: upload CSV/XLSX
  UI->>API: POST .../imports/analyze (multipart)
  API-->>UI: headers, sample_rows, suggested_mapping

  User->>UI: confirm mapping
  UI->>API: POST .../imports (file + mapping JSON)
  API->>DB: INSERT import_jobs pending + source_rows
  API->>Q: Enqueue import.process (bulk)
  API-->>UI: ImportResponse pending

  Q->>W: deliver job
  W->>DB: MarkProcessing
  W->>W: LoadSpec once (fields + rules)
  loop each row (batch COPY every 100)
    W->>W: map → ValidateWithSpec
    W->>DB: CreateBatch records
  end
  W->>DB: Finish completed/failed + errors
```

## Export (async)

```mermaid
sequenceDiagram
  participant UI as Exports page
  participant API as cmd/api
  participant Q as Redis asynq
  participant W as cmd/worker
  participant DB as PostgreSQL

  UI->>API: POST /modules/{id}/exports { format, columns, filters }
  API->>DB: INSERT export_jobs pending
  API->>Q: Enqueue export.process (bulk)
  API-->>UI: ExportResponse pending

  Q->>W: deliver job
  W->>DB: page records SkipTotal (no COUNT)
  W->>W: serialize CSV/XLSX
  W->>DB: store content BYTEA + status ready
  UI->>API: GET .../exports/{id}/download
  API->>DB: read content
  API-->>UI: file bytes
```

Sync alternative: `GET /modules/{id}/export` builds and streams immediately
(no job row).

## Notification send

```mermaid
sequenceDiagram
  participant UI as Notifications
  participant API as cmd/api
  participant Q as Redis asynq
  participant W as cmd/worker
  participant Disp as notify.Dispatcher
  participant DB as PostgreSQL

  UI->>API: POST /notifications { channel, to, body, ... }
  API->>DB: INSERT notifications status=queued
  API->>Q: Enqueue notification.send (critical)
  API-->>UI: NotificationResponse queued

  Q->>W: deliver job
  W->>DB: load notification
  W->>Disp: Dispatch(channel, message)
  alt success
    W->>DB: MarkSent + activity
  else failure
    W->>DB: MarkFailed + error
  end
```

Providers: email is simulation in the default worker wiring; WhatsApp uses Meta
Cloud API when `WHATSAPP_PROVIDER=meta` and credentials are set, otherwise
simulation.

## Lead create → dashboard invalidate

```mermaid
sequenceDiagram
  participant UI as Leads
  participant API as lead service
  participant Q as Redis asynq
  participant Cache as Redis cache
  participant DB as PostgreSQL

  UI->>API: POST /leads
  API->>DB: INSERT leads
  API->>DB: activity LEAD_CREATED
  API->>Q: lead.created (+ optional email.send)
  API->>Cache: DEL dashboard:{ownerID}
  API-->>UI: LeadResponse
```

## Related

- [Architecture](./architecture.md)
- [Import / export guide](./import-export-guide.md)
- [Automation guide](./automation-guide.md)
