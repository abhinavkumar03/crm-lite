# Import / export guide

Bulk data moves for **dynamic** modules only (`storage_strategy = dynamic`).

Requires worker process (`make run-worker`) and Redis. Feature flags:
`FEATURE_IMPORT`, `FEATURE_EXPORT`.

## Permissions

| Action | Permission | Module ACL |
| --- | --- | --- |
| Import | `import.run` | `can_create` on the module |
| Export | `export.run` | `can_view` on the module |

## Import

### Lifecycle

```text
analyze (optional) → create job (pending) → worker (processing) → completed|failed
```

| Status | Meaning |
| --- | --- |
| `pending` | Staged; waiting on queue |
| `processing` | Worker owns the job |
| `completed` | Finished (may still have row errors) |
| `failed` | Could not run, or zero successful rows |

### API

Base: `/api/v1/modules/:moduleId/imports`

| Step | Method | Body |
| --- | --- | --- |
| Analyze | `POST …/analyze` | multipart `file` |
| Create | `POST …/` | multipart `file`, `mapping` (JSON string), optional `options` |
| List | `GET …/` | `page`, `page_size`, `status` |
| Get | `GET …/:importId` | — |

**Analyze** returns headers, sample rows, suggested mapping, row count — nothing
persisted.

**Create** stages `source_rows` + mapping on `import_jobs`, enqueues
`import.process` on the **bulk** queue (MaxRetry 3, Timeout 10m).

### Mapping

JSON object: CSV/XLSX **header → field `api_name`**.

```json
{
  "Company Name": "name",
  "Deal Size": "amount",
  "Stage": "stage"
}
```

Empty cells are omitted so required-field validation can fire.

### Worker behaviour (Phase 17)

1. Idempotent skip if already terminal.
2. `LoadSpec` once (fields + active validation rules).
3. Per row: map → coerce types → `ValidateWithSpec`.
4. Flush validated rows with Postgres `COPY` every 100 rows; fall back to
   single inserts if a batch fails.
5. Persist up to 500 row errors in `errors` JSONB; counters reflect true totals.

UI: `/imports`.

### Practical tips

- Prefer UTF-8 CSV for portability.
- Run Analyze first to catch header typos.
- Large files: watch worker logs (`import: completed`).
- Re-upload creates a **new** job; completed jobs are never re-run.

## Export

### Modes

| Mode | Endpoint | When |
| --- | --- | --- |
| Sync | `GET /modules/:id/export` | Small sets; streams file immediately |
| Async | `POST /modules/:id/exports` | Larger sets; worker builds file |

Async statuses: `pending` → `processing` → ready/completed → download via
`GET …/exports/:exportId/download`. Failed jobs store `error` text.

### Spec

```json
{
  "format": "csv",
  "columns": ["name", "amount", "stage"],
  "filters": [{ "field": "stage", "operator": "eq", "value": "proposal" }],
  "search": "acme",
  "sort": "amount",
  "order": "desc",
  "expand": false
}
```

Formats: `csv`, `xlsx`. Max rows per build is capped (see exporter service
`MaxRows`). Async paging uses `SkipTotal` so each page skips `COUNT(*)`.

### Templates

CRUD under `/modules/:id/export-templates` — named presets for columns/filters/
format reused by the UI (`/exports`).

### Queue

`export.process` → **bulk** queue (same retry/timeout profile as import).

## End-to-end checklist

1. `make migrate-up && make seed`
2. Start API + worker
3. Login as admin → open a dynamic module (e.g. Deals)
4. Import: analyze sample file → map → create → poll job until `completed`
5. Export: sync download for a quick check; async for larger filters
6. Confirm RBAC: viewer without `import.run` receives 403

## Related

- [Sequences](./sequences.md) — import & export diagrams
- [Metadata guide](./metadata-guide.md)
- [`COMMANDS.md`](../COMMANDS.md) import/export sections
- OpenAPI tags **Import** / **Export**
