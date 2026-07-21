DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions WHERE key IN (
        'workflow.view',
        'workflow.create',
        'workflow.edit',
        'workflow.delete',
        'workflow.publish',
        'workflow.execute',
        'workflow.logs.view'
    )
);

DELETE FROM permissions WHERE key IN (
    'workflow.view',
    'workflow.create',
    'workflow.edit',
    'workflow.delete',
    'workflow.publish',
    'workflow.execute',
    'workflow.logs.view'
);
