-- Stronger copy + action labels for workspace / note / timeline mentored steps.

UPDATE demo_workflow_steps SET
    description = 'We open a Tutorial Lead record workspace for you. Skim Overview, then continue.',
    action_label = 'Open Tutorial Lead workspace',
    hint = 'You should be on /tables/:moduleId/:recordId — then Validate.',
    failure_message = 'Open a Tutorial Lead record workspace first.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'record_workspace';

UPDATE demo_workflow_steps SET
    description = 'On the Notes tab, a follow-up note is ready. Click Add to save it and continue.',
    action_label = 'Open Notes tab',
    hint = 'Click the highlighted Add button on the Notes tab.',
    failure_message = 'Add a note on a Tutorial Lead record (Notes tab → Add).',
    target_selector = '[data-tutorial-action="add-note"]',
    placement = 'left',
    required_action = 'create_resource',
    allow_selectors = '["[data-tutorial-action=\"add-note\"]"]'::jsonb
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'add_note';

UPDATE demo_workflow_steps SET
    description = 'Open the Timeline tab and confirm create / note activity appears for this record.',
    action_label = 'Open Timeline tab',
    hint = 'Open Timeline, review activities, then Validate.',
    failure_message = 'Open the Timeline tab on a Tutorial Lead record.',
    target_selector = '[data-tutorial-action="open-timeline-tab"]',
    placement = 'bottom',
    required_action = 'navigate'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'timeline';
