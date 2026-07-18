DELETE FROM layouts WHERE layout_type = 'detail' AND name = 'Default Detail';

DROP INDEX IF EXISTS idx_activities_entity;
DROP INDEX IF EXISTS idx_attachments_entity;
DROP INDEX IF EXISTS idx_notes_entity;
DROP INDEX IF EXISTS idx_activities_org_module_entity;
DROP INDEX IF EXISTS idx_attachments_org_module_entity;
DROP INDEX IF EXISTS idx_notes_org_module_entity;

ALTER TABLE activities
    DROP COLUMN IF EXISTS module_id,
    DROP COLUMN IF EXISTS organization_id;

ALTER TABLE attachments
    DROP COLUMN IF EXISTS module_id,
    DROP COLUMN IF EXISTS organization_id;

ALTER TABLE notes
    DROP COLUMN IF EXISTS title,
    DROP COLUMN IF EXISTS module_id,
    DROP COLUMN IF EXISTS organization_id;
