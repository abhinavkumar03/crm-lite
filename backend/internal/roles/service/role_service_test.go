package service

import (
	"context"
	"testing"
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/roles/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/roles/entity"
)

type fakeRepo struct {
	roles map[string]*entity.Role
	perms map[string][]string
}

func newFake() *fakeRepo {
	return &fakeRepo{
		roles: map[string]*entity.Role{},
		perms: map[string][]string{},
	}
}

func (f *fakeRepo) ListPermissions(context.Context) ([]entity.Permission, error) {
	return []entity.Permission{{ID: "1", Key: "module.view", Category: "module"}}, nil
}
func (f *fakeRepo) ListRoles(context.Context, string) ([]entity.Role, []int, error) {
	return nil, nil, nil
}
func (f *fakeRepo) GetByID(_ context.Context, _, roleID string) (*entity.Role, error) {
	return f.roles[roleID], nil
}
func (f *fakeRepo) MemberCount(context.Context, string) (int, error) { return 0, nil }
func (f *fakeRepo) Create(_ context.Context, role *entity.Role) error {
	role.ID = "role-1"
	role.CreatedAt = time.Now()
	role.UpdatedAt = role.CreatedAt
	f.roles[role.ID] = role
	return nil
}
func (f *fakeRepo) Update(_ context.Context, role *entity.Role) error {
	f.roles[role.ID] = role
	return nil
}
func (f *fakeRepo) Delete(context.Context, string, string) error { return nil }
func (f *fakeRepo) PermissionKeys(_ context.Context, roleID string) ([]string, error) {
	return f.perms[roleID], nil
}
func (f *fakeRepo) SetPermissions(_ context.Context, roleID string, keys []string) error {
	f.perms[roleID] = keys
	return nil
}
func (f *fakeRepo) SetModuleAccess(context.Context, string, []rbac.ModuleAccess) error {
	return nil
}
func (f *fakeRepo) SetFieldAccess(context.Context, string, []rbac.FieldAccess) error {
	return nil
}
func (f *fakeRepo) SlugExists(context.Context, string, string) (bool, error) { return false, nil }

type fakeAccess struct{}

func (fakeAccess) ListModuleAccess(context.Context, string) ([]rbac.ModuleAccess, error) {
	return []rbac.ModuleAccess{}, nil
}
func (fakeAccess) ListFieldAccess(context.Context, string) ([]rbac.FieldAccess, error) {
	return []rbac.FieldAccess{}, nil
}

func TestCreateAndSetPermissions(t *testing.T) {
	repo := newFake()
	svc := New(repo, fakeAccess{}, nil)

	created, err := svc.Create(context.Background(), "org1", dto.CreateRoleRequest{
		Name: "Custom",
		Slug: "custom_role",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.Slug != "custom_role" {
		t.Errorf("unexpected slug %q", created.Slug)
	}

	detail, err := svc.SetPermissions(context.Background(), "org1", created.ID, dto.SetPermissionsRequest{
		Permissions: []string{"module.view", "record.view"},
	})
	if err != nil {
		t.Fatalf("SetPermissions: %v", err)
	}
	if len(detail.Permissions) != 2 {
		t.Errorf("expected 2 perms, got %v", detail.Permissions)
	}
}

func TestCreateRejectsBadSlug(t *testing.T) {
	svc := New(newFake(), fakeAccess{}, nil)
	_, err := svc.Create(context.Background(), "org1", dto.CreateRoleRequest{
		Name: "Bad",
		Slug: "NOT VALID",
	})
	if err != ErrInvalidSlug {
		t.Fatalf("expected ErrInvalidSlug, got %v", err)
	}
}

func TestDeleteSystemRole(t *testing.T) {
	repo := newFake()
	repo.roles["sys"] = &entity.Role{ID: "sys", OrganizationID: "org1", IsSystem: true}
	svc := New(repo, fakeAccess{}, nil)

	if err := svc.Delete(context.Background(), "org1", "sys"); err != ErrSystemRole {
		t.Fatalf("expected ErrSystemRole, got %v", err)
	}
}
