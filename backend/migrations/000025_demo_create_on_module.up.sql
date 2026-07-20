-- Create records from module listing pages; Form Designer is preview-only.

UPDATE demo_workflow_steps SET
    route = '/m/tutorial_lead?create=1',
    target_selector = '[data-tutorial-action="create-record"]',
    action_label = 'Add Tutorial Lead',
    hint = 'Open Tutorial Lead → Add Tutorial Lead → Create record.',
    title = 'Create a Tutorial Lead record',
    description = 'From the Tutorial Lead module page, open Add and save a record. Forms settings is preview-only.'
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'create_record';

UPDATE demo_workflow_steps SET
    route = '/settings/forms',
    target_selector = '[data-tutorial-action="open-forms"]',
    title = 'Form Designer',
    description = 'Preview how create forms are generated from field metadata. Create real records from each module Add button.',
    placement = 'right'
WHERE workflow_key = 'crm_orientation_tour' AND step_key = 'forms';
