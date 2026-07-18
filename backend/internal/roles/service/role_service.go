package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/roles/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/roles/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/cache"
)

var (
	ErrNotFound    = errors.New("role not found")
	ErrSlugTaken   = errors.New("role slug already exists")
	ErrSystemRole  = errors.New("system roles cannot be deleted")
	ErrInvalidSlug = errors.New("slug must be lowercase letters, numbers and underscores")
	ErrHasMembers  = errors.New("role still has active members")
)

var slugRE = regexp.MustCompile(`^[a-z][a-z0-9_]{1,99}$`)

// Repository is the persistence contract for roles.
type Repository interface {
	ListPermissions(ctx context.Context) ([]entity.Permission, error)
	ListRoles(ctx context.Context, orgID string) ([]entity.Role, []int, error)
	GetByID(ctx context.Context, orgID, roleID string) (*entity.Role, error)
	MemberCount(ctx context.Context, roleID string) (int, error)
	Create(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, orgID, roleID string) error
	PermissionKeys(ctx context.Context, roleID string) ([]string, error)
	SetPermissions(ctx context.Context, roleID string, keys []string) error
	SetModuleAccess(ctx context.Context, roleID string, access []rbac.ModuleAccess) error
	SetFieldAccess(ctx context.Context, roleID string, access []rbac.FieldAccess) error
	SlugExists(ctx context.Context, orgID, slug string) (bool, error)
}

// AccessReader loads ACL rows (satisfied by *rbac.Guard).
type AccessReader interface {
	ListModuleAccess(ctx context.Context, roleID string) ([]rbac.ModuleAccess, error)
	ListFieldAccess(ctx context.Context, roleID string) ([]rbac.FieldAccess, error)
}

type Service struct {
	repo   Repository
	access AccessReader
	cache  *cache.Cache
}

func New(repo Repository, access AccessReader, c *cache.Cache) *Service {
	return &Service{repo: repo, access: access, cache: c}
}

func (s *Service) ListPermissions(ctx context.Context) ([]dto.PermissionResponse, error) {
	items, err := s.repo.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.PermissionResponse, 0, len(items))
	for _, p := range items {
		out = append(out, dto.PermissionResponse{
			ID: p.ID, Key: p.Key, Category: p.Category, Description: p.Description,
		})
	}
	return out, nil
}

func (s *Service) List(ctx context.Context, orgID string) ([]dto.RoleSummary, error) {
	roles, counts, err := s.repo.ListRoles(ctx, orgID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.RoleSummary, 0, len(roles))
	for i := range roles {
		out = append(out, toSummary(&roles[i], counts[i]))
	}
	return out, nil
}

func (s *Service) Get(ctx context.Context, orgID, roleID string) (*dto.RoleDetail, error) {
	role, err := s.repo.GetByID(ctx, orgID, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrNotFound
	}
	return s.detail(ctx, role)
}

func (s *Service) Create(ctx context.Context, orgID string, req dto.CreateRoleRequest) (*dto.RoleDetail, error) {
	slug := strings.TrimSpace(strings.ToLower(req.Slug))
	if !slugRE.MatchString(slug) {
		return nil, ErrInvalidSlug
	}
	exists, err := s.repo.SlugExists(ctx, orgID, slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlugTaken
	}

	level := 80
	if req.HierarchyLevel != nil {
		level = *req.HierarchyLevel
	}
	role := &entity.Role{
		OrganizationID: orgID,
		Name:           strings.TrimSpace(req.Name),
		Slug:           slug,
		Description:    req.Description,
		HierarchyLevel: level,
	}
	if err := s.repo.Create(ctx, role); err != nil {
		return nil, err
	}
	return s.detail(ctx, role)
}

func (s *Service) Update(ctx context.Context, orgID, roleID string, req dto.UpdateRoleRequest) (*dto.RoleDetail, error) {
	role, err := s.repo.GetByID(ctx, orgID, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrNotFound
	}
	if req.Name != nil {
		role.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		role.Description = req.Description
	}
	if req.HierarchyLevel != nil {
		role.HierarchyLevel = *req.HierarchyLevel
	}
	if err := s.repo.Update(ctx, role); err != nil {
		return nil, err
	}
	return s.detail(ctx, role)
}

func (s *Service) Delete(ctx context.Context, orgID, roleID string) error {
	role, err := s.repo.GetByID(ctx, orgID, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrNotFound
	}
	if role.IsSystem {
		return ErrSystemRole
	}
	n, err := s.repo.MemberCount(ctx, roleID)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrHasMembers
	}
	oldModules, _ := s.access.ListModuleAccess(ctx, roleID)
	oldFields, _ := s.access.ListFieldAccess(ctx, roleID)

	if err := s.repo.Delete(ctx, orgID, roleID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	s.cache.InvalidateRole(ctx, roleID)
	s.invalidateModuleACL(ctx, roleID, oldModules, nil)
	s.invalidateFieldACL(ctx, roleID, oldFields, nil)
	return nil
}

