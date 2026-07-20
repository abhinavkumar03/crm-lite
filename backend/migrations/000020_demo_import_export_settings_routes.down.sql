UPDATE demo_workflow_steps SET
    route = '/imports',
    target_selector = '[data-tour="nav-imports"]',
    action_label = 'Open Import',
    validator_params = '{"route":"/imports"}'::jsonb,
    hint = 'Click Import in the sidebar.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'import_engine';

UPDATE demo_workflow_steps SET
    route = '/exports',
    target_selector = '[data-tour="nav-exports"]',
    action_label = 'Open Export',
    validator_params = '{"route":"/exports"}'::jsonb,
    hint = 'Click Export in the sidebar.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'export_engine';

UPDATE demo_workflow_steps SET
    route = NULL,
    target_selector = '[data-tour="nav-imports"]',
    placement = 'right'
WHERE workflow_key = 'crm_orientation_tour' AND step_key = 'imports';

UPDATE demo_workflow_steps SET
    route = NULL,
    target_selector = '[data-tour="nav-exports"]',
    placement = 'right'
WHERE workflow_key = 'crm_orientation_tour' AND step_key = 'exports';
