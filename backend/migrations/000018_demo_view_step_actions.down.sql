UPDATE demo_workflow_steps SET
    required_action = 'navigate',
    hint = 'Open Tables, select Tutorial Lead, click a row.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'record_workspace';

UPDATE demo_workflow_steps SET
    required_action = 'acknowledge',
    hint = ''
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'timeline';

UPDATE demo_workflow_steps SET
    required_action = 'acknowledge',
    hint = ''
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'product_demo_module';

UPDATE demo_workflow_steps SET
    required_action = 'acknowledge',
    hint = ''
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'roles_glance';
