package rbac

// PermissionDef is one entry in the global permission catalog.
type PermissionDef struct {
	Key      string
	Category string
	Desc     string
}

// PermissionCatalog is the full set of permission keys the app understands.
// Kept in sync with route guards (permissions.go) and seeders.
var PermissionCatalog = []PermissionDef{
	{"module.view", "module", "View modules and navigation"},
	{"module.manage", "module", "Create, update, reorder and delete modules"},
	{"field.manage", "field", "Create, update and delete fields"},
	{"record.view", "record", "View records"},
	{"record.create", "record", "Create records"},
	{"record.update", "record", "Update records"},
	{"record.delete", "record", "Delete records"},
	{"import.run", "import", "Run data imports"},
	{"export.run", "export", "Run data exports"},
	{"automation.manage", "automation", "Create and manage automation rules"},
	{"workflow.view", "workflow", "View workflows and builder metadata"},
	{"workflow.create", "workflow", "Create workflow drafts"},
	{"workflow.edit", "workflow", "Edit workflow drafts and definitions"},
	{"workflow.delete", "workflow", "Archive or delete workflows"},
	{"workflow.publish", "workflow", "Publish and disable workflows"},
	{"workflow.execute", "workflow", "Manually run workflows"},
	{"workflow.logs.view", "workflow", "View workflow execution logs"},
	{"validation.manage", "validation", "Create and manage validation rules"},
	{"settings.manage", "settings", "Manage organization settings"},
	{"user.manage", "user", "Invite and manage users"},
	{"role.manage", "role", "Create and manage roles and permissions"},
	{"organization.manage", "organization", "Create organizations and manage membership structure"},
	{"analytics.view", "analytics", "View dashboards and analytics"},
	{"notification.view", "notification", "View notification center and delivery history"},
	{"notification.send", "notification", "Compose, send, schedule, retry and cancel notifications"},
	{"notification.templates.manage", "notification", "Create and manage notification templates"},
	{"communication.providers.manage", "notification", "Manage email/WhatsApp providers and sender identities"},
}
