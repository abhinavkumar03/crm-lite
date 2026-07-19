DELETE FROM demo_workflow_steps
WHERE workflow_key = 'crm_interactive_walkthrough'
  AND step_key IN (
    'validation_rules', 'import_engine', 'export_engine',
    'automation_settings', 'settings_sweep'
  );

UPDATE demo_workflow_steps
SET sort_order = 12
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'completion';

UPDATE demo_workflows
SET description = 'Hands-on sandbox tutorial covering metadata modules, records, workspace, and settings.',
    duration_min = 15
WHERE workflow_key = 'crm_interactive_walkthrough';
