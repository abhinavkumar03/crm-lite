package dto

import (
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
)

// PermissionResponse is one catalog entry.
type PermissionResponse struct {
	ID          string  `json:"id"`
	Key         string  `json:"key"`
	Category    string  `json:"category"`
	Description *string `json:"description"`
}

// RoleSummary is the list-view representation of a role.
type RoleSummary struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	Description    *string   `json:"description"`
	IsSystem       bool      `json:"is_system"`
	HierarchyLevel int       `json:"hierarchy_level"`
	ParentRoleID   *string   `json:"parent_role_id,omitempty"`
	MemberCount    int       `json:"member_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// RoleDetail is a role plus its full permission matrix and ACL.
type RoleDetail struct {
	RoleSummary
	Permissions  []string            `json:"permissions"`
	ModuleAccess []rbac.ModuleAccess `json:"module_access"`
	FieldAccess  []rbac.FieldAccess  `json:"field_access"`
}

// CreateRoleRequest creates a custom (non-system) role.
type CreateRoleRequest struct {
	Name           string  `json:"name" validate:"required,max=100"`
	Slug           string  `json:"slug" validate:"required,max=100"`
	Description    *string `json:"description"`
	HierarchyLevel *int    `json:"hierarchy_level" validate:"omitempty,min=0,max=1000"`
}

// UpdateRoleRequest is a partial update of a role's metadata.
type UpdateRoleRequest struct {
	Name           *string `json:"name" validate:"omitempty,max=100"`
	Description    *string `json:"description"`
	HierarchyLevel *int    `json:"hierarchy_level" validate:"omitempty,min=0,max=1000"`
}

// SetPermissionsRequest replaces the role's global permission grants.
type SetPermissionsRequest struct {
	Permissions []string `json:"permissions" validate:"required,dive,max=120"`
}

// SetModuleAccessRequest replaces all module ACL rows for the role.
type SetModuleAccessRequest struct {
	Access []rbac.ModuleAccess `json:"access" validate:"required,dive"`
}

// SetFieldAccessRequest replaces all field ACL rows for the role.
type SetFieldAccessRequest struct {
	Access []rbac.FieldAccess `json:"access" validate:"required,dive"`
}

// MeResponse is the caller's effective RBAC context (role + permissions).
type MeResponse struct {
	RoleID       string              `json:"role_id"`
	RoleSlug     string              `json:"role_slug"`
	Permissions  []string            `json:"permissions"`
	ModuleAccess []rbac.ModuleAccess `json:"module_access"`
	FieldAccess  []rbac.FieldAccess  `json:"field_access"`
}
