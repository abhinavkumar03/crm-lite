-- Workflow Automation Engine RBAC keys + Owner/Admin grants.

INSERT INTO permissions (key, category, description) VALUES
    ('workflow.view', 'workflow', 'View workflows and builder metadata'),
    ('workflow.create', 'workflow', 'Create workflow drafts'),
    ('workflow.edit', 'workflow', 'Edit workflow drafts and definitions'),
    ('workflow.delete', 'workflow', 'Archive or delete workflows'),
    ('workflow.publish', 'workflow', 'Publish and disable workflows'),
    ('workflow.execute', 'workflow', 'Manually run workflows'),
    ('workflow.logs.view', 'workflow', 'View workflow execution logs')
ON CONFLICT (key) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.is_system = TRUE
  AND r.slug IN ('owner', 'super_admin', 'admin')
  AND p.key IN (
      'workflow.view',
      'workflow.create',
      'workflow.edit',
      'workflow.delete',
      'workflow.publish',
      'workflow.execute',
      'workflow.logs.view'
  )
ON CONFLICT DO NOTHING;
