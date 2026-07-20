-- Point Forms / Tables demo & orientation steps at Settings routes.

UPDATE demo_workflow_steps SET
    route = '/settings/forms',
    target_selector = '[data-tutorial-action="create-record"]',
    action_label = 'Create record (Settings → Forms)',
    hint = 'Open Settings → Forms, then Create record.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'create_record';

UPDATE demo_workflow_steps SET
    route = '/settings/tables',
    action_label = 'Open Tables (Settings)',
    hint = 'Open a module from the workspace sidebar, or Settings → Tables.'
WHERE workflow_key = 'crm_interactive_walkthrough'
  AND step_key IN ('record_workspace', 'add_note', 'timeline', 'product_demo_module');

UPDATE demo_workflow_steps SET
    route = '/settings/forms',
    target_selector = '[data-tutorial-action="open-forms"]',
    placement = 'right'
WHERE workflow_key = 'crm_orientation_tour' AND step_key = 'forms';

UPDATE demo_workflow_steps SET
    route = '/settings/tables',
    target_selector = '[data-tutorial-action="open-tables"]',
    placement = 'right'
WHERE workflow_key = 'crm_orientation_tour' AND step_key = 'tables';
