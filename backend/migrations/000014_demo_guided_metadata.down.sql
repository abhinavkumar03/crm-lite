ALTER TABLE demo_step_progress DROP CONSTRAINT IF EXISTS demo_step_progress_status_check;
ALTER TABLE demo_step_progress
    ADD CONSTRAINT demo_step_progress_status_check
    CHECK (status IN ('locked', 'active', 'completed', 'skipped', 'failed'));

ALTER TABLE demo_workflow_steps
    DROP COLUMN IF EXISTS required_action,
    DROP COLUMN IF EXISTS success_event,
    DROP COLUMN IF EXISTS failure_message,
    DROP COLUMN IF EXISTS hint,
    DROP COLUMN IF EXISTS max_attempts,
    DROP COLUMN IF EXISTS allow_selectors,
    DROP COLUMN IF EXISTS placement;
