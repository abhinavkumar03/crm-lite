export type FlowStep = {
  title: string;
  detail: string;
};

export type DocBlock =
  | { type: "intro"; text: string }
  | { type: "bullets"; title?: string; items: string[] }
  | { type: "table"; title?: string; headers: string[]; rows: string[][] }
  | { type: "flow"; title: string; steps: FlowStep[] }
  | { type: "map"; title: string; nodes: { id: string; label: string; hint?: string }[]; edges: string[] }
  | { type: "groups"; title: string; groups: { name: string; items: string[] }[] }
  | { type: "callout"; tone: "tip" | "warn" | "note"; text: string };

export type DocGuide = {
  id: string;
  title: string;
  eyebrow: string;
  summary: string;
  readingTime: string;
  blocks: DocBlock[];
};

/** Swagger UI path on the configured API base (no hardcoded host). */
export function apiSwaggerDocsUrl(): string {
  const base = process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "");
  return base ? `${base}/docs` : "/api/v1/docs";
}

export const DOC_GUIDES: DocGuide[] = [
  {
    id: "architecture",
    title: "Architecture",
    eyebrow: "System design",
    summary:
      "How the API, worker, Postgres, and Redis fit together — and how a request travels through auth, tenancy, and RBAC.",
    readingTime: "4 min",
    blocks: [
      {
        type: "intro",
        text: "CRM Lite is a multi-tenant, metadata-driven CRM: Go/Gin serves HTTP, a separate asynq worker runs bulk and notification jobs, PostgreSQL is the source of truth, and Redis backs cache + queues.",
      },
      {
        type: "map",
        title: "System context",
        nodes: [
          { id: "ui", label: "Next.js", hint: "Browser" },
          { id: "api", label: "cmd/api", hint: "Gin" },
          { id: "worker", label: "cmd/worker", hint: "asynq" },
          { id: "pg", label: "PostgreSQL", hint: "SoR" },
          { id: "redis", label: "Redis", hint: "Cache + queue" },
        ],
        edges: [
          "UI → API (REST + JWT)",
          "API → Postgres",
          "API → Redis (cache / enqueue)",
          "Worker ← Redis → Postgres",
        ],
      },
      {
        type: "table",
        title: "Processes",
        headers: ["Binary", "Role", "Postgres", "Redis"],
        rows: [
          ["cmd/api", "HTTP + enqueue + cache", "Yes", "Yes"],
          ["cmd/worker", "Import / export / notify", "Yes", "Yes"],
          ["cmd/migrate", "Schema migrations", "Yes", "—"],
          ["cmd/seed", "Demo + catalog data", "Yes", "—"],
        ],
      },
      {
        type: "flow",
        title: "Request path",
        steps: [
          { title: "Middleware", detail: "RequestID, logger, recovery, CORS, security headers" },
          { title: "JWT", detail: "Bearer token → userID on context" },
          { title: "Tenant", detail: "Resolve org + role (cached ~2m)" },
          { title: "RBAC", detail: "Load permission keys; Require / RequireModule" },
          { title: "Handler → Service → Repo", detail: "Vertical slice; SQL only in repositories" },
        ],
      },
      {
        type: "bullets",
        title: "Storage strategies",
        items: [
          "Product CRM is dynamic-only — modules + records.data JSONB",
          "Legacy native tables (leads/contacts/tasks) remain in migrations but APIs are unwired",
          "Vertical slices under internal/<feature>/{handler,service,repository}",
        ],
      },
      {
        type: "callout",
        tone: "note",
        text: "API and worker scale independently. Never run migrations from the API process.",
      },
    ],
  },
  {
    id: "erd",
    title: "Entity model (ERD)",
    eyebrow: "Data model",
    summary:
      "Tables grouped by domain — identity, metadata engine, jobs, RBAC, and tour.",
    readingTime: "5 min",
    blocks: [
      {
        type: "intro",
        text: "Schema lives in backend/migrations. Everything org-scoped carries organization_id. Active product data is modules + records.",
      },
      {
        type: "groups",
        title: "Domain clusters",
        groups: [
          {
            name: "Identity & tenancy",
            items: ["users", "organizations", "organization_members", "roles"],
          },
          {
            name: "Legacy tables (unused by API)",
            items: ["leads", "contacts", "tasks", "notes", "call_logs"],
          },
          {
            name: "Metadata engine",
            items: ["modules", "fields", "views", "validation_rules", "layouts", "records"],
          },
          {
            name: "Jobs & notifications",
            items: ["notifications", "import_jobs", "export_jobs", "export_templates"],
          },
          {
            name: "RBAC",
            items: ["permissions", "role_permissions", "role_module_access", "role_field_access"],
          },
          {
            name: "Tour",
            items: ["tour_steps", "tour_progress"],
          },
        ],
      },
      {
        type: "bullets",
        title: "Key relationships",
        items: [
          "organization_members links users ↔ organizations (+ optional role_id)",
          "modules → fields → validation_rules; dynamic rows in records.data (GIN)",
          "role_module_access / role_field_access overlay permissions (absence = unrestricted layer)",
          "notes / call_logs / attachments / activities are polymorphic (entity_type + entity_id)",
        ],
      },
      {
        type: "callout",
        tone: "tip",
        text: "Full Mermaid diagrams are in docs/erd.md — this UI highlights the clusters you’ll touch day to day.",
      },
    ],
  },
  {
    id: "sequences",
    title: "Sequence flows",
    eyebrow: "Runtime",
    summary:
      "Happy-path flows for login, protected requests, records, import, export, and notifications.",
    readingTime: "6 min",
    blocks: [
      {
        type: "flow",
        title: "Login",
        steps: [
          { title: "Credentials", detail: "UI posts email + password" },
          { title: "Verify", detail: "API loads user, bcrypt compare" },
          { title: "JWT", detail: "Sign access_token + user payload" },
          { title: "Session", detail: "Frontend stores token → /dashboard" },
        ],
      },
      {
        type: "flow",
        title: "Protected org request",
        steps: [
          { title: "Bearer JWT", detail: "Validate → userID" },
          { title: "Membership cache", detail: "tenant:membership:{userID} or Postgres" },
          { title: "Permissions cache", detail: "rbac:perms:{roleID} or Postgres" },
          { title: "Authorize", detail: "Require(permission) / module ACL" },
          { title: "Execute", detail: "Handler → service → repository" },
        ],
      },
      {
        type: "flow",
        title: "Import",
        steps: [
          { title: "Analyze", detail: "Multipart file → headers + suggested mapping" },
          { title: "Create job", detail: "Stage source_rows; enqueue import.process (bulk)" },
          { title: "Worker", detail: "LoadSpec once → validate rows → COPY batches" },
          { title: "Finish", detail: "Counters + up to 500 row errors" },
        ],
      },
      {
        type: "flow",
        title: "Async export",
        steps: [
          { title: "Create", detail: "POST export job + enqueue export.process" },
          { title: "Build", detail: "Page records with SkipTotal; write CSV/XLSX" },
          { title: "Store", detail: "Persist BYTEA content on export_jobs" },
          { title: "Download", detail: "GET …/exports/:id/download" },
        ],
      },
      {
        type: "flow",
        title: "Notification",
        steps: [
          { title: "Send", detail: "INSERT notifications status=queued" },
          { title: "Enqueue", detail: "notification.send on critical queue" },
          { title: "Dispatch", detail: "Worker → notify.Dispatcher (email/WhatsApp)" },
          { title: "Settle", detail: "MarkSent or MarkFailed + activity" },
        ],
      },
    ],
  },
  {
    id: "migrations",
    title: "Migration guide",
    eyebrow: "Schema ops",
    summary:
      "How golang-migrate versions work, how to recover a dirty DB, and how to author a safe change.",
    readingTime: "4 min",
    blocks: [
      {
        type: "intro",
        text: "Numbered SQL pairs under backend/migrations are embedded and applied by cmd/migrate. Seeders are separate (cmd/seed).",
      },
      {
        type: "table",
        title: "Everyday commands",
        headers: ["Command", "Purpose"],
        rows: [
          ["make migrate-up", "Apply all pending"],
          ["make migrate-down n=1", "Roll back last n"],
          ["make migrate-version", "Show version + dirty"],
          ["make migrate-force v=N", "Set version without SQL"],
          ["make migrate-create name=…", "Scaffold next pair"],
          ["make db-reset", "drop → up → seed (dev)"],
        ],
      },
      {
        type: "flow",
        title: "Dirty database recovery",
        steps: [
          { title: "Inspect", detail: "Find the failed statement / partial objects" },
          { title: "Repair", detail: "Fix schema/data by hand if needed" },
          { title: "Force", detail: "make migrate-force v=<last-good>" },
          { title: "Retry", detail: "make migrate-up" },
        ],
      },
      {
        type: "table",
        title: "Version map (summary)",
        headers: ["#", "Name", "Adds"],
        rows: [
          ["1", "init_schema", "Users + native CRM"],
          ["3", "core_metadata", "Modules, fields, records"],
          ["4–6", "jobs", "Notifications, import, export"],
          ["7", "tour_progress", "Guided tour"],
          ["8", "rbac_access", "Module/field ACL"],
          ["9", "perf_indexes", "Composite list indexes"],
        ],
      },
      {
        type: "callout",
        tone: "warn",
        text: "Never edit an already-applied migration in shared environments — add a new numbered pair. migrate drop is blocked in production.",
      },
    ],
  },
  {
    id: "metadata",
    title: "Metadata guide",
    eyebrow: "Customization",
    summary:
      "Modules, fields, validation, views, and the dynamic record runtime that powers module pages and Form Designer.",
    readingTime: "5 min",
    blocks: [
      {
        type: "intro",
        text: "Custom objects are described in metadata and executed by engines. Settings UI is a thin admin over those engines.",
      },
      {
        type: "table",
        title: "Core concepts",
        headers: ["Concept", "Table", "Meaning"],
        rows: [
          ["Module", "modules", "Object type (deal, company, …)"],
          ["Field", "fields", "Typed attribute (api_name + field_type)"],
          ["Record", "records", "Dynamic row (data JSONB)"],
          ["Rule", "validation_rules", "Field or module constraint"],
          ["View", "views", "Saved columns / filters / sort"],
        ],
      },
      {
        type: "flow",
        title: "Create a custom object",
        steps: [
          { title: "Module", detail: "Settings → Modules (storage_strategy = dynamic)" },
          { title: "Fields", detail: "Add typed fields + options / lookups" },
          { title: "Validation", detail: "Attach rules (required, min_length, …)" },
          { title: "Use", detail: "Module listings / Form Designer / Import / Export consume the catalog" },
        ],
      },
      {
        type: "bullets",
        title: "Record API highlights",
        items: [
          "GET/POST /modules/:id/records — list & create",
          "Query: page, page_size, search, sort, order, filters, expand",
          "expand=true resolves lookup/user labels into relations",
          "RBAC: record.* plus optional module/field ACL",
        ],
      },
      {
        type: "callout",
        tone: "tip",
        text: "Frontend caches modules/fields/schemas for 60s; Settings mutations invalidate the cache.",
      },
    ],
  },
  {
    id: "import-export",
    title: "Import / export",
    eyebrow: "Bulk data",
    summary:
      "File pipelines for dynamic modules — analyze mapping, worker validation, sync vs async export.",
    readingTime: "5 min",
    blocks: [
      {
        type: "intro",
        text: "Import and export require storage_strategy = dynamic, Redis, and a running worker. Permissions: import.run / export.run.",
      },
      {
        type: "flow",
        title: "Import lifecycle",
        steps: [
          { title: "Analyze", detail: "POST …/imports/analyze → headers + suggested mapping" },
          { title: "Create", detail: "Multipart file + mapping JSON → import_jobs pending" },
          { title: "Process", detail: "Worker: LoadSpec → ValidateWithSpec → COPY ×100" },
          { title: "Result", detail: "completed/failed + row error report" },
        ],
      },
      {
        type: "table",
        title: "Export modes",
        headers: ["Mode", "Endpoint", "When"],
        rows: [
          ["Sync", "GET …/export", "Small sets; stream immediately"],
          ["Async", "POST …/exports", "Larger sets; worker builds file"],
          ["Download", "GET …/exports/:id/download", "After async ready"],
        ],
      },
      {
        type: "bullets",
        title: "Operator tips",
        items: [
          "Prefer UTF-8 CSV; run Analyze before Create",
          "Mapping is header → field api_name",
          "Completed import jobs are never re-run (new upload = new job)",
          "Async export paging skips COUNT(*) (SkipTotal)",
        ],
      },
      {
        type: "callout",
        tone: "warn",
        text: "If jobs stay pending, the worker is not running or Redis is unreachable.",
      },
    ],
  },
  {
    id: "automation",
    title: "Automation",
    eyebrow: "Notifications",
    summary:
      "Org automation settings plus the notification send pipeline (email / WhatsApp providers).",
    readingTime: "4 min",
    blocks: [
      {
        type: "intro",
        text: "Shipped automation centers on notifications and organizations.settings.automation. The automation_rules table is ready for a fuller rules engine.",
      },
      {
        type: "table",
        title: "Settings slice",
        headers: ["Key", "Meaning"],
        rows: [
          ["notifications_enabled", "UI / send affordances"],
          ["default_channel", "email | whatsapp"],
          ["daily_digest", "Digest preference flag"],
        ],
      },
      {
        type: "flow",
        title: "Send notification",
        steps: [
          { title: "API", detail: "POST /notifications → row status=queued" },
          { title: "Queue", detail: "notification.send on critical (retry 5 / 30s)" },
          { title: "Provider", detail: "Dispatcher → simulation or Meta WhatsApp" },
          { title: "Settle", detail: "MarkSent / MarkFailed + activity log" },
        ],
      },
      {
        type: "bullets",
        title: "Providers",
        items: [
          "Email: simulation logger in default worker wiring",
          "WhatsApp: WHATSAPP_PROVIDER=meta + token/phone ID for Meta Cloud API",
          "Notification sends enqueue notification.send on the critical queue"
        ],
      },
      {
        type: "callout",
        tone: "note",
        text: "Permission automation.manage gates notification management; settings.manage gates PUT /settings.",
      },
    ],
  },
  {
    id: "onboarding",
    title: "Developer onboarding",
    eyebrow: "Getting started",
    summary:
      "From zero to a running stack — Compose or hybrid local API — plus demo logins and first-change habits.",
    readingTime: "5 min",
    blocks: [
      {
        type: "flow",
        title: "Hybrid local setup",
        steps: [
          { title: "Infra", detail: "docker compose up -d postgres redis" },
          { title: "Schema", detail: "cd backend && make migrate-up && make seed" },
          { title: "API + worker", detail: "make run-api · make run-worker" },
          { title: "Frontend", detail: "NEXT_PUBLIC_API_URL=… npm run dev" },
        ],
      },
      {
        type: "table",
        title: "Demo logins",
        headers: ["Email", "Password", "Role"],
        rows: [
          ["admin@crm.com", "Admin@123", "Organization Owner"],
          ["priya@crmlite.com", "Password@123", "Sales Manager"],
          ["vikram@crmlite.com", "Password@123", "Sales Rep"],
          ["arjun@crmlite.com", "Password@123", "Viewer"],
        ],
      },
      {
        type: "bullets",
        title: "First 30 minutes",
        items: [
          "Login as admin → Dashboard, module pages, Form Designer, Settings → Modules",
          "Open Swagger at /api/v1/docs and Authorize with a JWT",
          "Walk Settings → Modules / Fields / Validation / Roles",
          "Run an import with the worker running",
        ],
      },
      {
        type: "bullets",
        title: "Definition of done",
        items: [
          "Tests for new service logic",
          "Migration pair + ERD note if schema changed",
          "make openapi if routes/DTOs changed",
          "Update the matching guide under docs/",
        ],
      },
      {
        type: "callout",
        tone: "tip",
        text: `Swagger UI: ${apiSwaggerDocsUrl()} · Markdown source of truth: /docs in the repo.`,
      },
    ],
  },
];
