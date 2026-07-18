-- Phase 14: Guided CRM tour.
--
-- tour_progress persists a user's onboarding progress so the guided tour can be
-- resumed across devices/sessions and restarted on demand. Progress is scoped to
-- (organization, user, tour_key); tour_key lets a single user carry independent
-- progress for multiple named tours (the app ships one, "app", but the schema is
-- open for feature-specific tours later). The step catalogue itself lives on the
-- client — only lightweight progress is stored server-side.

CREATE TABLE tour_progress (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tour_key         VARCHAR(80)  NOT NULL DEFAULT 'app',
    status           VARCHAR(20)  NOT NULL DEFAULT 'active'
                         CHECK (status IN ('active', 'completed', 'skipped')),
    current_step     INTEGER NOT NULL DEFAULT 0,
    completed_steps  JSONB   NOT NULL DEFAULT '[]'::jsonb,
    started_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at     TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (organization_id, user_id, tour_key)
);

CREATE INDEX idx_tour_progress_user ON tour_progress(user_id);
