-- Guided interaction metadata + stricter step statuses for the Interactive Tutorial Engine.

ALTER TABLE demo_workflow_steps
    ADD COLUMN IF NOT EXISTS required_action VARCHAR(40) NOT NULL DEFAULT 'acknowledge',
    ADD COLUMN IF NOT EXISTS success_event VARCHAR(160),
    ADD COLUMN IF NOT EXISTS failure_message TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS hint TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS max_attempts INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS allow_selectors JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS placement VARCHAR(20) NOT NULL DEFAULT 'center';

ALTER TABLE demo_step_progress DROP CONSTRAINT IF EXISTS demo_step_progress_status_check;
ALTER TABLE demo_step_progress
    ADD CONSTRAINT demo_step_progress_status_check
    CHECK (status IN ('locked', 'active', 'waiting', 'validating', 'completed', 'skipped', 'failed'));

-- Tighten target selectors for mentor UX (click-through spotlight).
UPDATE demo_workflow_steps SET
    target_selector = '[data-tour="dashboard"]',
    placement = 'bottom',
    required_action = 'navigate',
    hint = 'Open the dashboard from the sidebar, then confirm.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'dashboard';

UPDATE demo_workflow_steps SET
    target_selector = '[data-tour="nav-settings"]',
    placement = 'right',
    required_action = 'navigate',
    hint = 'Click Settings in the sidebar, then open Modules.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'modules_settings';

UPDATE demo_workflow_steps SET
    target_selector = '[data-tutorial-action="create-module"]',
    placement = 'left',
    required_action = 'create_resource',
    success_event = 'api:module_created',
    failure_message = 'Create a module with api_name tutorial_lead first.',
    hint = 'Click New module, set api_name to tutorial_lead, then Create module.',
    allow_selectors = '["[data-tutorial-action=\"create-module\"]","[data-tutorial-action=\"submit-module\"]"]'::jsonb
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'create_module';

UPDATE demo_workflow_steps SET
    target_selector = '[data-tutorial-action="add-field"]',
    placement = 'left',
    required_action = 'create_resource',
    success_event = 'api:field_created',
    failure_message = 'Add a text field company_name on tutorial_lead.',
    hint = 'Select Tutorial Lead, click New field, api_name company_name.',
    allow_selectors = '["[data-tutorial-action=\"add-field\"]","[data-tutorial-action=\"submit-field\"]"]'::jsonb
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'create_field';

UPDATE demo_workflow_steps SET
    target_selector = '[data-tutorial-action="create-record"]',
    placement = 'top',
    required_action = 'create_resource',
    success_event = 'api:record_created',
    failure_message = 'Create at least one Tutorial Lead record.',
    hint = 'Pick Tutorial Lead on Forms and submit Create record.',
    allow_selectors = '["[data-tutorial-action=\"create-record\"]"]'::jsonb
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'create_record';

UPDATE demo_workflow_steps SET
    target_selector = '[data-tour="nav-tables"]',
    placement = 'right',
    required_action = 'navigate',
    hint = 'Open Tables, select Tutorial Lead, click a row.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'record_workspace';

UPDATE demo_workflow_steps SET
    target_selector = '[data-tutorial-action="add-note"]',
    placement = 'left',
    required_action = 'create_resource',
    failure_message = 'Add a note on a Tutorial Lead record.',
    hint = 'Open Notes tab and click Add.',
    allow_selectors = '["[data-tutorial-action=\"add-note\"]"]'::jsonb
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'add_note';

UPDATE demo_workflow_steps SET
    required_action = 'acknowledge',
    placement = 'center'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key IN ('welcome', 'completion');
