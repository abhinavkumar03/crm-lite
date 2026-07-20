package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/bootstrap"
	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/repository"
	settingsentity "github.com/abhinavkumar03/crm-lite/backend/internal/settings/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrConflict      = errors.New("conflict")
	ErrForbidden     = errors.New("forbidden")
	ErrLastWorkspace = errors.New("cannot delete last workspace")
)

type Service struct {
	repo      *repository.Repository
	bootstrap *bootstrap.Service
	tenant    *tenant.Resolver
}

func New(repo *repository.Repository, boot *bootstrap.Service, tenantResolver *tenant.Resolver) *Service {
	return &Service{repo: repo, bootstrap: boot, tenant: tenantResolver}
}

func (s *Service) ListMyOrgs(ctx context.Context, userID string) ([]dto.OrgSummary, error) {
	return s.repo.ListOrgsForUser(ctx, userID)
}

func (s *Service) SwitchOrg(ctx context.Context, userID, orgID string) error {
	return s.tenant.SetActiveOrganization(ctx, userID, orgID)
}

func (s *Service) CreateOrg(ctx context.Context, userID string, req dto.CreateOrgRequest) (string, error) {
	opts := bootstrap.CreateOptions{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Industry:    req.Industry,
		CompanySize: req.CompanySize,
		Country:     req.Country,
		LogoURL:     req.LogoURL,
	}
	if req.General != nil {
		opts.Timezone = req.General.Timezone
		opts.Currency = req.General.Currency
		opts.Locale = req.General.Locale
		opts.DateFormat = req.General.DateFormat
	}
	orgID, err := s.bootstrap.CreateOrganization(ctx, opts, userID)
	if err != nil {
		return "", err
	}
	_ = s.tenant.SetActiveOrganization(ctx, userID, orgID)
	return orgID, nil
}

func (s *Service) GetCurrentOrg(ctx context.Context, orgID string) (*dto.OrgDetail, error) {
	row, err := s.repo.GetOrgByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNotFound
	}
	return orgRowToDetail(row), nil
}

func (s *Service) UpdateCurrentOrg(ctx context.Context, orgID string, req dto.UpdateOrgRequest) (*dto.OrgDetail, error) {
	row, err := s.repo.GetOrgByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, ErrNotFound
	}

	name := row.Name
	if req.Name != nil {
		name = strings.TrimSpace(*req.Name)
	}
	logo := row.LogoURL
	if req.LogoURL != nil {
		v := strings.TrimSpace(*req.LogoURL)
		if v == "" {
			logo = nil
		} else {
			logo = &v
		}
	}
	desc := row.Description
	if req.Description != nil {
		v := strings.TrimSpace(*req.Description)
		if v == "" {
			desc = nil
		} else {
			desc = &v
		}
	}
	industry := row.Industry
	if req.Industry != nil {
		industry = req.Industry
	}
	size := row.CompanySize
	if req.CompanySize != nil {
		size = req.CompanySize
	}
	country := row.Country
	if req.Country != nil {
		country = req.Country
	}

	general := parseGeneral(row.Settings)
	if req.General != nil {
		mergeGeneral(&general, *req.General)
	}
	settings, err := marshalSettings(row.Settings, general)
	if err != nil {
		return nil, err
	}

	updated, err := s.repo.UpdateOrg(ctx, orgID, repository.OrgProfileUpdate{
		Name:        name,
		LogoURL:     logo,
		Description: desc,
		Industry:    industry,
		CompanySize: size,
		Country:     country,
		Settings:    settings,
	})
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, ErrNotFound
	}
	return orgRowToDetail(updated), nil
}

