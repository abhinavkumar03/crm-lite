package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/repository"
)

type fakeRepo struct {
	org *entity.Organization
}

func (f *fakeRepo) GetByID(_ context.Context, _ string) (*entity.Organization, error) {
	return f.org, nil
}

func (f *fakeRepo) Update(_ context.Context, _ string, p repository.ProfileUpdate) (*entity.Organization, error) {
	f.org.Name = p.Name
	f.org.LogoURL = p.LogoURL
	f.org.Industry = p.Industry
	f.org.CompanySize = p.CompanySize
	f.org.Country = p.Country
	f.org.Status = p.Status
	f.org.Settings = p.Settings
	return f.org, nil
}

func newRepo() *fakeRepo {
	return &fakeRepo{org: &entity.Organization{
		ID:     "org1",
		Name:   "Acme",
		Slug:   "acme",
		Plan:   "free",
		Status: "active",
	}}
}

func TestGetFillsDefaultsForEmptySettings(t *testing.T) {
	svc := New(newRepo())

	resp, err := svc.Get(context.Background(), "org1")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if resp.General.Timezone != "UTC" || resp.General.TimeFormat != "24h" {
		t.Errorf("expected general defaults, got %+v", resp.General)
	}
	if !resp.Automation.NotificationsEnabled || resp.Automation.DefaultChannel != "whatsapp" {
		t.Errorf("expected automation defaults, got %+v", resp.Automation)
	}
}

func TestUpdatePersistsSectionsAndNormalizes(t *testing.T) {
	repo := newRepo()
	svc := New(repo)

	name := "Acme Corp"
	resp, err := svc.Update(context.Background(), "org1", dto.UpdateSettingsRequest{
		Name: &name,
		General: &entity.GeneralSettings{
			Timezone:   "America/New_York",
			DateFormat: "MM/DD/YYYY",
			TimeFormat: "invalid", // should normalize to default
			Currency:   "EUR",
			Locale:     "en-GB",
			WeekStart:  "sunday",
		},
		Automation: &entity.AutomationSettings{
			NotificationsEnabled: false,
			DefaultChannel:       "email",
			DailyDigest:          true,
		},
	})
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}

	if resp.Name != "Acme Corp" {
		t.Errorf("expected name persisted, got %q", resp.Name)
	}
	if resp.General.TimeFormat != "24h" {
		t.Errorf("expected invalid time_format normalized to 24h, got %q", resp.General.TimeFormat)
	}
	if resp.General.Currency != "EUR" || resp.General.WeekStart != "sunday" {
		t.Errorf("unexpected general values: %+v", resp.General)
	}
	if resp.Automation.DefaultChannel != "email" || resp.Automation.NotificationsEnabled {
		t.Errorf("unexpected automation values: %+v", resp.Automation)
	}

	// The persisted blob must round-trip cleanly.
	var blob settingsBlob
	if err := json.Unmarshal(repo.org.Settings, &blob); err != nil {
		t.Fatalf("stored settings not valid JSON: %v", err)
	}
	if len(blob.General) == 0 || len(blob.Automation) == 0 {
		t.Errorf("expected both sections persisted, got %s", repo.org.Settings)
	}
}

func TestUpdatePartialLeavesOtherSectionUntouched(t *testing.T) {
	svc := New(newRepo())
	ctx := context.Background()

	// Save automation first.
	if _, err := svc.Update(ctx, "org1", dto.UpdateSettingsRequest{
		Automation: &entity.AutomationSettings{DefaultChannel: "email", NotificationsEnabled: true},
	}); err != nil {
		t.Fatalf("Update(automation) error: %v", err)
	}

	// Now save only general; automation must be preserved.
	resp, err := svc.Update(ctx, "org1", dto.UpdateSettingsRequest{
		General: &entity.GeneralSettings{Currency: "GBP"},
	})
	if err != nil {
		t.Fatalf("Update(general) error: %v", err)
	}
	if resp.Automation.DefaultChannel != "email" {
		t.Errorf("expected automation preserved, got %+v", resp.Automation)
	}
	if resp.General.Currency != "GBP" {
		t.Errorf("expected general currency GBP, got %q", resp.General.Currency)
	}
}
