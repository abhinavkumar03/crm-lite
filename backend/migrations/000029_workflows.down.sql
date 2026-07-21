DROP TABLE IF EXISTS workflow_execution_steps;
DROP TABLE IF EXISTS workflow_executions;
DROP TABLE IF EXISTS workflow_actions;
DROP TABLE IF EXISTS workflow_conditions;
DROP TABLE IF EXISTS workflow_triggers;
ALTER TABLE workflows DROP CONSTRAINT IF EXISTS workflows_published_version_fk;
DROP TABLE IF EXISTS workflow_versions;
DROP TABLE IF EXISTS workflows;
DROP TABLE IF EXISTS workflow_templates;
