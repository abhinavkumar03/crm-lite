// Package rbac implements role-based access control: a permission catalog,
// middleware that loads a role's grants into the request context, Require()
// guards for route handlers, and helpers for module-/field-level ACL.
package rbac

// Global permission keys. Kept in sync with the seed catalog
// (internal/seed/seeders/catalog.go).
const (
	PermModuleView       = "module.view"
	PermModuleManage     = "module.manage"
	PermFieldManage      = "field.manage"
	PermRecordView       = "record.view"
	PermRecordCreate     = "record.create"
	PermRecordUpdate     = "record.update"
	PermRecordDelete     = "record.delete"
	PermImportRun        = "import.run"
	PermExportRun        = "export.run"
	PermAutomationManage = "automation.manage"
	PermValidationManage = "validation.manage"
	PermSettingsManage     = "settings.manage"
	PermUserManage         = "user.manage"
	PermRoleManage         = "role.manage"
	PermOrganizationManage = "organization.manage"
	PermAnalyticsView      = "analytics.view"
	PermNotificationView             = "notification.view"
	PermNotificationSend             = "notification.send"
	PermNotificationTemplatesManage  = "notification.templates.manage"
	PermCommunicationProvidersManage = "communication.providers.manage"
)

// ModuleAction is a CRUD verb checked against role_module_access.
type ModuleAction string

const (
	ActionView   ModuleAction = "view"
	ActionCreate ModuleAction = "create"
	ActionUpdate ModuleAction = "update"
	ActionDelete ModuleAction = "delete"
)

// FieldAccess levels for role_field_access.
const (
	FieldHidden = "hidden"
	FieldRead   = "read"
	FieldWrite  = "write"
)

// ContextPermissions is the Gin context key holding the caller's permission keys.
const ContextPermissions = "permissions"
