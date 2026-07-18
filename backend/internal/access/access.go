package access

import (
	"context"
	"fmt"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/repository"
)

// Actor is the authenticated principal inside an organization.
type Actor struct {
	UserID         string
	OrgID          string
	RoleSlug       string
	HierarchyLevel int
	DepartmentID   *string
	TeamID         *string
	SubordinateIDs []string // includes self when loaded from hierarchy CTE
}

// Service builds reusable visibility filters for dynamic records.
type Service struct {
	orgRepo *repository.Repository
}

func New(orgRepo *repository.Repository) *Service {
	return &Service{orgRepo: orgRepo}
}

// LoadActor resolves membership meta + subordinate tree for visibility checks.
func (s *Service) LoadActor(ctx context.Context, orgID, userID string) (Actor, error) {
	dept, team, slug, level, err := s.orgRepo.MemberMeta(ctx, orgID, userID)
	if err != nil {
		return Actor{}, err
	}
	subs, err := s.orgRepo.SubordinateUserIDs(ctx, orgID, userID)
	if err != nil {
		return Actor{}, err
	}
	if len(subs) == 0 {
		subs = []string{userID}
	}
	return Actor{
		UserID:         userID,
		OrgID:          orgID,
		RoleSlug:       slug,
		HierarchyLevel: level,
		DepartmentID:   dept,
		TeamID:         team,
		SubordinateIDs: subs,
	}, nil
}

func (a Actor) SeesAllInOrg() bool {
	switch a.RoleSlug {
	case "owner", "super_admin", "admin":
		return true
	}
	return a.HierarchyLevel <= 20
}

// VisibilitySQL appends a boolean predicate (no leading AND) that restricts
// rows the actor may see. startArg is the next $N placeholder index.
func VisibilitySQL(actor Actor, startArg int) (sql string, args []any, next int) {
	if actor.SeesAllInOrg() {
		return "TRUE", nil, startArg
	}

	args = []any{actor.UserID}
	userPh := fmt.Sprintf("$%d", startArg)
	next = startArg + 1

	subs := actor.SubordinateIDs
	if len(subs) == 0 {
		subs = []string{actor.UserID}
	}
	subsPh := make([]string, len(subs))
	for i, id := range subs {
		subsPh[i] = fmt.Sprintf("$%d", next)
		args = append(args, id)
		next++
	}
	subsList := strings.Join(subsPh, ",")

	parts := []string{
		fmt.Sprintf("visibility IN ('organization','public')"),
		fmt.Sprintf("(owner_id = %s OR assigned_to = %s OR created_by = %s)", userPh, userPh, userPh),
		fmt.Sprintf("(visibility IN ('manager','hierarchy') AND owner_id IN (%s))", subsList),
	}

	if actor.DepartmentID != nil && *actor.DepartmentID != "" {
		parts = append(parts, fmt.Sprintf(
			"(visibility = 'department' AND department_id = $%d)", next,
		))
		args = append(args, *actor.DepartmentID)
		next++
	}
	if actor.TeamID != nil && *actor.TeamID != "" {
		parts = append(parts, fmt.Sprintf(
			"(visibility = 'team' AND team_id = $%d)", next,
		))
		args = append(args, *actor.TeamID)
		next++
	}
	// private/owner still covered by owner/assignee/created_by match above
	_ = "private"

	return "(" + strings.Join(parts, " OR ") + ")", args, next
}

// CanViewRecord evaluates a single row (used by Get).
func CanViewRecord(actor Actor, ownerID, assignedTo, createdBy, deptID, teamID, visibility *string) bool {
	if actor.SeesAllInOrg() {
		return true
	}
	vis := "organization"
	if visibility != nil && *visibility != "" {
		vis = *visibility
	}
	if vis == "organization" || vis == "public" {
		return true
	}
	if eq(ownerID, actor.UserID) || eq(assignedTo, actor.UserID) || eq(createdBy, actor.UserID) {
		return true
	}
	if vis == "manager" || vis == "hierarchy" {
		oid := deref(ownerID)
		for _, id := range actor.SubordinateIDs {
			if id == oid {
				return true
			}
		}
	}
	if vis == "department" && actor.DepartmentID != nil && eq(deptID, *actor.DepartmentID) {
		return true
	}
	if vis == "team" && actor.TeamID != nil && eq(teamID, *actor.TeamID) {
		return true
	}
	return false
}

func eq(p *string, v string) bool {
	return p != nil && *p == v
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
