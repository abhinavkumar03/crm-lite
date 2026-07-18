-- Phase 3: Core metadata database.
--
-- Introduces multi-tenancy (organizations + membership), RBAC (roles /
-- permissions), the metadata catalog (modules, fields, layouts, views,
-- validation & automation rules, import/export templates, tour), and the
-- generic JSONB-backed record engine that lets new modules exist without any
-- schema change.
--
-- New tables use TIMESTAMPTZ (timezone-aware) and gen_random_uuid() defaults
-- (pgcrypto was enabled in 000001).

-- =====================================================================
-- Multi-tenancy
-- =====================================================================

CREATE TABLE organizations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(200) NOT NULL,
    slug        VARCHAR(120) NOT NULL UNIQUE,
    plan        VARCHAR(50)  NOT NULL DEFAULT 'free',
    settings    JSONB        NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- =====================================================================
-- RBAC
-- =====================================================================

CREATE TABLE roles (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name             VARCHAR(100) NOT NULL,
    slug             VARCHAR(100) NOT NULL,
    description      TEXT,
    is_system        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, slug)
);

-- Global catalog of permission keys (e.g. "module.create", "field.update",
-- "import.run"). Seeded in a later phase.
CREATE TABLE permissions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key         VARCHAR(120) NOT NULL UNIQUE,
    category    VARCHAR(60)  NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE role_permissions (
    role_id        UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id  UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE organization_members (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id          UUID REFERENCES roles(id) ON DELETE SET NULL,
    status           VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, user_id)
);

-- =====================================================================
-- Metadata catalog
-- =====================================================================

-- A module is a metadata-defined object type. storage_strategy = 'native' for
-- the existing first-class tables (leads/contacts/tasks) surfaced through the
-- metadata layer; 'dynamic' modules store their rows in the records table.
CREATE TABLE modules (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id      UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    api_name             VARCHAR(80)  NOT NULL,
    singular_label       VARCHAR(80)  NOT NULL,
    plural_label         VARCHAR(80)  NOT NULL,
    description          TEXT,
    icon                 VARCHAR(60),
    color                VARCHAR(20),
    storage_strategy     VARCHAR(20)  NOT NULL DEFAULT 'dynamic'
                             CHECK (storage_strategy IN ('native', 'dynamic')),
    native_table         VARCHAR(80),
    is_system            BOOLEAN NOT NULL DEFAULT FALSE,
    is_enabled           BOOLEAN NOT NULL DEFAULT TRUE,
    is_visible_sidebar   BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order           INTEGER NOT NULL DEFAULT 0,
    default_sort_field   VARCHAR(80) NOT NULL DEFAULT 'created_at',
    default_sort_order   VARCHAR(4)  NOT NULL DEFAULT 'desc'
                             CHECK (default_sort_order IN ('asc', 'desc')),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, api_name)
);

CREATE TABLE fields (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id    UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    module_id          UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    api_name           VARCHAR(80) NOT NULL,
    label              VARCHAR(120) NOT NULL,
    field_type         VARCHAR(30) NOT NULL CHECK (field_type IN (
                           'text','textarea','email','phone','number','currency',
                           'date','datetime','boolean','dropdown','multiselect',
                           'radio','checkbox','url','file','image','user','lookup',
                           'formula','json','richtext'
                       )),
    is_required        BOOLEAN NOT NULL DEFAULT FALSE,
    is_unique          BOOLEAN NOT NULL DEFAULT FALSE,
    is_read_only       BOOLEAN NOT NULL DEFAULT FALSE,
    default_value      TEXT,
    placeholder        VARCHAR(200),
    description        TEXT,
    help_text          TEXT,
    min_length         INTEGER,
    max_length         INTEGER,
    regex              TEXT,
    validation_message TEXT,
    options            JSONB NOT NULL DEFAULT '[]'::jsonb,
    lookup_module_id   UUID REFERENCES modules(id) ON DELETE SET NULL,
    sort_order         INTEGER NOT NULL DEFAULT 0,
    is_visible         BOOLEAN NOT NULL DEFAULT TRUE,
    is_searchable      BOOLEAN NOT NULL DEFAULT FALSE,
    is_filterable      BOOLEAN NOT NULL DEFAULT FALSE,
    is_nullable        BOOLEAN NOT NULL DEFAULT TRUE,
    is_indexed         BOOLEAN NOT NULL DEFAULT FALSE,
    is_system          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (module_id, api_name)
);

