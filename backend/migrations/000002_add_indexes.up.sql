CREATE INDEX idx_users_email
ON users(email);

CREATE INDEX idx_leads_owner
ON leads(owner_id);

CREATE INDEX idx_leads_status
ON leads(status);

CREATE INDEX idx_contacts_owner
ON contacts(owner_id);

CREATE INDEX idx_contacts_email
ON contacts(email);

CREATE INDEX idx_contacts_company
ON contacts(company);

CREATE INDEX idx_tasks_owner
ON tasks(owner_id);

CREATE INDEX idx_tasks_status
ON tasks(status);

CREATE INDEX idx_tasks_due_date
ON tasks(due_date);

CREATE INDEX idx_tasks_lead
ON tasks(lead_id);

CREATE INDEX idx_tasks_contact
ON tasks(contact_id);

CREATE INDEX idx_activity_logs_lead
ON activity_logs(lead_id);