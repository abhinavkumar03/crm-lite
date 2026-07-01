package dto

import (
	"encoding/json"
	"time"
)

type RecentActivityResponse struct {
	ID string `json:"id"`

	EntityType string `json:"entity_type"`

	EntityID string `json:"entity_id"`

	Action string `json:"action"`

	Description string `json:"description"`

	Metadata json.RawMessage `json:"metadata,omitempty"`

	PerformedBy string `json:"performed_by"`

	CreatedAt time.Time `json:"created_at"`
}
