-- Phase 13: Export engine.
--
-- export_jobs is the durable record of a (possibly asynchronous) export. The
-- generated file is stored inline in `content` so a completed export is
-- downloadable straight from the database without a separate object store — the
-- worker is stateless with respect to the filesystem, and history/re-download is
-- trivially auditable. export_templates persist reusable export configurations
-- (format + columns + filters + sort) so a common export is one click away.

CREATE TABLE export_jobs (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    module_id        UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    filename         VARCHAR(255) NOT NULL,
    format           VARCHAR(10)  NOT NULL CHECK (format IN ('csv', 'xlsx')),
    status           VARCHAR(20)  NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    columns          JSONB NOT NULL DEFAULT '[]'::jsonb,
    filters          JSONB NOT NULL DEFAULT '[]'::jsonb,
    options          JSONB NOT NULL DEFAULT '{}'::jsonb,
    row_count        INTEGER NOT NULL DEFAULT 0,
    byte_size        INTEGER NOT NULL DEFAULT 0,
    content          BYTEA,
    error            TEXT,
    created_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    started_at       TIMESTAMPTZ,
    finished_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_export_jobs_org_module ON export_jobs(organization_id, module_id, created_at DESC);
CREATE INDEX idx_export_jobs_status     ON export_jobs(status);

CREATE TABLE export_templates (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    module_id        UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    name             VARCHAR(120) NOT NULL,
    format           VARCHAR(10)  NOT NULL DEFAULT 'csv' CHECK (format IN ('csv', 'xlsx')),
    columns          JSONB NOT NULL DEFAULT '[]'::jsonb,
    filters          JSONB NOT NULL DEFAULT '[]'::jsonb,
    sort             JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_export_templates_org_module ON export_templates(organization_id, module_id, name);
