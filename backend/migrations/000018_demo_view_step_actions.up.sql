-- Mark inspect/view steps so clients wait for explicit Continue (no silent auto-advance).

UPDATE demo_workflow_steps SET
    required_action = 'inspect',
    hint = 'Skim the record Overview, then press Continue.',
    is_skippable = TRUE
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'record_workspace';

UPDATE demo_workflow_steps SET
    required_action = 'inspect',
    target_selector = '[data-tutorial-action="open-timeline-tab"]',
    hint = 'Review Timeline activities, then press Continue.',
    is_skippable = TRUE
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'timeline';

UPDATE demo_workflow_steps SET
    required_action = 'inspect',
    hint = 'Browse Product Demo rows and fields, then press Continue.',
    is_skippable = TRUE
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'product_demo_module';

UPDATE demo_workflow_steps SET
    required_action = 'navigate',
    hint = 'Open Roles, look around, then press Continue.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'roles_glance';
