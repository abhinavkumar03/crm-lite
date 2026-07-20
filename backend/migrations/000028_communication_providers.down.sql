DROP TABLE IF EXISTS notification_preferences;

ALTER TABLE notification_templates DROP CONSTRAINT IF EXISTS notification_templates_category_check;
ALTER TABLE notification_templates
    ADD CONSTRAINT notification_templates_category_check CHECK (category IN (
        'sales', 'follow_up', 'welcome', 'proposal', 'invoice', 'reminder', 'custom'
    ));

ALTER TABLE notification_templates
    DROP COLUMN IF EXISTS whatsapp_language,
    DROP COLUMN IF EXISTS whatsapp_template_name,
    DROP COLUMN IF EXISTS version,
    DROP COLUMN IF EXISTS status;

DROP INDEX IF EXISTS idx_notifications_next_retry;
DROP INDEX IF EXISTS idx_notifications_open_token;
DROP INDEX IF EXISTS idx_notifications_provider_message_id;

ALTER TABLE notifications
    DROP COLUMN IF EXISTS open_tracking_token,
    DROP COLUMN IF EXISTS next_retry_at,
    DROP COLUMN IF EXISTS reply_to,
    DROP COLUMN IF EXISTS from_address,
    DROP COLUMN IF EXISTS provider_message_id,
    DROP COLUMN IF EXISTS provider_id;

DROP TABLE IF EXISTS communication_sender_identities;
DROP TABLE IF EXISTS communication_providers;

DELETE FROM role_permissions
WHERE permission_id IN (SELECT id FROM permissions WHERE key = 'communication.providers.manage');
DELETE FROM permissions WHERE key = 'communication.providers.manage';
