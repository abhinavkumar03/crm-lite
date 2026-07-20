-- Notification Center RBAC keys + Owner/Admin grants.

INSERT INTO permissions (key, category, description) VALUES
    ('notification.view', 'notification', 'View notification center and delivery history'),
    ('notification.send', 'notification', 'Compose, send, schedule, retry and cancel notifications'),
    ('notification.templates.manage', 'notification', 'Create and manage notification templates')
ON CONFLICT (key) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.is_system = TRUE
  AND r.slug IN ('owner', 'super_admin', 'admin')
  AND p.key IN (
      'notification.view',
      'notification.send',
      'notification.templates.manage'
  )
ON CONFLICT DO NOTHING;
