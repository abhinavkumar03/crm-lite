ALTER TABLE records DROP CONSTRAINT IF EXISTS records_visibility_check;
ALTER TABLE records
    DROP COLUMN IF EXISTS assigned_to,
    DROP COLUMN IF EXISTS team_id,
    DROP COLUMN IF EXISTS department_id,
    DROP COLUMN IF EXISTS visibility;

DROP INDEX IF EXISTS idx_records_visibility;
DROP INDEX IF EXISTS idx_records_owner;

ALTER TABLE organization_members
    DROP COLUMN IF EXISTS manager_user_id,
    DROP COLUMN IF EXISTS department_id,
    DROP COLUMN IF EXISTS team_id,
    DROP COLUMN IF EXISTS branch_id,
    DROP COLUMN IF EXISTS designation,
    DROP COLUMN IF EXISTS hierarchy_level;

ALTER TABLE roles
    DROP COLUMN IF EXISTS hierarchy_level,
    DROP COLUMN IF EXISTS parent_role_id;

DROP TABLE IF EXISTS organization_invitations;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS branches;
DROP TABLE IF EXISTS departments;

ALTER TABLE users DROP COLUMN IF EXISTS active_organization_id;

ALTER TABLE organizations DROP CONSTRAINT IF EXISTS organizations_status_check;
ALTER TABLE organizations
    DROP COLUMN IF EXISTS logo_url,
    DROP COLUMN IF EXISTS industry,
    DROP COLUMN IF EXISTS company_size,
    DROP COLUMN IF EXISTS country,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS created_by;
