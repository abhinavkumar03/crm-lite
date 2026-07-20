-- Enterprise Notification Center: richer lifecycle, templates, delivery events.

-- Templates first so notifications.template_id can reference them.
CREATE TABLE notification_templates (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    channel          VARCHAR(20)  NOT NULL CHECK (channel IN ('email', 'whatsapp')),
    name             VARCHAR(120) NOT NULL,
    category         VARCHAR(40)  NOT NULL DEFAULT 'custom'
                         CHECK (category IN (
                             'sales', 'follow_up', 'welcome', 'proposal',
                             'invoice', 'reminder', 'custom'
                         )),
    subject          VARCHAR(255),
    body             TEXT NOT NULL DEFAULT '',
    body_html        TEXT,
    variables        JSONB NOT NULL DEFAULT '[]'::jsonb,
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, channel, name)
);

CREATE INDEX idx_notification_templates_org
    ON notification_templates(organization_id, channel, is_active);

-- Expand notifications status / columns for the communication platform.
ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_status_check;
ALTER TABLE notifications
    ADD CONSTRAINT notifications_status_check CHECK (status IN (
        'draft', 'scheduled', 'queued', 'processing', 'sent', 'delivered',
        'opened', 'read', 'failed', 'retrying', 'cancelled'
    ));

ALTER TABLE notifications
    ADD COLUMN IF NOT EXISTS module_id UUID REFERENCES modules(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS cc TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS bcc TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS body_html TEXT,
    ADD COLUMN IF NOT EXISTS template_id UUID REFERENCES notification_templates(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS scheduled_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS cancelled_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS retry_count INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS max_retries INT NOT NULL DEFAULT 3,
    ADD COLUMN IF NOT EXISTS last_error TEXT,
    ADD COLUMN IF NOT EXISTS provider_response JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS variables_used JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS attachment_ids UUID[] NOT NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS queued_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS processing_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS delivered_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS opened_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS read_at TIMESTAMPTZ;

-- Backfill last_error from legacy error column.
UPDATE notifications SET last_error = error WHERE error IS NOT NULL AND last_error IS NULL;
UPDATE notifications SET queued_at = created_at WHERE status = 'queued' AND queued_at IS NULL;
UPDATE notifications SET sent_at = COALESCE(sent_at, updated_at) WHERE status = 'sent' AND sent_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_notifications_org_status_created
    ON notifications(organization_id, status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_org_channel_created
    ON notifications(organization_id, channel, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_org_module_entity
    ON notifications(organization_id, module_id, entity_id)
    WHERE entity_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_notifications_scheduled_due
    ON notifications(organization_id, scheduled_at)
    WHERE status = 'scheduled';

CREATE TABLE notification_delivery_events (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    notification_id  UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    event            VARCHAR(40) NOT NULL,
    provider         VARCHAR(60),
    payload          JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notification_delivery_events_notification
    ON notification_delivery_events(notification_id, created_at DESC);
CREATE INDEX idx_notification_delivery_events_org
    ON notification_delivery_events(organization_id, created_at DESC);

CREATE TABLE notification_attachments (
    notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    attachment_id   UUID NOT NULL REFERENCES attachments(id) ON DELETE CASCADE,
    PRIMARY KEY (notification_id, attachment_id)
);
