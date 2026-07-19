-- Expand Interactive CRM Walkthrough curriculum (metadata-only steps).

-- Shift completion to the end after new steps.
UPDATE demo_workflow_steps
SET sort_order = 20
WHERE workflow_key = 'crm_interactive_walkthrough' AND step_key = 'completion';

INSERT INTO demo_workflow_steps (
    workflow_key, step_key, sort_order, title, description, why_it_matters, how_it_works,
    expected_result, route, target_selector, action_label, validator_key, validator_params,
    is_skippable, required_action, hint, placement
) VALUES
(
    'crm_interactive_walkthrough', 'validation_rules', 12,
    'Open validation settings',
    'Visit Settings → Validation to see how rules attach to dynamic fields.',
    'Validation keeps bad data out of the CRM without hardcoding per-module logic.',
    'Rules are metadata evaluated by the validation engine on record write.',
    'Validation settings page is visible.',
    '/settings/validation', '[data-tour="nav-settings"]', 'Open Validation',
    'route_visited', '{"route":"/settings/validation"}', TRUE,
    'navigate', 'Open Settings → Validation from the sidebar.', 'right'
),
(
    'crm_interactive_walkthrough', 'import_engine', 13,
    'Explore Import',
    'Open the Import page. Upload is optional — understand column mapping and async jobs.',
    'Bulk onboarding is a core CRM workflow; the import engine validates every row.',
    'POST /imports creates a job; the worker maps, validates, and inserts records.',
    'Import page loads.',
    '/imports', '[data-tour="nav-imports"]', 'Open Import',
    'route_visited', '{"route":"/imports"}', TRUE,
    'navigate', 'Click Import in the sidebar.', 'right'
),
(
    'crm_interactive_walkthrough', 'export_engine', 14,
    'Explore Export',
    'Open Export to see filtered CSV/Excel jobs and history.',
    'Recruiters and ops teams need portable extracts without raw SQL.',
    'POST /exports enqueues a worker job that serializes module records.',
    'Export page loads.',
    '/exports', '[data-tour="nav-exports"]', 'Open Export',
    'route_visited', '{"route":"/exports"}', TRUE,
    'navigate', 'Click Export in the sidebar.', 'right'
),
(
    'crm_interactive_walkthrough', 'automation_settings', 15,
    'Glance at automation',
    'Open Settings → Automation (and Notifications) to see WhatsApp/email pipeline controls.',
    'CRM value compounds when records trigger outbound messages automatically.',
    'Automation settings live on the organization; notifications enqueue to Asynq workers.',
    'Automation settings page loads.',
    '/settings/automation', NULL, 'Open Automation',
    'route_visited', '{"route":"/settings/automation"}', TRUE,
    'navigate', 'Open Settings → Automation.', 'center'
),
(
    'crm_interactive_walkthrough', 'settings_sweep', 16,
    'Settings center sweep',
    'Open the Settings home to see modules, fields, roles, members, and data tools in one place.',
    'A single admin surface reduces hunting across the app.',
    'Settings pages compose existing metadata engines under one nav.',
    'Settings overview is visible.',
    '/settings', '[data-tour="nav-settings"]', 'Open Settings',
    'route_visited', '{"route":"/settings"}', TRUE,
    'navigate', 'Click Settings in the sidebar.', 'right'
)
ON CONFLICT (workflow_key, step_key) DO NOTHING;

UPDATE demo_workflows
SET description = 'Hands-on sandbox tutorial covering metadata modules, fields, records, workspace, validation, import/export, automation, and settings.',
    duration_min = 20
WHERE workflow_key = 'crm_interactive_walkthrough';
