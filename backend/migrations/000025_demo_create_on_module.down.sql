-- Restore create_record / forms orientation routes to Settings → Forms create flow.

UPDATE demo_workflow_steps SET
    route = '/settings/forms',
    target_selector = '[data-tutorial-action="create-record"]',
    action_label = 'Create record (Settings → Forms)',
    hint = 'Open Settings → Forms, then Create record.',
    title = 'Create a record',
    description = 'Open Settings → Forms and create a Tutorial Lead record.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'create_record';

UPDATE demo_workflow_steps SET
    route = '/settings/forms',
    target_selector = '[data-tutorial-action="open-forms"]',
    title = 'Dynamic forms',
    description = 'Create records from Settings → Forms — forms are generated from module metadata.',
    placement = 'right'
WHERE workflow_key = 'crm_orientation_tour' AND step_key = 'forms';
