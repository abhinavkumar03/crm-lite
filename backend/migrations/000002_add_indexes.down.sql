DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_leads_owner;
DROP INDEX IF EXISTS idx_leads_status;

DROP INDEX IF EXISTS idx_contacts_owner;
DROP INDEX IF EXISTS idx_contacts_email;
DROP INDEX IF EXISTS idx_contacts_company;

DROP INDEX IF EXISTS idx_tasks_owner;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_due_date;
DROP INDEX IF EXISTS idx_tasks_lead;
DROP INDEX IF EXISTS idx_tasks_contact;

DROP INDEX IF EXISTS idx_activity_logs_lead;

DROP INDEX IF EXISTS idx_notes_entity;
DROP INDEX IF EXISTS idx_notes_created_by;
DROP INDEX IF EXISTS idx_notes_created_at;

DROP INDEX IF EXISTS idx_call_logs_entity;
DROP INDEX IF EXISTS idx_call_logs_created_by;
DROP INDEX IF EXISTS idx_call_logs_followup;

DROP INDEX IF EXISTS idx_attachment_entity;
DROP INDEX IF EXISTS idx_attachment_public_id;
DROP INDEX IF EXISTS idx_attachment_uploaded_by;

DROP INDEX IF EXISTS idx_activity_entity;
DROP INDEX IF EXISTS idx_activity_created_at;
DROP INDEX IF EXISTS idx_activity_user;
