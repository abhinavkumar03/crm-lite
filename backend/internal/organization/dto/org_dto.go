package dto

import "time"

type OrgSummary struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	LogoURL     *string `json:"logo_url,omitempty"`
	Description *string `json:"description,omitempty"`
	RoleSlug    string  `json:"role_slug"`
	IsActive    bool    `json:"is_active"`
}

// OrgGeneralPrefs are locale defaults stored under settings.general.
type OrgGeneralPrefs struct {
	Timezone   string `json:"timezone" validate:"omitempty,max=80"`
	Currency   string `json:"currency" validate:"omitempty,max=10"`
	Locale     string `json:"locale" validate:"omitempty,max=20"`
	DateFormat string `json:"date_format" validate:"omitempty,max=40"`
	TimeFormat string `json:"time_format" validate:"omitempty,oneof=12h 24h"`
	WeekStart  string `json:"week_start" validate:"omitempty,oneof=sunday monday"`
}

// OrgDetail is the active workspace profile returned by GET /organizations/current.
type OrgDetail struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Plan        string           `json:"plan"`
	LogoURL     *string          `json:"logo_url,omitempty"`
	Description *string          `json:"description,omitempty"`
	Industry    *string          `json:"industry,omitempty"`
	CompanySize *string          `json:"company_size,omitempty"`
	Country     *string          `json:"country,omitempty"`
	Status      string           `json:"status"`
	CreatedBy   *string          `json:"created_by,omitempty"`
	General     OrgGeneralPrefs  `json:"general"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// UpdateOrgRequest partially updates the active workspace.
type UpdateOrgRequest struct {
	Name        *string          `json:"name" validate:"omitempty,min=2,max=200"`
	LogoURL     *string          `json:"logo_url"`
	Description *string          `json:"description" validate:"omitempty,max=2000"`
	Industry    *string          `json:"industry" validate:"omitempty,max=120"`
	CompanySize *string          `json:"company_size" validate:"omitempty,max=40"`
	Country     *string          `json:"country" validate:"omitempty,max=80"`
	General     *OrgGeneralPrefs `json:"general"`
}

// CreateOrgGeneralPrefs are locale defaults applied into settings.general on create.
type CreateOrgGeneralPrefs struct {
	Timezone   string `json:"timezone" validate:"omitempty,max=80"`
	Currency   string `json:"currency" validate:"omitempty,max=10"`
	Locale     string `json:"locale" validate:"omitempty,max=20"`
	DateFormat string `json:"date_format" validate:"omitempty,max=40"`
}

// CreateOrgRequest creates a workspace and bootstraps roles/modules for the caller.
type CreateOrgRequest struct {
	Name        string                 `json:"name" validate:"required,min=2,max=200"`
	Slug        string                 `json:"slug" validate:"omitempty,max=120"`
	Description string                 `json:"description" validate:"omitempty,max=2000"`
	Industry    string                 `json:"industry" validate:"omitempty,max=120"`
	CompanySize string                 `json:"company_size" validate:"omitempty,max=40"`
	Country     string                 `json:"country" validate:"omitempty,max=80"`
	LogoURL     string                 `json:"logo_url" validate:"omitempty,max=500"`
	General     *CreateOrgGeneralPrefs `json:"general"`
}

type SwitchOrgRequest struct {
	OrganizationID string `json:"organization_id" validate:"required,uuid"`
}

type CreateInviteRequest struct {
	Email         string  `json:"email" validate:"required,email"`
	RoleID        string  `json:"role_id" validate:"required,uuid"`
	ManagerUserID *string `json:"manager_user_id" validate:"omitempty,uuid"`
	DepartmentID  *string `json:"department_id" validate:"omitempty,uuid"`
	TeamID        *string `json:"team_id" validate:"omitempty,uuid"`
}

type InviteResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
	// SimulatedEmail is logged for portfolio demos (no real SMTP).
	SimulatedEmail string `json:"simulated_email_body"`
}

type AcceptInviteRequest struct {
	Token    string `json:"token" validate:"required"`
	Name     string `json:"name" validate:"omitempty,min=1,max=255"`
	Password string `json:"password" validate:"omitempty,min=8,max=128"`
}

type MemberResponse struct {
	UserID         string  `json:"user_id"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	RoleID         *string `json:"role_id,omitempty"`
	RoleSlug       string  `json:"role_slug"`
	ManagerUserID  *string `json:"manager_user_id,omitempty"`
	DepartmentID   *string `json:"department_id,omitempty"`
	TeamID         *string `json:"team_id,omitempty"`
	BranchID       *string `json:"branch_id,omitempty"`
	Designation    *string `json:"designation,omitempty"`
	HierarchyLevel int     `json:"hierarchy_level"`
	Status         string  `json:"status"`
}

type StructureItem struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	Location     *string `json:"location,omitempty"`
	DepartmentID *string `json:"department_id,omitempty"`
}

type CreateDepartmentRequest struct {
	Name        string  `json:"name" validate:"required,max=120"`
	Description *string `json:"description"`
}

type CreateTeamRequest struct {
	Name         string  `json:"name" validate:"required,max=120"`
	Description  *string `json:"description"`
	DepartmentID *string `json:"department_id" validate:"omitempty,uuid"`
}

type CreateBranchRequest struct {
	Name     string  `json:"name" validate:"required,max=120"`
	Location *string `json:"location"`
}
