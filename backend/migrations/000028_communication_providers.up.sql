-- Communication providers, sender identities, notification correlation fields,
-- template versioning, preferences, and open-tracking tokens.

CREATE TABLE IF NOT EXISTS communication_providers (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id   UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    channel           VARCHAR(20) NOT NULL CHECK (channel IN ('email', 'whatsapp')),
    provider_type     VARCHAR(40) NOT NULL,
    name              VARCHAR(120) NOT NULL,
    config            JSONB NOT NULL DEFAULT '{}'::jsonb,
    secrets_encrypted BYTEA,
    is_default        BOOLEAN NOT NULL DEFAULT FALSE,
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    last_health_at    TIMESTAMPTZ,
    last_error        TEXT,
    created_by        UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, channel, name)
);

CREATE INDEX idx_communication_providers_org_channel
    ON communication_providers(organization_id, channel, is_active);

CREATE UNIQUE INDEX idx_communication_providers_default
    ON communication_providers(organization_id, channel)
    WHERE is_default = TRUE AND is_active = TRUE;

CREATE TABLE IF NOT EXISTS communication_sender_identities (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id   UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    provider_id       UUID REFERENCES communication_providers(id) ON DELETE SET NULL,
    channel           VARCHAR(20) NOT NULL CHECK (channel IN ('email', 'whatsapp')),
    display_name      VARCHAR(160),
    from_address      VARCHAR(255) NOT NULL,
    reply_to          VARCHAR(255),
    is_default        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_communication_sender_identities_org
    ON communication_sender_identities(organization_id, channel);

ALTER TABLE notifications
    ADD COLUMN IF NOT EXISTS provider_id UUID REFERENCES communication_providers(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS provider_message_id VARCHAR(255),
    ADD COLUMN IF NOT EXISTS from_address VARCHAR(255),
    ADD COLUMN IF NOT EXISTS reply_to VARCHAR(255),
    ADD COLUMN IF NOT EXISTS next_retry_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS open_tracking_token VARCHAR(64);

CREATE INDEX IF NOT EXISTS idx_notifications_provider_message_id
    ON notifications(provider_message_id)
    WHERE provider_message_id IS NOT NULL AND provider_message_id <> '';

CREATE INDEX IF NOT EXISTS idx_notifications_open_token
    ON notifications(open_tracking_token)
    WHERE open_tracking_token IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_notifications_next_retry
    ON notifications(next_retry_at)
    WHERE status = 'retrying' AND next_retry_at IS NOT NULL;

-- Template draft/publish + expanded categories.
ALTER TABLE notification_templates
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'published'
        CHECK (status IN ('draft', 'published')),
    ADD COLUMN IF NOT EXISTS version INT NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS whatsapp_template_name VARCHAR(120),
    ADD COLUMN IF NOT EXISTS whatsapp_language VARCHAR(20);

ALTER TABLE notification_templates DROP CONSTRAINT IF EXISTS notification_templates_category_check;
ALTER TABLE notification_templates
    ADD CONSTRAINT notification_templates_category_check CHECK (category IN (
        'sales', 'follow_up', 'welcome', 'proposal', 'invoice', 'reminder',
        'quotation', 'marketing', 'support', 'custom'
    ));

CREATE TABLE IF NOT EXISTS notification_preferences (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id   UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id           UUID REFERENCES users(id) ON DELETE CASCADE,
    email_enabled     BOOLEAN NOT NULL DEFAULT TRUE,
    whatsapp_enabled  BOOLEAN NOT NULL DEFAULT TRUE,
    transactional     BOOLEAN NOT NULL DEFAULT TRUE,
    marketing         BOOLEAN NOT NULL DEFAULT TRUE,
    do_not_disturb    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, user_id)
);

INSERT INTO permissions (key, category, description) VALUES
    ('communication.providers.manage', 'notification', 'Manage email/WhatsApp providers and sender identities')
ON CONFLICT (key) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.is_system = TRUE
  AND r.slug IN ('owner', 'super_admin', 'admin')
  AND p.key = 'communication.providers.manage'
ON CONFLICT DO NOTHING;
