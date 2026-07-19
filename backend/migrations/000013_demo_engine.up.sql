-- Interactive Demo Engine: metadata workflows + sandbox sessions.

CREATE TABLE IF NOT EXISTS demo_workflows (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_key VARCHAR(80) NOT NULL UNIQUE,
    name         VARCHAR(160) NOT NULL,
    description  TEXT,
    version      INT NOT NULL DEFAULT 1,
    duration_min INT NOT NULL DEFAULT 15,
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS demo_workflow_steps (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_key     VARCHAR(80) NOT NULL REFERENCES demo_workflows(workflow_key) ON DELETE CASCADE,
    step_key         VARCHAR(80) NOT NULL,
    sort_order       INT NOT NULL,
    title            VARCHAR(200) NOT NULL,
    description      TEXT NOT NULL DEFAULT '',
    why_it_matters   TEXT NOT NULL DEFAULT '',
    how_it_works     TEXT NOT NULL DEFAULT '',
    expected_result  TEXT NOT NULL DEFAULT '',
    route            VARCHAR(255),
    target_selector  VARCHAR(120),
    action_label     VARCHAR(120),
    validator_key    VARCHAR(80) NOT NULL DEFAULT 'none',
    validator_params JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_skippable     BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (workflow_key, step_key)
);

CREATE INDEX IF NOT EXISTS idx_demo_workflow_steps_order
    ON demo_workflow_steps(workflow_key, sort_order);

CREATE TABLE IF NOT EXISTS demo_sessions (
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workflow_key              VARCHAR(80) NOT NULL REFERENCES demo_workflows(workflow_key),
    workflow_version          INT NOT NULL DEFAULT 1,
    sandbox_organization_id   UUID REFERENCES organizations(id) ON DELETE SET NULL,
    previous_organization_id  UUID REFERENCES organizations(id) ON DELETE SET NULL,
    status                    VARCHAR(30) NOT NULL DEFAULT 'active'
                                  CHECK (status IN ('active', 'completed', 'abandoned', 'cleaned')),
    current_step_key          VARCHAR(80),
    started_at                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at              TIMESTAMPTZ,
    updated_at                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    keep_data                 BOOLEAN,
    stats                     JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_demo_sessions_user_status
    ON demo_sessions(user_id, status);

CREATE TABLE IF NOT EXISTS demo_step_progress (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id   UUID NOT NULL REFERENCES demo_sessions(id) ON DELETE CASCADE,
    step_key     VARCHAR(80) NOT NULL,
    status       VARCHAR(20) NOT NULL DEFAULT 'locked'
                     CHECK (status IN ('locked', 'active', 'completed', 'skipped', 'failed')),
    attempts     INT NOT NULL DEFAULT 0,
    last_error   TEXT,
    completed_at TIMESTAMPTZ,
    UNIQUE (session_id, step_key)
);

CREATE TABLE IF NOT EXISTS demo_resources (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id    UUID NOT NULL REFERENCES demo_sessions(id) ON DELETE CASCADE,
    resource_type VARCHAR(60) NOT NULL,
    resource_id   UUID NOT NULL,
    meta          JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_demo_resources_session
    ON demo_resources(session_id);

CREATE TABLE IF NOT EXISTS demo_events (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES demo_sessions(id) ON DELETE CASCADE,
    event_type VARCHAR(80) NOT NULL,
    payload    JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_demo_events_session
    ON demo_events(session_id, created_at);

-- Canonical CRM Interactive Walkthrough workflow
INSERT INTO demo_workflows (workflow_key, name, description, version, duration_min)
VALUES (
    'crm_interactive_walkthrough',
    'Interactive CRM Walkthrough',
    'Hands-on sandbox tutorial covering metadata modules, records, workspace, and settings.',
    1,
    15
)
ON CONFLICT (workflow_key) DO NOTHING;

INSERT INTO demo_workflow_steps (
    workflow_key, step_key, sort_order, title, description, why_it_matters, how_it_works,
    expected_result, route, target_selector, action_label, validator_key, validator_params, is_skippable
) VALUES
(
    'crm_interactive_walkthrough', 'welcome', 1,
    'Welcome to the CRM sandbox',
    'You are inside an isolated demo organization. Nothing here affects your real workspace.',
    'Recruiters and first-time users need a safe place to explore without fear of breaking data.',
    'The demo engine provisioned a sandbox organization, switched your active tenant context, and seeded a Product Demo module.',
    'You see this instruction panel and can navigate freely inside the sandbox.',
    '/dashboard', NULL, 'I understand — continue', 'none', '{}', TRUE
),
(
    'crm_interactive_walkthrough', 'dashboard', 2,
    'Explore the dashboard',
    'Open the dashboard and notice module counts and recent activity for the sandbox org.',
    'Dashboards summarize tenant-scoped metrics from the record engine (with Redis caching in production).',
    'GET /dashboard reads organization-scoped aggregates from dynamic records.',
    'Dashboard cards render for the sandbox organization.',
    '/dashboard', 'data-tour=dashboard', 'Open Dashboard', 'route_visited', '{"route":"/dashboard"}', TRUE
),
(
    'crm_interactive_walkthrough', 'modules_settings', 3,
    'Open Modules settings',
    'Go to Settings → Modules. This is where metadata-driven object types are managed.',
    'CRM Lite treats business objects as data (modules/fields), not hardcoded screens.',
    'Modules live in the modules table; the UI loads them via GET /modules.',
    'The Modules settings page is visible.',
    '/settings/modules', NULL, 'Open Modules', 'route_visited', '{"route":"/settings/modules"}', TRUE
),
(
    'crm_interactive_walkthrough', 'create_module', 4,
    'Create a custom module',
    'Create a module with API name tutorial_lead (singular Tutorial Lead). Enable sidebar visibility.',
    'Custom modules prove the metadata engine can extend the CRM without shipping new frontend pages.',
    'POST /modules inserts a dynamic module; tables/forms/workspace pick it up automatically.',
    'A module with api_name tutorial_lead exists in the sandbox org.',
    '/settings/modules', NULL, 'Create module then Validate', 'module_exists', '{"api_name":"tutorial_lead"}', FALSE
),
(
    'crm_interactive_walkthrough', 'create_field', 5,
    'Add a dynamic field',
    'On Fields settings for tutorial_lead, add a required text field with api_name company_name.',
    'Fields drive forms, tables, validation, import mapping, and workspace overview sections.',
    'POST /modules/:id/fields stores field metadata consumed by DynamicForm and DynamicTable.',
    'Field company_name exists on tutorial_lead.',
    '/settings/fields', NULL, 'Add field then Validate', 'field_exists', '{"module_api_name":"tutorial_lead","api_name":"company_name"}', FALSE
),
(
    'crm_interactive_walkthrough', 'create_record', 6,
    'Create your first record',
    'Open Forms (or Tables → Add), select Tutorial Lead, and create a record with a company name.',
    'Records are JSONB rows scoped by organization_id + module_id — one runtime for every module.',
    'POST /modules/:id/records validates via the validation engine then persists into records.data.',
    'At least one tutorial_lead record exists.',
    '/forms', NULL, 'Create record then Validate', 'record_exists', '{"module_api_name":"tutorial_lead"}', FALSE
),
(
    'crm_interactive_walkthrough', 'record_workspace', 7,
    'Open the Record Workspace',
    'From Tables, click a Tutorial Lead row to open the metadata-driven Record Workspace.',
    'List views are for scanning; the workspace is where users collaborate on a single record lifecycle.',
    'GET /modules/:id/records/:rid loads the record; tabs reuse notes/attachments/activities APIs.',
    'You can open a record detail URL under /tables/:moduleId/:recordId.',
    '/tables', NULL, 'Open a record then Validate', 'workspace_opened', '{"module_api_name":"tutorial_lead"}', TRUE
),
(
    'crm_interactive_walkthrough', 'add_note', 8,
    'Add a note on the record',
    'In the workspace Notes tab, add a note describing a follow-up.',
    'Collaboration artifacts are polymorphic and attached to any dynamic record.',
    'POST /modules/:id/records/:rid/notes writes entity_type=RECORD and emits a timeline activity.',
    'At least one note exists for a tutorial_lead record.',
    '/tables', NULL, 'Add note then Validate', 'note_exists', '{"module_api_name":"tutorial_lead"}', FALSE
),
(
    'crm_interactive_walkthrough', 'timeline', 9,
    'Inspect the timeline',
    'Open the Timeline tab and confirm create/update/note events appear.',
    'A unified activity stream explains who changed what — critical for CRM auditability.',
    'GET …/activities lists org-scoped RECORD activities written by the workspace + record engines.',
    'Timeline shows at least one activity for the sandbox records.',
    '/tables', NULL, 'Check timeline then Validate', 'activity_exists', '{"module_api_name":"tutorial_lead"}', TRUE
),
(
    'crm_interactive_walkthrough', 'product_demo_module', 10,
    'Explore the seeded Product Demo module',
    'Open Tables and select the Product Demo module seeded for this sandbox. Browse its fields and records.',
    'The demo engine pre-seeds a showcase module so you can see lookups, dropdowns, and sample data immediately.',
    'Sandbox bootstrap created module api_name=product_demo with sample fields and records.',
    'product_demo module is visible in navigation/tables.',
    '/tables', NULL, 'Open Product Demo then Validate', 'module_exists', '{"api_name":"product_demo"}', TRUE
),
(
    'crm_interactive_walkthrough', 'roles_glance', 11,
    'Glance at roles & permissions',
    'Visit Settings → Roles to see hierarchy levels and the permission matrix for the sandbox org.',
    'RBAC is grant-based with hierarchy_level for visibility defaults — not silent inheritance.',
    'GET /roles returns org-scoped roles seeded at organization bootstrap (Owner, Admin, …).',
    'Roles settings page loads.',
    '/settings/roles', NULL, 'Open Roles', 'route_visited', '{"route":"/settings/roles"}', TRUE
),
(
    'crm_interactive_walkthrough', 'completion', 12,
    'Walkthrough complete',
    'You have exercised the metadata CRM loop: modules → fields → records → workspace → collaboration → RBAC.',
    'A guided sandbox is the fastest way for reviewers to understand architecture without reading docs.',
    'The demo state machine marks the session completed and offers keep/delete for sandbox data.',
    'Completion screen with stats and cleanup choices.',
    '/dashboard', NULL, 'Finish', 'none', '{}', TRUE
)
ON CONFLICT (workflow_key, step_key) DO NOTHING;
