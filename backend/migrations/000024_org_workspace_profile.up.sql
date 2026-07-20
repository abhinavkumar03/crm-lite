-- Workspace (organization) profile: description + soft delete.

ALTER TABLE organizations
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ NULL;

CREATE INDEX IF NOT EXISTS idx_organizations_deleted_at
    ON organizations(deleted_at)
    WHERE deleted_at IS NULL;
