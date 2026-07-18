-- Drop in reverse dependency order. Polymorphic tables (no FKs) first, then
-- child tables, then parents.
DROP TABLE IF EXISTS activities;

DROP TABLE IF EXISTS attachments;

DROP TABLE IF EXISTS call_logs;

DROP TABLE IF EXISTS notes;

DROP TABLE IF EXISTS activity_logs;

DROP TABLE IF EXISTS tasks;

DROP TABLE IF EXISTS contacts;

DROP TABLE IF EXISTS leads;

DROP TABLE IF EXISTS users;

-- Extension intentionally left in place: it may be shared by other schemas and
-- is harmless to keep.
