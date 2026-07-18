-- export_templates belongs to 000003; only roll back the Phase 13 jobs table.
DROP INDEX IF EXISTS idx_export_templates_org_module;

ALTER TABLE export_templates
    DROP COLUMN IF EXISTS sort;

DROP TABLE IF EXISTS export_jobs;