-- Form / detail layouts (arrangement of fields into sections).
CREATE TABLE layouts (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    module_id        UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    name             VARCHAR(120) NOT NULL,
    layout_type      VARCHAR(20)  NOT NULL DEFAULT 'form'
                         CHECK (layout_type IN ('form', 'detail')),
    is_default       BOOLEAN NOT NULL DEFAULT FALSE,
    config           JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Saved table views (visible columns, filters, sorting).
CREATE TABLE views (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    module_id        UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    name             VARCHAR(120) NOT NULL,
    columns          JSONB NOT NULL DEFAULT '[]'::jsonb,
    filters          JSONB NOT NULL DEFAULT '[]'::jsonb,
    sort             JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_default       BOOLEAN NOT NULL DEFAULT FALSE,
    is_public        BOOLEAN NOT NULL DEFAULT TRUE,
    owner_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Database-driven validation rules (consumed by the validation engine in a
-- later phase). field_id NULL = a module-level (cross-field) rule.
CREATE TABLE validation_rules (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    module_id        UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    field_id         UUID REFERENCES fields(id) ON DELETE CASCADE,
    rule_type        VARCHAR(40) NOT NULL,
    params           JSONB NOT NULL DEFAULT '{}'::jsonb,
    error_message    TEXT,
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order       INTEGER NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Event-driven automation (trigger -> conditions -> actions).
CREATE TABLE automation_rules (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    module_id        UUID REFERENCES modules(id) ON DELETE CASCADE,
    name             VARCHAR(150) NOT NULL,
    trigger_event    VARCHAR(80)  NOT NULL,
    conditions       JSONB NOT NULL DEFAULT '[]'::jsonb,
    actions          JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE import_templates (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    module_id        UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    name             VARCHAR(150) NOT NULL,
    mapping          JSONB NOT NULL DEFAULT '{}'::jsonb,
    options          JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE export_templates (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    module_id        UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    name             VARCHAR(150) NOT NULL,
    columns          JSONB NOT NULL DEFAULT '[]'::jsonb,
    filters          JSONB NOT NULL DEFAULT '[]'::jsonb,
    format           VARCHAR(10) NOT NULL DEFAULT 'csv'
                         CHECK (format IN ('csv', 'xlsx')),
    created_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Guided product tour. Steps can be global (organization_id NULL) or per-org.
CREATE TABLE tour_steps (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID REFERENCES organizations(id) ON DELETE CASCADE,
    step_key         VARCHAR(80)  NOT NULL,
    title            VARCHAR(150) NOT NULL,
    body             TEXT NOT NULL,
    target_selector  VARCHAR(200),
    route            VARCHAR(200),
    placement        VARCHAR(20) NOT NULL DEFAULT 'bottom',
    sort_order       INTEGER NOT NULL DEFAULT 0,
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tour_progress (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tour_key     VARCHAR(80) NOT NULL,
    status       VARCHAR(20) NOT NULL DEFAULT 'in_progress'
                     CHECK (status IN ('in_progress', 'completed', 'skipped')),
    completed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, tour_key)
);

-- =====================================================================
-- Generic record engine (JSONB storage strategy)
-- =====================================================================

CREATE TABLE records (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    module_id        UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    data             JSONB NOT NULL DEFAULT '{}'::jsonb,
    owner_id         UUID REFERENCES users(id) ON DELETE SET NULL,
    created_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- =====================================================================
-- Indexes
-- =====================================================================

CREATE INDEX idx_roles_org               ON roles(organization_id);
CREATE INDEX idx_org_members_org         ON organization_members(organization_id);
CREATE INDEX idx_org_members_user        ON organization_members(user_id);
CREATE INDEX idx_role_permissions_role   ON role_permissions(role_id);

CREATE INDEX idx_modules_org             ON modules(organization_id);
CREATE INDEX idx_fields_module           ON fields(module_id);
CREATE INDEX idx_fields_org              ON fields(organization_id);
CREATE INDEX idx_layouts_module          ON layouts(module_id);
CREATE INDEX idx_views_module            ON views(module_id);
CREATE INDEX idx_validation_rules_module ON validation_rules(module_id);
CREATE INDEX idx_validation_rules_field  ON validation_rules(field_id);
CREATE INDEX idx_automation_rules_module ON automation_rules(module_id);
CREATE INDEX idx_automation_rules_event  ON automation_rules(trigger_event);
CREATE INDEX idx_import_templates_module ON import_templates(module_id);
CREATE INDEX idx_export_templates_module ON export_templates(module_id);
CREATE INDEX idx_tour_steps_org          ON tour_steps(organization_id);
CREATE INDEX idx_tour_progress_user      ON tour_progress(user_id);

CREATE INDEX idx_records_org_module      ON records(organization_id, module_id);
CREATE INDEX idx_records_owner           ON records(owner_id);
CREATE INDEX idx_records_created_at      ON records(created_at DESC);
-- GIN index enables fast containment/existence queries against dynamic fields.
CREATE INDEX idx_records_data_gin        ON records USING GIN (data jsonb_path_ops);
