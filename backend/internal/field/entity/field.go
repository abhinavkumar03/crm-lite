package entity

import "time"

// Field types supported by the engine. Kept in sync with the CHECK constraint
// on the fields table (migration 000003 + 000023).
const (
	TypeText         = "text"
	TypeTextarea     = "textarea"
	TypeEmail        = "email"
	TypePhone        = "phone"
	TypeNumber       = "number"
	TypeCurrency     = "currency"
	TypePercentage   = "percentage"
	TypeDate         = "date"
	TypeDatetime     = "datetime"
	TypeTime         = "time"
	TypeBoolean      = "boolean"
	TypeToggle       = "toggle"
	TypeDropdown     = "dropdown"
	TypeMultiselect  = "multiselect"
	TypeRadio        = "radio"
	TypeCheckbox     = "checkbox"
	TypeURL          = "url"
	TypeFile         = "file"
	TypeImage        = "image"
	TypeUser         = "user"
	TypeLookup       = "lookup"
	TypeFormula      = "formula"
	TypeJSON         = "json"
	TypeRichtext     = "richtext"
	TypeGST          = "gst"
	TypePAN          = "pan"
	TypeAddress      = "address"
	TypeAutoNumber   = "auto_number"
	TypeBarcode      = "barcode"
	TypeSerialNumber = "serial_number"
)

// Lock modes control when a field value becomes non-editable.
const (
	LockNever       = "never"
	LockAfterCreate = "after_create"
	LockAlways      = "always"
)

// AllTypes is the canonical set of valid field types.
var AllTypes = []string{
	TypeText, TypeTextarea, TypeEmail, TypePhone, TypeNumber, TypeCurrency,
	TypePercentage, TypeDate, TypeDatetime, TypeTime, TypeBoolean, TypeToggle,
	TypeDropdown, TypeMultiselect, TypeRadio, TypeCheckbox, TypeURL, TypeFile,
	TypeImage, TypeUser, TypeLookup, TypeFormula, TypeJSON, TypeRichtext,
	TypeGST, TypePAN, TypeAddress, TypeAutoNumber, TypeBarcode, TypeSerialNumber,
}

// AllLockModes is the set of valid lock_mode values.
var AllLockModes = []string{LockNever, LockAfterCreate, LockAlways}

// Field is a metadata-defined attribute of a module.
type Field struct {
	ID                string
	OrganizationID    string
	ModuleID          string
	APIName           string
	Label             string
	FieldType         string
	IsRequired        bool
	IsUnique          bool
	IsReadOnly        bool
	DefaultValue      *string
	Placeholder       *string
	Description       *string
	HelpText          *string
	MinLength         *int
	MaxLength         *int
	Regex             *string
	ValidationMessage *string
	Options           []byte // raw JSONB; normalized at the service layer
	LookupModuleID    *string
	SortOrder         int
	IsVisible         bool
	IsSearchable      bool
	IsFilterable      bool
	IsNullable        bool
	IsIndexed         bool
	IsSystem          bool
	LockMode          string
	EditableBy        string
	ViewableBy        string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// SortPosition is a single (id, sort_order) pair used for reordering.
type SortPosition struct {
	ID        string
	SortOrder int
}