func (s *Service) SoftDeleteCurrentOrg(ctx context.Context, userID, orgID, roleSlug string) error {
	if roleSlug != "owner" && roleSlug != "super_admin" && roleSlug != "admin" {
		return ErrForbidden
	}
	n, err := s.repo.CountActiveMemberships(ctx, userID)
	if err != nil {
		return err
	}
	if n <= 1 {
		return ErrLastWorkspace
	}

	memberIDs, err := s.repo.SoftDeleteOrg(ctx, orgID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	for _, uid := range memberIDs {
		s.tenant.InvalidateMembershipCache(ctx, uid)
	}
	return nil
}

func (s *Service) ListMembers(ctx context.Context, orgID string) ([]dto.MemberResponse, error) {
	return s.repo.ListMembers(ctx, orgID)
}

func (s *Service) Invite(ctx context.Context, orgID, invitedBy string, req dto.CreateInviteRequest) (*dto.InviteResponse, error) {
	return s.repo.CreateInvitation(
		ctx, orgID, strings.ToLower(strings.TrimSpace(req.Email)), req.RoleID, invitedBy,
		req.ManagerUserID, req.DepartmentID, req.TeamID,
	)
}

func (s *Service) AcceptInvite(ctx context.Context, req dto.AcceptInviteRequest) (string, error) {
	inv, err := s.repo.GetPendingInvite(ctx, req.Token)
	if err != nil {
		return "", err
	}
	if inv == nil {
		return "", ErrNotFound
	}
	userID, err := s.repo.AcceptInvite(ctx, inv, strings.TrimSpace(req.Name), req.Password)
	if errors.Is(err, repository.ErrInviteExpired) || errors.Is(err, repository.ErrPasswordRequired) {
		return "", err
	}
	return userID, err
}

func (s *Service) ListDepartments(ctx context.Context, orgID string) ([]dto.StructureItem, error) {
	return s.repo.ListDepartments(ctx, orgID)
}
func (s *Service) CreateDepartment(ctx context.Context, orgID string, req dto.CreateDepartmentRequest) (*dto.StructureItem, error) {
	return s.repo.CreateDepartment(ctx, orgID, req)
}
func (s *Service) ListTeams(ctx context.Context, orgID string) ([]dto.StructureItem, error) {
	return s.repo.ListTeams(ctx, orgID)
}
func (s *Service) CreateTeam(ctx context.Context, orgID string, req dto.CreateTeamRequest) (*dto.StructureItem, error) {
	return s.repo.CreateTeam(ctx, orgID, req)
}
func (s *Service) ListBranches(ctx context.Context, orgID string) ([]dto.StructureItem, error) {
	return s.repo.ListBranches(ctx, orgID)
}
func (s *Service) CreateBranch(ctx context.Context, orgID string, req dto.CreateBranchRequest) (*dto.StructureItem, error) {
	return s.repo.CreateBranch(ctx, orgID, req)
}

func orgRowToDetail(row *repository.OrgRow) *dto.OrgDetail {
	g := parseGeneral(row.Settings)
	return &dto.OrgDetail{
		ID:          row.ID,
		Name:        row.Name,
		Slug:        row.Slug,
		Plan:        row.Plan,
		LogoURL:     row.LogoURL,
		Description: row.Description,
		Industry:    row.Industry,
		CompanySize: row.CompanySize,
		Country:     row.Country,
		Status:      row.Status,
		CreatedBy:   row.CreatedBy,
		General: dto.OrgGeneralPrefs{
			Timezone:   g.Timezone,
			Currency:   g.Currency,
			Locale:     g.Locale,
			DateFormat: g.DateFormat,
			TimeFormat: g.TimeFormat,
			WeekStart:  g.WeekStart,
		},
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

type settingsBlob struct {
	General    json.RawMessage `json:"general,omitempty"`
	Automation json.RawMessage `json:"automation,omitempty"`
}

func parseGeneral(raw []byte) settingsentity.GeneralSettings {
	g := settingsentity.DefaultGeneral()
	if len(raw) == 0 {
		return g
	}
	var blob settingsBlob
	if err := json.Unmarshal(raw, &blob); err != nil || len(blob.General) == 0 {
		return g
	}
	_ = json.Unmarshal(blob.General, &g)
	def := settingsentity.DefaultGeneral()
	if g.Timezone == "" {
		g.Timezone = def.Timezone
	}
	if g.Currency == "" {
		g.Currency = def.Currency
	}
	if g.DateFormat == "" {
		g.DateFormat = def.DateFormat
	}
	if g.TimeFormat == "" {
		g.TimeFormat = def.TimeFormat
	}
	if g.Locale == "" {
		g.Locale = def.Locale
	}
	if g.WeekStart == "" {
		g.WeekStart = def.WeekStart
	}
	return g
}

func mergeGeneral(dst *settingsentity.GeneralSettings, src dto.OrgGeneralPrefs) {
	if v := strings.TrimSpace(src.Timezone); v != "" {
		dst.Timezone = v
	}
	if v := strings.TrimSpace(src.Currency); v != "" {
		dst.Currency = v
	}
	if v := strings.TrimSpace(src.Locale); v != "" {
		dst.Locale = v
	}
	if v := strings.TrimSpace(src.DateFormat); v != "" {
		dst.DateFormat = v
	}
	if v := strings.TrimSpace(src.TimeFormat); v != "" {
		dst.TimeFormat = v
	}
	if v := strings.TrimSpace(src.WeekStart); v != "" {
		dst.WeekStart = v
	}
}

func marshalSettings(existing []byte, general settingsentity.GeneralSettings) ([]byte, error) {
	var blob settingsBlob
	if len(existing) > 0 {
		_ = json.Unmarshal(existing, &blob)
	}
	gRaw, err := json.Marshal(general)
	if err != nil {
		return nil, err
	}
	blob.General = gRaw
	if len(blob.Automation) == 0 {
		aRaw, err := json.Marshal(settingsentity.DefaultAutomation())
		if err != nil {
			return nil, err
		}
		blob.Automation = aRaw
	}
	return json.Marshal(blob)
}
