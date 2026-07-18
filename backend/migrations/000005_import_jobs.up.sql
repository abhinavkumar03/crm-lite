-- Phase 12: Import engine.
--
-- import_jobs is the durable record of a bulk import. The uploaded file is
-- parsed at request time and its rows are staged in source_rows (JSONB) so the
-- worker stays stateless with respect to the file — it can process, retry and
-- report entirely from the database. mapping stores the chosen source-header ->
-- field api_name association; errors holds the per-row failure report; the
-- counters make progress observable while the job runs asynchronously.

CREATE TABLE import_jobs (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    module_id        UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    filename         VARCHAR(255) NOT NULL,
    status           VARCHAR(20)  NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    mapping          JSONB NOT NULL DEFAULT '{}'::jsonb,
    options          JSONB NOT NULL DEFAULT '{}'::jsonb,
    source_rows      JSONB NOT NULL DEFAULT '[]'::jsonb,
    total_rows       INTEGER NOT NULL DEFAULT 0,
    processed_rows   INTEGER NOT NULL DEFAULT 0,
    success_rows     INTEGER NOT NULL DEFAULT 0,
    error_rows       INTEGER NOT NULL DEFAULT 0,
    errors           JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    started_at       TIMESTAMPTZ,
    finished_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_import_jobs_org_module ON import_jobs(organization_id, module_id, created_at DESC);
CREATE INDEX idx_import_jobs_status     ON import_jobs(status);
