DROP TABLE IF EXISTS notification_attachments;
DROP TABLE IF EXISTS notification_delivery_events;

DROP INDEX IF EXISTS idx_notifications_scheduled_due;
DROP INDEX IF EXISTS idx_notifications_org_module_entity;
DROP INDEX IF EXISTS idx_notifications_org_channel_created;
DROP INDEX IF EXISTS idx_notifications_org_status_created;

ALTER TABLE notifications
    DROP COLUMN IF EXISTS module_id,
    DROP COLUMN IF EXISTS cc,
    DROP COLUMN IF EXISTS bcc,
    DROP COLUMN IF EXISTS body_html,
    DROP COLUMN IF EXISTS template_id,
    DROP COLUMN IF EXISTS scheduled_at,
    DROP COLUMN IF EXISTS cancelled_at,
    DROP COLUMN IF EXISTS retry_count,
    DROP COLUMN IF EXISTS max_retries,
    DROP COLUMN IF EXISTS last_error,
    DROP COLUMN IF EXISTS provider_response,
    DROP COLUMN IF EXISTS variables_used,
    DROP COLUMN IF EXISTS attachment_ids,
    DROP COLUMN IF EXISTS queued_at,
    DROP COLUMN IF EXISTS processing_at,
    DROP COLUMN IF EXISTS delivered_at,
    DROP COLUMN IF EXISTS opened_at,
    DROP COLUMN IF EXISTS read_at;

-- Restore original status check (rows with new statuses must be cleaned first).
UPDATE notifications
SET status = CASE
    WHEN status IN ('draft', 'scheduled', 'cancelled') THEN 'queued'
    WHEN status IN ('processing', 'retrying') THEN 'queued'
    WHEN status IN ('delivered', 'opened', 'read') THEN 'sent'
    ELSE status
END
WHERE status NOT IN ('queued', 'sent', 'failed');

ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_status_check;
ALTER TABLE notifications
    ADD CONSTRAINT notifications_status_check CHECK (status IN ('queued', 'sent', 'failed'));

DROP TABLE IF EXISTS notification_templates;
