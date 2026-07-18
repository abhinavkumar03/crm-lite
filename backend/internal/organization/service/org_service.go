package service

import (
	"context"
	"errors"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/bootstrap"
	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
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
		Industry:    req.Industry,
		CompanySize: req.CompanySize,
		Country:     req.Country,
		LogoURL:     req.LogoURL,
	}
	if req.General != nil {
		opts.Timezone = req.General.Timezone
		opts.Currency = req.General.Currency
		opts.Locale = req.General.Locale
	}
	orgID, err := s.bootstrap.CreateOrganization(ctx, opts, userID)
	if err != nil {
		return "", err
	}
	// Ensure tenant cache does not keep a pre-create miss.
	_ = s.tenant.SetActiveOrganization(ctx, userID, orgID)
	return orgID, nil
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
