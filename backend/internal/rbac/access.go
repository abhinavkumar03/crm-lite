package rbac

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/cache"
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
	FieldID  string `json:"field_id"`
	ModuleID string `json:"module_id,omitempty"`
	Access   string `json:"access"` // hidden | read | write
}

// moduleACLEntry is the cached shape for a single role+module ACL lookup.
// Found=false means "no row" (unrestricted).
type moduleACLEntry struct {
	Found     bool `json:"found"`
	CanView   bool `json:"can_view"`
	CanCreate bool `json:"can_create"`
	CanUpdate bool `json:"can_update"`
	CanDelete bool `json:"can_delete"`
}

// ModuleAllowed reports whether the caller's role may perform action on the
// module. No ACL row means unrestricted (true).
func (g *Guard) ModuleAllowed(ctx context.Context, c *gin.Context, moduleID string, action ModuleAction) (bool, error) {
	roleID := tenant.RoleID(c)
	if roleID == "" {
		return true, nil
	}

	entry, err := g.moduleACL(ctx, roleID, moduleID)
	if err != nil {
		return false, err
	}
	if !entry.Found {
		return true, nil
	}

	switch action {
	case ActionView:
		return entry.CanView, nil
	case ActionCreate:
		return entry.CanCreate, nil
	case ActionUpdate:
		return entry.CanUpdate, nil
	case ActionDelete:
		return entry.CanDelete, nil
	default:
		return false, nil
	}
}

func (g *Guard) moduleACL(ctx context.Context, roleID, moduleID string) (moduleACLEntry, error) {
	key := cache.Key("rbac", "module", roleID, moduleID)
	var cached moduleACLEntry
	if g.cache.GetJSON(ctx, key, &cached) {
		return cached, nil
	}

	var entry moduleACLEntry
	err := g.db.QueryRow(ctx, `
		SELECT can_view, can_create, can_update, can_delete
		FROM role_module_access
		WHERE role_id = $1 AND module_id = $2
	`, roleID, moduleID).Scan(&entry.CanView, &entry.CanCreate, &entry.CanUpdate, &entry.CanDelete)

	if errors.Is(err, pgx.ErrNoRows) {
		entry.Found = false
		g.cache.SetJSON(ctx, key, entry, cache.TTLShort)
		return entry, nil
	}
	if err != nil {
		return entry, err
	}
	entry.Found = true
	g.cache.SetJSON(ctx, key, entry, cache.TTLShort)
	return entry, nil
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
// module. Used to strip/reject payload keys on create/update. Cached per
// (role, module).
func (g *Guard) FieldAccessByAPIName(ctx context.Context, roleID, moduleID string) (map[string]string, error) {
	out := map[string]string{}
	if roleID == "" {
		return out, nil
	}

	key := cache.FieldAccessKey(roleID, moduleID)
	if g.cache.GetJSON(ctx, key, &out) {
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	g.cache.SetJSON(ctx, key, out, cache.TTLShort)
	return out, nil
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
		SELECT rfa.field_id::text, f.module_id::text, rfa.access
		FROM role_field_access rfa
		JOIN fields f ON f.id = rfa.field_id
		WHERE rfa.role_id = $1
		ORDER BY rfa.field_id
	`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]FieldAccess, 0)
	for rows.Next() {
		var f FieldAccess
		if err := rows.Scan(&f.FieldID, &f.ModuleID, &f.Access); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	return items, rows.Err()
}
