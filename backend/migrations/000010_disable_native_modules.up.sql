-- Retire seeded native CRM modules (lead/contact/task). Product surface is
-- dynamic-only; legacy SQL tables remain. Child metadata (fields, views,
-- rules, ACL, records, import/export jobs) cascades via ON DELETE CASCADE.
DELETE FROM modules WHERE api_name IN ('lead', 'contact', 'task');
