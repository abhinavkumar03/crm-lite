-- Enterprise tenancy: org profile, active org, invitations, structure,
-- hierarchy, and record visibility (roadmap Phases 1–4 schema).

-- =====================================================================
-- Phase 1 — Organization profile
-- =====================================================================

ALTER TABLE organizations
    ADD COLUMN IF NOT EXISTS logo_url      TEXT,
    ADD COLUMN IF NOT EXISTS industry      VARCHAR(120),
    ADD COLUMN IF NOT EXISTS company_size  VARCHAR(40),
    ADD COLUMN IF NOT EXISTS country       VARCHAR(80),
    ADD COLUMN IF NOT EXISTS status        VARCHAR(20) NOT NULL DEFAULT 'active',
    ADD COLUMN IF NOT EXISTS created_by    UUID REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE organizations
    DROP CONSTRAINT IF EXISTS organizations_status_check;
ALTER TABLE organizations
    ADD CONSTRAINT organizations_status_check
        CHECK (status IN ('active', 'suspended', 'trial', 'inactive'));

-- =====================================================================
-- Phase 2 — Active organization + invitations
-- =====================================================================

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS active_organization_id UUID
        REFERENCES organizations(id) ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS organization_invitations (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email            VARCHAR(255) NOT NULL,
    role_id          UUID REFERENCES roles(id) ON DELETE SET NULL,
    manager_user_id  UUID REFERENCES users(id) ON DELETE SET NULL,
    department_id    UUID, -- FK added after departments exist
    team_id          UUID,
    token            VARCHAR(64) NOT NULL UNIQUE,
    status           VARCHAR(20) NOT NULL DEFAULT 'pending'
                         CHECK (status IN ('pending', 'accepted', 'revoked', 'expired')),
    invited_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    expires_at       TIMESTAMPTZ NOT NULL,
    accepted_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_org_invitations_org_email
    ON organization_invitations(organization_id, email);
CREATE INDEX IF NOT EXISTS idx_org_invitations_token
    ON organization_invitations(token);

-- =====================================================================
-- Phase 3 — Departments, teams, branches + hierarchy
-- =====================================================================

CREATE TABLE IF NOT EXISTS departments (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name             VARCHAR(120) NOT NULL,
    description      TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, name)
);

CREATE TABLE IF NOT EXISTS branches (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name             VARCHAR(120) NOT NULL,
    location         VARCHAR(200),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, name)
);

CREATE TABLE IF NOT EXISTS teams (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    department_id    UUID REFERENCES departments(id) ON DELETE SET NULL,
    name             VARCHAR(120) NOT NULL,
    description      TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, name)
);

ALTER TABLE organization_invitations
    DROP CONSTRAINT IF EXISTS organization_invitations_department_id_fkey;
ALTER TABLE organization_invitations
    ADD CONSTRAINT organization_invitations_department_id_fkey
        FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL;
ALTER TABLE organization_invitations
    DROP CONSTRAINT IF EXISTS organization_invitations_team_id_fkey;
ALTER TABLE organization_invitations
    ADD CONSTRAINT organization_invitations_team_id_fkey
        FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE SET NULL;

ALTER TABLE roles
    ADD COLUMN IF NOT EXISTS hierarchy_level INT NOT NULL DEFAULT 100,
    ADD COLUMN IF NOT EXISTS parent_role_id UUID REFERENCES roles(id) ON DELETE SET NULL;

ALTER TABLE organization_members
    ADD COLUMN IF NOT EXISTS manager_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS department_id   UUID REFERENCES departments(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS team_id         UUID REFERENCES teams(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS branch_id       UUID REFERENCES branches(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS designation     VARCHAR(120),
    ADD COLUMN IF NOT EXISTS hierarchy_level INT NOT NULL DEFAULT 100;

CREATE INDEX IF NOT EXISTS idx_org_members_manager
    ON organization_members(organization_id, manager_user_id);
CREATE INDEX IF NOT EXISTS idx_org_members_department
    ON organization_members(organization_id, department_id);

-- =====================================================================
-- Phase 4 — Record visibility / ownership
-- =====================================================================

ALTER TABLE records
    ADD COLUMN IF NOT EXISTS assigned_to   UUID REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS team_id       UUID REFERENCES teams(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS department_id UUID REFERENCES departments(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS visibility    VARCHAR(20) NOT NULL DEFAULT 'organization';

ALTER TABLE records
    DROP CONSTRAINT IF EXISTS records_visibility_check;
ALTER TABLE records
    ADD CONSTRAINT records_visibility_check
        CHECK (visibility IN (
            'private', 'owner', 'manager', 'hierarchy',
            'department', 'organization', 'team', 'public'
        ));

CREATE INDEX IF NOT EXISTS idx_records_owner
    ON records(organization_id, owner_id);
CREATE INDEX IF NOT EXISTS idx_records_visibility
    ON records(organization_id, visibility);
