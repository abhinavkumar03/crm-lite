UPDATE demo_workflow_steps SET
    route = '/forms',
    action_label = 'Create record then Validate',
    hint = ''
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'create_record';

UPDATE demo_workflow_steps SET
    route = '/tables'
WHERE workflow_key = 'crm_interactive_walkthrough'
  AND step_key IN ('record_workspace', 'add_note', 'timeline', 'product_demo_module');

UPDATE demo_workflow_steps SET
    route = NULL,
    target_selector = '[data-tour="nav-forms"]'
WHERE workflow_key = 'crm_orientation_tour' AND step_key = 'forms';

UPDATE demo_workflow_steps SET
    route = NULL,
    target_selector = '[data-tour="nav-tables"]'
WHERE workflow_key = 'crm_orientation_tour' AND step_key = 'tables';
