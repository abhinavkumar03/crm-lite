package entity

import "time"

// Template categories for reusable message content.
const (
	CategorySales     = "sales"
	CategoryFollowUp  = "follow_up"
	CategoryWelcome   = "welcome"
	CategoryProposal  = "proposal"
	CategoryInvoice   = "invoice"
	CategoryReminder  = "reminder"
	CategoryQuotation = "quotation"
	CategoryMarketing = "marketing"
	CategorySupport   = "support"
	CategoryCustom    = "custom"
)

const (
	TemplateStatusDraft     = "draft"
	TemplateStatusPublished = "published"
)

// Template is an org-scoped reusable notification body.
type Template struct {
	ID             string
	OrganizationID string
	Channel        string
	Name           string
	Category       string
	Subject        *string
	Body           string
	BodyHTML       *string
	Variables      []byte // JSONB array of declared keys
	IsActive       bool
	Status         string
	Version        int
	CreatedBy      *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
