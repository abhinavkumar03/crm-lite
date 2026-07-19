-- Orientation tour as a metadata workflow (acknowledge-only steps).
-- The FE TourProvider prefers this catalogue when available.

INSERT INTO demo_workflows (workflow_key, name, description, version, duration_min)
VALUES (
    'crm_orientation_tour',
    'CRM Orientation Tour',
    '60-second spotlight walkthrough of workspace navigation. No sandbox; progress still uses tour_progress.',
    1,
    1
)
ON CONFLICT (workflow_key) DO NOTHING;

INSERT INTO demo_workflow_steps (
    workflow_key, step_key, sort_order, title, description, why_it_matters, how_it_works,
    expected_result, route, target_selector, action_label, validator_key, validator_params,
    is_skippable, required_action, placement
) VALUES
(
    'crm_orientation_tour', 'welcome', 1,
    'Welcome to CRM Lite',
    'Take a 60-second tour of the workspace. You can skip anytime and restart later from your profile menu.',
    'Orientation reduces first-login friction.',
    'Client spotlight; progress persisted in tour_progress.',
    'User understands they can skip or continue.',
    '/dashboard', NULL, 'Next', 'acknowledge', '{}', TRUE, 'acknowledge', 'center'
),
(
    'crm_orientation_tour', 'navigation', 2,
    'Your workspace',
    'Everything lives in this sidebar — dynamic modules, forms, tables, and data tools.',
    'Navigation is the map of the product.',
    'Sidebar items are stable data-tour anchors.',
    'Sidebar is highlighted.',
    '/dashboard', '[data-tour="sidebar-nav"]', 'Next', 'acknowledge', '{}', TRUE, 'acknowledge', 'right'
),
(
    'crm_orientation_tour', 'forms', 3,
    'Dynamic forms',
    'Create records with forms generated from module metadata and a backend validation schema.',
    'Forms are metadata-driven — no hardcoded screens per object.',
    'DynamicForm renders ModuleField metadata.',
    'Forms nav highlighted.',
    NULL, '[data-tour="nav-forms"]', 'Next', 'acknowledge', '{}', TRUE, 'acknowledge', 'right'
),
(
    'crm_orientation_tour', 'tables', 4,
    'Dynamic tables',
    'Metadata-driven tables with sorting, filtering, and saved views that persist per module.',
    'List views are how users scan CRM data at scale.',
    'DynamicTable + saved views.',
    'Tables nav highlighted.',
    NULL, '[data-tour="nav-tables"]', 'Next', 'acknowledge', '{}', TRUE, 'acknowledge', 'right'
),
(
    'crm_orientation_tour', 'imports', 5,
    'Import engine',
    'Bring in CSV or Excel files. Columns are auto-mapped and rows are validated and processed in the background.',
    'Bulk load is essential for CRM adoption.',
    'Import jobs via asynq worker.',
    'Import nav highlighted.',
    NULL, '[data-tour="nav-imports"]', 'Next', 'acknowledge', '{}', TRUE, 'acknowledge', 'right'
),
(
    'crm_orientation_tour', 'exports', 6,
    'Export engine',
    'Export any module to CSV or Excel — instantly or as a background job — and reuse saved export templates.',
    'Portable extracts without SQL.',
    'Export jobs via asynq worker.',
    'Export nav highlighted.',
    NULL, '[data-tour="nav-exports"]', 'Next', 'acknowledge', '{}', TRUE, 'acknowledge', 'right'
),
(
    'crm_orientation_tour', 'search', 7,
    'Global search',
    'Jump to any record fast. Search spans dynamic module data from anywhere in the app.',
    'Power users live in search.',
    'Topbar global search.',
    'Search highlighted.',
    '/dashboard', '[data-tour="global-search"]', 'Next', 'acknowledge', '{}', TRUE, 'acknowledge', 'bottom'
),
(
    'crm_orientation_tour', 'notifications', 8,
    'Notifications',
    'Delivery updates for WhatsApp and email automations show up here as they are sent.',
    'Outbound messaging needs visibility.',
    'Notification bell + list.',
    'Bell highlighted.',
    NULL, '[data-tour="notification-bell"]', 'Next', 'acknowledge', '{}', TRUE, 'acknowledge', 'bottom'
),
(
    'crm_orientation_tour', 'restart', 9,
    'Restart anytime',
    'Re-run this tour from your profile menu. For a hands-on sandbox tutorial, use Explore CRM.',
    'Users forget UI — restart is cheap.',
    'UserMenu → Take a tour.',
    'Profile menu highlighted.',
    NULL, '[data-tour="user-menu"]', 'Next', 'acknowledge', '{}', TRUE, 'acknowledge', 'bottom'
),
(
    'crm_orientation_tour', 'done', 10,
    'You are ready',
    'That is the lay of the land. Launch Explore CRM when you want a guided sandbox walkthrough.',
    'Clear handoff to the interactive tutorial.',
    'Tour marks completed in tour_progress.',
    'Tour finished.',
    NULL, NULL, 'Finish', 'acknowledge', '{}', TRUE, 'acknowledge', 'center'
)
ON CONFLICT (workflow_key, step_key) DO NOTHING;
