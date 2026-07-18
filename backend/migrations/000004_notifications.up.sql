-- Phase 11: WhatsApp automation / notification pipeline.
--
-- The notifications table is the durable outbound-message log. Every message
-- (email or WhatsApp) is persisted with a lifecycle status (queued -> sent /
-- failed) so delivery is auditable and retryable independently of the request
-- that created it. Rendering happens at enqueue time; the worker only dispatches
-- and updates status.

CREATE TABLE notifications (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    channel          VARCHAR(20)  NOT NULL CHECK (channel IN ('email', 'whatsapp')),
    recipient        VARCHAR(255) NOT NULL,
    subject          VARCHAR(255),
    body             TEXT NOT NULL,
    template         VARCHAR(120),
    data             JSONB NOT NULL DEFAULT '{}'::jsonb,
    status           VARCHAR(20)  NOT NULL DEFAULT 'queued'
                         CHECK (status IN ('queued', 'sent', 'failed')),
    provider         VARCHAR(60),
    error            TEXT,
    entity_type      VARCHAR(40),
    entity_id        UUID,
    created_by       UUID REFERENCES users(id) ON DELETE SET NULL,
    sent_at          TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_org_status  ON notifications(organization_id, status);
CREATE INDEX idx_notifications_entity      ON notifications(entity_type, entity_id);
CREATE INDEX idx_notifications_created_at  ON notifications(created_at DESC);
