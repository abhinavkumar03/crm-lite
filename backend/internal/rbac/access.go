package rbac

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

// ModuleAccess is one row of role_module_access.
type ModuleAccess struct {
	ModuleID  string `json:"module_id"`
	CanView   bool   `json:"can_view"`
	CanCreate bool   `json:"can_create"`
	CanUpdate bool   `json:"can_update"`
	CanDelete bool   `json:"can_delete"`
}

// FieldAccess is one row of role_field_access.
type FieldAccess struct {
	FieldID string `json:"field_id"`
	Access  string `json:"access"` // hidden | read | write
}

// ModuleAllowed reports whether the caller's role may perform action on the
// module. No ACL row means unrestricted (true).
func (g *Guard) ModuleAllowed(ctx context.Context, c *gin.Context, moduleID string, action ModuleAction) (bool, error) {
	roleID := tenant.RoleID(c)
	if roleID == "" {
		return true, nil
	}

	var canView, canCreate, canUpdate, canDelete bool
	err := g.db.QueryRow(ctx, `
		SELECT can_view, can_create, can_update, can_delete
		FROM role_module_access
		WHERE role_id = $1 AND module_id = $2
	`, roleID, moduleID).Scan(&canView, &canCreate, &canUpdate, &canDelete)

	if errors.Is(err, pgx.ErrNoRows) {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	switch action {
	case ActionView:
		return canView, nil
	case ActionCreate:
		return canCreate, nil
	case ActionUpdate:
		return canUpdate, nil
	case ActionDelete:
		return canDelete, nil
	default:
		return false, nil
	}
}

// FieldAccessMap returns field_id → access for the caller's role within a
// module. Only explicit ACL rows are returned; missing keys mean "write".
func (g *Guard) FieldAccessMap(ctx context.Context, roleID, moduleID string) (map[string]string, error) {
	out := map[string]string{}
	if roleID == "" {
		return out, nil
	}

	rows, err := g.db.Query(ctx, `
		SELECT rfa.field_id::text, rfa.access
		FROM role_field_access rfa
		JOIN fields f ON f.id = rfa.field_id
		WHERE rfa.role_id = $1 AND f.module_id = $2
	`, roleID, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fieldID, access string
		if err := rows.Scan(&fieldID, &access); err != nil {
			return nil, err
		}
		out[fieldID] = access
	}
	return out, rows.Err()
}

// FieldAccessByAPIName returns api_name → access for the caller's role within a
// module. Used to strip/reject payload keys on create/update.
func (g *Guard) FieldAccessByAPIName(ctx context.Context, roleID, moduleID string) (map[string]string, error) {
	out := map[string]string{}
	if roleID == "" {
		return out, nil
	}

	rows, err := g.db.Query(ctx, `
		SELECT f.api_name, rfa.access
		FROM role_field_access rfa
		JOIN fields f ON f.id = rfa.field_id
		WHERE rfa.role_id = $1 AND f.module_id = $2
	`, roleID, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var apiName, access string
		if err := rows.Scan(&apiName, &access); err != nil {
			return nil, err
		}
		out[apiName] = access
	}
	return out, rows.Err()
}

// ListModuleAccess returns every explicit module ACL row for a role.
func (g *Guard) ListModuleAccess(ctx context.Context, roleID string) ([]ModuleAccess, error) {
	rows, err := g.db.Query(ctx, `
		SELECT module_id::text, can_view, can_create, can_update, can_delete
		FROM role_module_access
		WHERE role_id = $1
		ORDER BY module_id
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ModuleAccess, 0)
	for rows.Next() {
		var m ModuleAccess
		if err := rows.Scan(&m.ModuleID, &m.CanView, &m.CanCreate, &m.CanUpdate, &m.CanDelete); err != nil {
			return nil, err
		}
		items = append(items, m)
	}
	return items, rows.Err()
}

// ListFieldAccess returns every explicit field ACL row for a role.
func (g *Guard) ListFieldAccess(ctx context.Context, roleID string) ([]FieldAccess, error) {
	rows, err := g.db.Query(ctx, `
		SELECT field_id::text, access
		FROM role_field_access
		WHERE role_id = $1
		ORDER BY field_id
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]FieldAccess, 0)
	for rows.Next() {
		var f FieldAccess
		if err := rows.Scan(&f.FieldID, &f.Access); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	return items, rows.Err()
}
