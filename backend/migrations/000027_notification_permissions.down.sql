DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT id FROM permissions
    WHERE key IN (
        'notification.view',
        'notification.send',
        'notification.templates.manage'
    )
);

DELETE FROM permissions
WHERE key IN (
    'notification.view',
    'notification.send',
    'notification.templates.manage'
);
