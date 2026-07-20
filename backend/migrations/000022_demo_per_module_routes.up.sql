-- Point demo/orientation table steps at per-module /m/{api_name} pages.

UPDATE demo_workflow_steps SET
    route = '/m/tutorial_lead',
    action_label = 'Open Tutorial Lead module',
    hint = 'Open Tutorial Lead from the workspace sidebar (or /m/tutorial_lead).'
WHERE workflow_key = 'crm_interactive_walkthrough'
  AND step_key IN ('record_workspace', 'add_note', 'timeline');

UPDATE demo_workflow_steps SET
    route = '/m/product_demo',
    action_label = 'Open Product Demo module',
    hint = 'Open Product Demo from the workspace sidebar.'
WHERE workflow_key = 'crm_interactive_walkthrough'
  AND step_key = 'product_demo_module';

UPDATE demo_workflow_steps SET
    route = '/dashboard',
    target_selector = '[data-tour="sidebar-nav"]',
    placement = 'right',
    title = 'Your modules',
    description = 'Each module in the workspace sidebar has its own page — view, edit, and delete records there.'
WHERE workflow_key = 'crm_orientation_tour' AND step_key = 'tables';
