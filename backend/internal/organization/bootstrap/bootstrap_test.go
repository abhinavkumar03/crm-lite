package bootstrap_test

import (
	"context"
	"testing"

	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/bootstrap"
	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
)

func TestDefaultRoles_OwnerHasAllPermissions(t *testing.T) {
	var owner *bootstrap.RoleSpec
	for i := range bootstrap.DefaultRoles {
		if bootstrap.DefaultRoles[i].Slug == "owner" {
			owner = &bootstrap.DefaultRoles[i]
			break
		}
	}
	if owner == nil {
		t.Fatal("owner role missing from DefaultRoles")
	}
	if !owner.AllPermissions {
		t.Fatal("organization creator (owner) must receive AllPermissions")
	}
}

func TestPermissionCatalog_IncludesModuleManage(t *testing.T) {
	found := false
	for _, p := range rbac.PermissionCatalog {
		if p.Key == rbac.PermModuleManage {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("%s must be in PermissionCatalog", rbac.PermModuleManage)
	}
}

func TestEnsurePermissionCatalog_NilDBPanicsNotCalled(t *testing.T) {
	// Compile-time smoke: method exists on Service.
	var _ interface {
		EnsurePermissionCatalog(context.Context) error
		RepairFullAccessRoles(context.Context) ([]string, error)
	} = (*bootstrap.Service)(nil)
}
