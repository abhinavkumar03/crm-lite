-- Restore the Phase 3 placeholder shape so a rollback to 000006 leaves the
-- schema exactly as it was before this migration.
DROP TABLE IF EXISTS tour_progress;

CREATE TABLE tour_progress (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tour_key     VARCHAR(80) NOT NULL,
    status       VARCHAR(20) NOT NULL DEFAULT 'in_progress'
                     CHECK (status IN ('in_progress', 'completed', 'skipped')),
    completed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, tour_key)
);

CREATE INDEX idx_tour_progress_user ON tour_progress(user_id);
