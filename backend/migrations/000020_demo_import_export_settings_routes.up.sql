-- Move Import / Export demo steps under Settings routes.

UPDATE demo_workflow_steps SET
    route = '/settings/imports',
    target_selector = '[data-tutorial-action="open-imports"]',
    action_label = 'Open Import (Settings)',
    validator_params = '{"route":"/settings/imports"}'::jsonb,
    hint = 'Open Settings → Import from the settings nav.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'import_engine';

UPDATE demo_workflow_steps SET
    route = '/settings/exports',
    target_selector = '[data-tutorial-action="open-exports"]',
    action_label = 'Open Export (Settings)',
    validator_params = '{"route":"/settings/exports"}'::jsonb,
    hint = 'Open Settings → Export from the settings nav.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'export_engine';

-- Orientation tour: spotlight Settings → Import / Export instead of sidebar.
UPDATE demo_workflow_steps SET
    route = '/settings/imports',
    target_selector = '[data-tutorial-action="open-imports"]',
    placement = 'right'
WHERE workflow_key = 'crm_orientation_tour' AND step_key = 'imports';

UPDATE demo_workflow_steps SET
    route = '/settings/exports',
    target_selector = '[data-tutorial-action="open-exports"]',
    placement = 'right'
WHERE workflow_key = 'crm_orientation_tour' AND step_key = 'exports';
