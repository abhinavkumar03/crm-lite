-- Phase 14: Guided CRM tour.
--
-- The Phase 3 schema shipped a placeholder `tour_progress` (user-only, no step
-- cursor, no org scoping). The guided-tour engine needs organization scoping and
-- a persisted step cursor + seen-step list, so we replace that placeholder with
-- the richer shape the engine actually uses. `tour_progress` holds only transient
-- onboarding state, so dropping and recreating it is safe.

DROP TABLE IF EXISTS tour_progress;

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
