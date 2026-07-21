-- Workflow Automation Engine (Phase A).
-- Supersedes the unused automation_rules placeholder with a normalized schema.

CREATE TABLE workflows (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id   UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    module_id         UUID REFERENCES modules(id) ON DELETE CASCADE,
    name              VARCHAR(150) NOT NULL,
    description       TEXT NOT NULL DEFAULT '',
    status            VARCHAR(20) NOT NULL DEFAULT 'draft'
                          CHECK (status IN ('draft', 'active', 'disabled', 'archived')),
    on_action_error   VARCHAR(20) NOT NULL DEFAULT 'stop'
                          CHECK (on_action_error IN ('continue', 'stop')),
    priority          INT NOT NULL DEFAULT 100,
    published_version_id UUID,
    created_by        UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by        UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE workflow_versions (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id          UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    organization_id      UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    version              INT NOT NULL,
    state                VARCHAR(20) NOT NULL DEFAULT 'draft'
                             CHECK (state IN ('draft', 'published', 'rolled_back')),
    definition_snapshot  JSONB NOT NULL DEFAULT '{}'::jsonb,
    changelog            TEXT NOT NULL DEFAULT '',
    published_at         TIMESTAMPTZ,
    published_by         UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (workflow_id, version)
);

ALTER TABLE workflows
    ADD CONSTRAINT workflows_published_version_fk
    FOREIGN KEY (published_version_id) REFERENCES workflow_versions(id) ON DELETE SET NULL;

CREATE TABLE workflow_triggers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version_id      UUID NOT NULL REFERENCES workflow_versions(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    type            VARCHAR(40) NOT NULL
                        CHECK (type IN (
                            'record_created', 'record_updated', 'record_deleted',
                            'field_updated', 'scheduled', 'date_based', 'manual'
                        )),
    config          JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE workflow_conditions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version_id      UUID NOT NULL REFERENCES workflow_versions(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    parent_id       UUID REFERENCES workflow_conditions(id) ON DELETE CASCADE,
    node_type       VARCHAR(20) NOT NULL CHECK (node_type IN ('group', 'predicate')),
    logic           VARCHAR(10) CHECK (logic IS NULL OR logic IN ('and', 'or')),
    field_api_name  VARCHAR(100),
    operator        VARCHAR(40),
    value           JSONB,
    sort_order      INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE workflow_actions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version_id        UUID NOT NULL REFERENCES workflow_versions(id) ON DELETE CASCADE,
    organization_id   UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    sort_order        INT NOT NULL DEFAULT 0,
    type              VARCHAR(40) NOT NULL
                          CHECK (type IN (
                              'update_record', 'create_record', 'delete_record',
                              'assign_owner', 'send_email', 'send_whatsapp',
                              'create_note', 'create_activity', 'webhook',
                              'delay', 'invoke_workflow', 'branch'
                          )),
    config            JSONB NOT NULL DEFAULT '{}'::jsonb,
    max_retries       INT NOT NULL DEFAULT 0,
    continue_on_error BOOLEAN,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE workflow_executions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    workflow_id     UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    version_id      UUID REFERENCES workflow_versions(id) ON DELETE SET NULL,
    module_id       UUID REFERENCES modules(id) ON DELETE SET NULL,
    record_id       UUID,
    trigger_type    VARCHAR(40) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'queued'
                        CHECK (status IN ('queued', 'running', 'succeeded', 'failed', 'partial', 'cancelled')),
    source          VARCHAR(20) NOT NULL DEFAULT 'user'
                        CHECK (source IN ('user', 'workflow', 'system', 'import', 'manual', 'scheduled')),
    depth           INT NOT NULL DEFAULT 0,
    error_summary   TEXT,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    duration_ms     INT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE workflow_execution_steps (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    execution_id    UUID NOT NULL REFERENCES workflow_executions(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    action_id       UUID,
    sort_order      INT NOT NULL DEFAULT 0,
    action_type     VARCHAR(40) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending', 'running', 'succeeded', 'failed', 'skipped')),
    input           JSONB NOT NULL DEFAULT '{}'::jsonb,
    output          JSONB NOT NULL DEFAULT '{}'::jsonb,
    error           TEXT,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE workflow_templates (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    key             VARCHAR(80) NOT NULL,
    name            VARCHAR(150) NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    module_api_name VARCHAR(80),
    definition      JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_builtin      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_workflow_templates_builtin_key
    ON workflow_templates (key) WHERE is_builtin = TRUE AND organization_id IS NULL;

CREATE INDEX idx_workflows_org_status ON workflows (organization_id, status);
CREATE INDEX idx_workflows_module ON workflows (module_id);
CREATE INDEX idx_workflow_versions_workflow ON workflow_versions (workflow_id, version DESC);
CREATE INDEX idx_workflow_triggers_version ON workflow_triggers (version_id);
CREATE INDEX idx_workflow_triggers_type ON workflow_triggers (type);
CREATE INDEX idx_workflow_conditions_version ON workflow_conditions (version_id);
CREATE INDEX idx_workflow_conditions_parent ON workflow_conditions (parent_id);
CREATE INDEX idx_workflow_actions_version ON workflow_actions (version_id, sort_order);
CREATE INDEX idx_workflow_executions_org_created ON workflow_executions (organization_id, created_at DESC);
CREATE INDEX idx_workflow_executions_workflow ON workflow_executions (workflow_id, created_at DESC);
CREATE INDEX idx_workflow_executions_record ON workflow_executions (record_id) WHERE record_id IS NOT NULL;
CREATE INDEX idx_workflow_execution_steps_exec ON workflow_execution_steps (execution_id, sort_order);

COMMENT ON TABLE automation_rules IS 'Deprecated placeholder; use workflows* tables instead.';
