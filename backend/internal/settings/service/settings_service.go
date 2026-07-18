package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/entity"
)

// ErrNotFound is returned when the organization row cannot be resolved.
var ErrNotFound = errors.New("organization not found")

// Repository is the persistence contract for organization settings.
type Repository interface {
	GetByID(ctx context.Context, orgID string) (*entity.Organization, error)
	Update(ctx context.Context, orgID, name string, settings []byte) (*entity.Organization, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

// settingsBlob is the on-disk shape of organizations.settings.
type settingsBlob struct {
	General    json.RawMessage `json:"general,omitempty"`
	Automation json.RawMessage `json:"automation,omitempty"`
}

// Get returns the organization's settings, filling in defaults for any section
// or key that has never been saved.
func (s *Service) Get(ctx context.Context, orgID string) (*dto.SettingsResponse, error) {
	org, err := s.repo.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, ErrNotFound
	}

	general, automation := parseBlob(org.Settings)
	resp := toResponse(org, general, automation)
	return &resp, nil
}

// Update applies the provided sections (partial) and persists them. Sections the
// client omits are left untouched.
func (s *Service) Update(ctx context.Context, orgID string, req dto.UpdateSettingsRequest) (*dto.SettingsResponse, error) {
	org, err := s.repo.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, ErrNotFound
	}

	general, automation := parseBlob(org.Settings)

	name := org.Name
	if req.Name != nil {
		name = *req.Name
	}
	if req.General != nil {
		general = normalizeGeneral(*req.General)
	}
	if req.Automation != nil {
		automation = normalizeAutomation(*req.Automation)
	}

	blob, err := marshalBlob(general, automation)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(ctx, orgID, name, blob)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, ErrNotFound
	}

	resp := toResponse(updated, general, automation)
	return &resp, nil
}

// parseBlob decodes the settings JSONB, layering saved values over defaults so a
// newly-added preference key is always populated.
func parseBlob(raw []byte) (entity.GeneralSettings, entity.AutomationSettings) {
	general := entity.DefaultGeneral()
	automation := entity.DefaultAutomation()

	if len(raw) == 0 {
		return general, automation
	}

	var b settingsBlob
	if err := json.Unmarshal(raw, &b); err != nil {
		return general, automation
	}
	if len(b.General) > 0 {
		_ = json.Unmarshal(b.General, &general)
	}
	if len(b.Automation) > 0 {
		_ = json.Unmarshal(b.Automation, &automation)
	}
	return normalizeGeneral(general), normalizeAutomation(automation)
}

func marshalBlob(general entity.GeneralSettings, automation entity.AutomationSettings) ([]byte, error) {
	g, err := json.Marshal(general)
	if err != nil {
		return nil, err
	}
	a, err := json.Marshal(automation)
	if err != nil {
		return nil, err
	}
	return json.Marshal(settingsBlob{General: g, Automation: a})
}

// normalizeGeneral guards enum-like fields and fills blanks with defaults so bad
// input can never corrupt the stored preferences.
func normalizeGeneral(g entity.GeneralSettings) entity.GeneralSettings {
	d := entity.DefaultGeneral()
	if g.Timezone == "" {
		g.Timezone = d.Timezone
	}
	if g.DateFormat == "" {
		g.DateFormat = d.DateFormat
	}
	if g.TimeFormat != "12h" && g.TimeFormat != "24h" {
		g.TimeFormat = d.TimeFormat
	}
	if g.Currency == "" {
		g.Currency = d.Currency
	}
	if g.Locale == "" {
		g.Locale = d.Locale
	}
	if g.WeekStart != "sunday" && g.WeekStart != "monday" {
		g.WeekStart = d.WeekStart
	}
	return g
}

func normalizeAutomation(a entity.AutomationSettings) entity.AutomationSettings {
	if a.DefaultChannel != "whatsapp" && a.DefaultChannel != "email" {
		a.DefaultChannel = entity.DefaultAutomation().DefaultChannel
	}
	return a
}

func toResponse(org *entity.Organization, general entity.GeneralSettings, automation entity.AutomationSettings) dto.SettingsResponse {
	return dto.SettingsResponse{
		ID:         org.ID,
		Name:       org.Name,
		Slug:       org.Slug,
		Plan:       org.Plan,
		General:    general,
		Automation: automation,
		UpdatedAt:  org.UpdatedAt,
	}
}
