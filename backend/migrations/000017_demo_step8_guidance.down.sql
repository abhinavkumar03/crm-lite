-- Revert step copy to prior guided-metadata defaults (000014 / 000013).

UPDATE demo_workflow_steps SET
    description = 'From Tables, click a Tutorial Lead row to open the metadata-driven Record Workspace.',
    action_label = 'Open a record then Validate',
    hint = 'Open Tables, select Tutorial Lead, click a row.',
    failure_message = ''
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'record_workspace';

UPDATE demo_workflow_steps SET
    description = 'In the workspace Notes tab, add a note describing a follow-up.',
    action_label = 'Add note then Validate',
    hint = 'Open Notes tab and click Add.',
    failure_message = 'Add a note on a Tutorial Lead record.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'add_note';

UPDATE demo_workflow_steps SET
    description = 'Open the Timeline tab and confirm create/update/note events appear.',
    action_label = 'Check timeline then Validate',
    hint = '',
    failure_message = '',
    target_selector = NULL,
    placement = NULL,
    required_action = NULL
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'timeline';
