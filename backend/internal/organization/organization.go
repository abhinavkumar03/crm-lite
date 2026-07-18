package organization

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/bootstrap"
	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/handler"
	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

type Module struct {
	Handler *handler.Handler
	auth    gin.HandlerFunc
	org     gin.HandlerFunc
	load    gin.HandlerFunc
	guard   *rbac.Guard
}

func NewModule(
	db *pgxpool.Pool,
	resolver *tenant.Resolver,
	auth, org, load gin.HandlerFunc,
	guard *rbac.Guard,
) *Module {
	repo := repository.New(db)
	boot := bootstrap.New(db)
	svc := service.New(repo, boot, resolver)
	h := handler.New(svc)
	return &Module{Handler: h, auth: auth, org: org, load: load, guard: guard}
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	// Public accept (token-based).
	api.POST("/invitations/accept", m.Handler.AcceptInvite)

	me := api.Group("/me")
	me.Use(m.auth)
	me.GET("/organizations", m.Handler.ListMyOrgs)
	me.POST("/organizations/switch", m.Handler.SwitchOrg)

	// Create org does not require existing tenant membership.
	orgs := api.Group("/organizations")
	orgs.Use(m.auth)
	orgs.POST("", m.Handler.CreateOrg)

	scoped := api.Group("")
	scoped.Use(m.auth, m.org, m.load)
	scoped.GET("/organizations/members", m.guard.Require(rbac.PermUserManage), m.Handler.ListMembers)
	scoped.POST("/organizations/invitations", m.guard.Require(rbac.PermUserManage), m.Handler.Invite)

	scoped.GET("/departments", m.Handler.ListDepartments)
	scoped.POST("/departments", m.guard.Require(rbac.PermOrganizationManage), m.Handler.CreateDepartment)
	scoped.GET("/teams", m.Handler.ListTeams)
	scoped.POST("/teams", m.guard.Require(rbac.PermOrganizationManage), m.Handler.CreateTeam)
	scoped.GET("/branches", m.Handler.ListBranches)
	scoped.POST("/branches", m.guard.Require(rbac.PermOrganizationManage), m.Handler.CreateBranch)
}
