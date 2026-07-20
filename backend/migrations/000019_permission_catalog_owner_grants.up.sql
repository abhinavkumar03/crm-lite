-- Ensure the global permission catalog exists and Owner/Admin system roles
-- always receive the full grant set. Fixes orgs created before seed ran
-- (role_permissions was empty → "Missing permission: module.manage").

INSERT INTO permissions (key, category, description) VALUES
    ('module.view', 'module', 'View modules and navigation'),
    ('module.manage', 'module', 'Create, update, reorder and delete modules'),
    ('field.manage', 'field', 'Create, update and delete fields'),
    ('record.view', 'record', 'View records'),
    ('record.create', 'record', 'Create records'),
    ('record.update', 'record', 'Update records'),
    ('record.delete', 'record', 'Delete records'),
    ('import.run', 'import', 'Run data imports'),
    ('export.run', 'export', 'Run data exports'),
    ('automation.manage', 'automation', 'Create and manage automation rules'),
    ('validation.manage', 'validation', 'Create and manage validation rules'),
    ('settings.manage', 'settings', 'Manage organization settings'),
    ('user.manage', 'user', 'Invite and manage users'),
    ('role.manage', 'role', 'Create and manage roles and permissions'),
    ('organization.manage', 'organization', 'Create organizations and manage membership structure'),
    ('analytics.view', 'analytics', 'View dashboards and analytics')
ON CONFLICT (key) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.is_system = TRUE
  AND r.slug IN ('owner', 'super_admin', 'admin')
ON CONFLICT DO NOTHING;
