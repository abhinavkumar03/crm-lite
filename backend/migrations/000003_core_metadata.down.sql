-- Reverse of 000003. Drop in reverse dependency order. Indexes are dropped
-- automatically with their tables.

DROP TABLE IF EXISTS records;

DROP TABLE IF EXISTS tour_progress;
DROP TABLE IF EXISTS tour_steps;

DROP TABLE IF EXISTS export_templates;
DROP TABLE IF EXISTS import_templates;

DROP TABLE IF EXISTS automation_rules;
DROP TABLE IF EXISTS validation_rules;

DROP TABLE IF EXISTS views;
DROP TABLE IF EXISTS layouts;

DROP TABLE IF EXISTS fields;
DROP TABLE IF EXISTS modules;

DROP TABLE IF EXISTS organization_members;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;

DROP TABLE IF EXISTS organizations;