func (s *Service) SetPermissions(ctx context.Context, orgID, roleID string, req dto.SetPermissionsRequest) (*dto.RoleDetail, error) {
	role, err := s.repo.GetByID(ctx, orgID, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrNotFound
	}
	if err := s.repo.SetPermissions(ctx, roleID, req.Permissions); err != nil {
		return nil, err
	}
	s.cache.Delete(ctx, cache.PermissionsKey(roleID))
	return s.detail(ctx, role)
}

func (s *Service) SetModuleAccess(ctx context.Context, orgID, roleID string, req dto.SetModuleAccessRequest) (*dto.RoleDetail, error) {
	role, err := s.repo.GetByID(ctx, orgID, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrNotFound
	}

	old, _ := s.access.ListModuleAccess(ctx, roleID)
	if err := s.repo.SetModuleAccess(ctx, roleID, req.Access); err != nil {
		return nil, err
	}
	s.invalidateModuleACL(ctx, roleID, old, req.Access)
	return s.detail(ctx, role)
}

func (s *Service) SetFieldAccess(ctx context.Context, orgID, roleID string, req dto.SetFieldAccessRequest) (*dto.RoleDetail, error) {
	role, err := s.repo.GetByID(ctx, orgID, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, ErrNotFound
	}

	old, _ := s.access.ListFieldAccess(ctx, roleID)
	if err := s.repo.SetFieldAccess(ctx, roleID, req.Access); err != nil {
		return nil, err
	}
	neu, _ := s.access.ListFieldAccess(ctx, roleID)
	s.invalidateFieldACL(ctx, roleID, old, neu)
	return s.detail(ctx, role)
}

func (s *Service) invalidateModuleACL(ctx context.Context, roleID string, old, neu []rbac.ModuleAccess) {
	seen := map[string]struct{}{}
	keys := make([]string, 0, len(old)+len(neu))
	for _, list := range [][]rbac.ModuleAccess{old, neu} {
		for _, a := range list {
			if _, ok := seen[a.ModuleID]; ok {
				continue
			}
			seen[a.ModuleID] = struct{}{}
			keys = append(keys, cache.Key("rbac", "module", roleID, a.ModuleID))
		}
	}
	s.cache.Delete(ctx, keys...)
}

func (s *Service) invalidateFieldACL(ctx context.Context, roleID string, old, neu []rbac.FieldAccess) {
	seen := map[string]struct{}{}
	ids := make([]string, 0)
	for _, list := range [][]rbac.FieldAccess{old, neu} {
		for _, a := range list {
			if a.ModuleID == "" {
				continue
			}
			if _, ok := seen[a.ModuleID]; ok {
				continue
			}
			seen[a.ModuleID] = struct{}{}
			ids = append(ids, a.ModuleID)
		}
	}
	s.cache.InvalidateRoleFieldAccess(ctx, roleID, ids...)
}

// Me returns the caller's effective RBAC context.
func (s *Service) Me(ctx context.Context, roleID, roleSlug string, permissions []string) (*dto.MeResponse, error) {
	moduleAccess, err := s.access.ListModuleAccess(ctx, roleID)
	if err != nil {
		return nil, err
	}
	fieldAccess, err := s.access.ListFieldAccess(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if permissions == nil {
		permissions = []string{}
	}
	return &dto.MeResponse{
		RoleID:       roleID,
		RoleSlug:     roleSlug,
		Permissions:  permissions,
		ModuleAccess: moduleAccess,
		FieldAccess:  fieldAccess,
	}, nil
}

func (s *Service) detail(ctx context.Context, role *entity.Role) (*dto.RoleDetail, error) {
	count, err := s.repo.MemberCount(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	perms, err := s.repo.PermissionKeys(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	moduleAccess, err := s.access.ListModuleAccess(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	fieldAccess, err := s.access.ListFieldAccess(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	if perms == nil {
		perms = []string{}
	}
	return &dto.RoleDetail{
		RoleSummary:  toSummary(role, count),
		Permissions:  perms,
		ModuleAccess: moduleAccess,
		FieldAccess:  fieldAccess,
	}, nil
}

func toSummary(role *entity.Role, count int) dto.RoleSummary {
	return dto.RoleSummary{
		ID:             role.ID,
		Name:           role.Name,
		Slug:           role.Slug,
		Description:    role.Description,
		IsSystem:       role.IsSystem,
		HierarchyLevel: role.HierarchyLevel,
		ParentRoleID:   role.ParentRoleID,
		MemberCount:    count,
		CreatedAt:      role.CreatedAt,
		UpdatedAt:      role.UpdatedAt,
	}
}
