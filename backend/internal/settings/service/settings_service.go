package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/repository"
)

var ErrNotFound = errors.New("organization not found")

type Repository interface {
	GetByID(ctx context.Context, orgID string) (*entity.Organization, error)
	Update(ctx context.Context, orgID string, p repository.ProfileUpdate) (*entity.Organization, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

type settingsBlob struct {
	General    json.RawMessage `json:"general,omitempty"`
	Automation json.RawMessage `json:"automation,omitempty"`
}

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
	logo := org.LogoURL
	if req.LogoURL != nil {
		v := strings.TrimSpace(*req.LogoURL)
		if v == "" {
			logo = nil
		} else {
			logo = &v
		}
	}
	desc := org.Description
	if req.Description != nil {
		v := strings.TrimSpace(*req.Description)
		if v == "" {
			desc = nil
		} else {
			desc = &v
		}
	}
	industry := org.Industry
	if req.Industry != nil {
		industry = req.Industry
	}
	size := org.CompanySize
	if req.CompanySize != nil {
		size = req.CompanySize
	}
	country := org.Country
	if req.Country != nil {
		country = req.Country
	}
	status := org.Status
	if status == "" {
		status = "active"
	}
	if req.Status != nil {
		status = *req.Status
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

	updated, err := s.repo.Update(ctx, orgID, repository.ProfileUpdate{
		Name:        name,
		LogoURL:     logo,
		Description: desc,
		Industry:    industry,
		CompanySize: size,
		Country:     country,
		Status:      status,
		Settings:    blob,
	})
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, ErrNotFound
	}

	resp := toResponse(updated, general, automation)
	return &resp, nil
}

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
	status := org.Status
	if status == "" {
		status = "active"
	}
	return dto.SettingsResponse{
		ID:               org.ID,
		Name:             org.Name,
		Slug:             org.Slug,
		Plan:             org.Plan,
		SubscriptionPlan: org.Plan,
		LogoURL:          org.LogoURL,
		Description:      org.Description,
		Industry:         org.Industry,
		CompanySize:      org.CompanySize,
		Country:          org.Country,
		Status:           status,
		General:          general,
		Automation:       automation,
		UpdatedAt:        org.UpdatedAt,
	}
}
