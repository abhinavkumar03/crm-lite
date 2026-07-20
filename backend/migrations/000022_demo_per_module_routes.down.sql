UPDATE demo_workflow_steps SET
    route = '/settings/tables',
    action_label = 'Open Tables (Settings)',
    hint = 'Open a module from the workspace sidebar, or Settings → Tables.'
WHERE workflow_key = 'crm_interactive_walkthrough'
  AND step_key IN ('record_workspace', 'add_note', 'timeline');

UPDATE demo_workflow_steps SET
    route = '/settings/tables',
    action_label = 'Open Tables (Settings)',
    hint = 'Open a module from the workspace sidebar, or Settings → Tables.'
WHERE workflow_key = 'crm_interactive_walkthrough'
  AND step_key = 'product_demo_module';

UPDATE demo_workflow_steps SET
    route = '/settings/tables',
    target_selector = '[data-tutorial-action="open-tables"]',
    placement = 'right',
    title = 'Dynamic tables',
    description = 'Browse records from Settings → Tables, or jump into a module from the workspace sidebar.'
WHERE workflow_key = 'crm_orientation_tour' AND step_key = 'tables';
