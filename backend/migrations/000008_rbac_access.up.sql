-- Phase 16: Roles & Permissions — module-level and field-level access control.
--
-- Global permission keys (module.manage, import.run, …) live in `permissions`
-- and are granted via `role_permissions`. These two tables add finer-grained
-- ACL on top of that:
--
--   role_module_access  — which CRUD actions a role may perform on a module
--   role_field_access   — whether a field is hidden / read-only / writable
--
-- Absence of a row means "unrestricted" (inherit full access). An explicit row
-- is always a restriction or an intentional grant for that role+resource pair.

CREATE TABLE role_module_access (
    role_id     UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    module_id   UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    can_view    BOOLEAN NOT NULL DEFAULT TRUE,
    can_create  BOOLEAN NOT NULL DEFAULT TRUE,
    can_update  BOOLEAN NOT NULL DEFAULT TRUE,
    can_delete  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, module_id)
);

CREATE INDEX idx_role_module_access_module ON role_module_access(module_id);

CREATE TABLE role_field_access (
    role_id     UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    field_id    UUID NOT NULL REFERENCES fields(id) ON DELETE CASCADE,
    access      VARCHAR(10) NOT NULL DEFAULT 'write'
                    CHECK (access IN ('hidden', 'read', 'write')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, field_id)
);

CREATE INDEX idx_role_field_access_field ON role_field_access(field_id);
