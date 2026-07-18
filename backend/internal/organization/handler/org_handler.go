package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/organization/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

type Handler struct {
	svc      *service.Service
	validate *validator.Validate
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc, validate: validator.New()}
}

func (h *Handler) ListMyOrgs(c *gin.Context) {
	orgs, err := h.svc.ListMyOrgs(c.Request.Context(), c.GetString("userID"))
	if err != nil {
		response.InternalServerError(c, "Unable to list organizations")
		return
	}
	response.OK(c, "Organizations fetched", orgs)
}

func (h *Handler) SwitchOrg(c *gin.Context) {
	var req dto.SwitchOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	err := h.svc.SwitchOrg(c.Request.Context(), c.GetString("userID"), req.OrganizationID)
	if errors.Is(err, tenant.ErrNotMember) {
		response.Forbidden(c, "Not a member of that organization")
		return
	}
	if err != nil {
		response.InternalServerError(c, "Unable to switch organization")
		return
	}
	response.OK(c, "Active organization updated", gin.H{"organization_id": req.OrganizationID})
}

func (h *Handler) CreateOrg(c *gin.Context) {
	var req dto.CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	id, err := h.svc.CreateOrg(c.Request.Context(), c.GetString("userID"), req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "Organization created", gin.H{"id": id})
}

func (h *Handler) ListMembers(c *gin.Context) {
	members, err := h.svc.ListMembers(c.Request.Context(), tenant.OrgID(c))
	if err != nil {
		response.InternalServerError(c, "Unable to list members")
		return
	}
	response.OK(c, "Members fetched", members)
}

func (h *Handler) Invite(c *gin.Context) {
	var req dto.CreateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "Validation failed", nil)
		return
	}
	inv, err := h.svc.Invite(c.Request.Context(), tenant.OrgID(c), c.GetString("userID"), req)
	if err != nil {
		response.InternalServerError(c, "Unable to create invitation")
		return
	}
	response.Created(c, "Invitation created", inv)
}

func (h *Handler) AcceptInvite(c *gin.Context) {
	var req dto.AcceptInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	userID, err := h.svc.AcceptInvite(c.Request.Context(), req)
	if errors.Is(err, service.ErrNotFound) {
		response.NotFound(c, "Invitation not found")
		return
	}
	if errors.Is(err, repository.ErrInviteExpired) {
		response.BadRequest(c, "Invitation expired", nil)
		return
	}
	if errors.Is(err, repository.ErrPasswordRequired) {
		response.BadRequest(c, "Name and password required for new users", nil)
		return
	}
	if err != nil {
		response.InternalServerError(c, "Unable to accept invitation")
		return
	}
	response.OK(c, "Invitation accepted", gin.H{"user_id": userID})
}

func (h *Handler) ListDepartments(c *gin.Context) {
	items, err := h.svc.ListDepartments(c.Request.Context(), tenant.OrgID(c))
	if err != nil {
		response.InternalServerError(c, "Unable to list departments")
		return
	}
	response.OK(c, "Departments fetched", items)
}

func (h *Handler) CreateDepartment(c *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	item, err := h.svc.CreateDepartment(c.Request.Context(), tenant.OrgID(c), req)
	if err != nil {
		response.BadRequest(c, "Unable to create department", nil)
		return
	}
	response.Created(c, "Department created", item)
}

func (h *Handler) ListTeams(c *gin.Context) {
	items, err := h.svc.ListTeams(c.Request.Context(), tenant.OrgID(c))
	if err != nil {
		response.InternalServerError(c, "Unable to list teams")
		return
	}
	response.OK(c, "Teams fetched", items)
}

func (h *Handler) CreateTeam(c *gin.Context) {
	var req dto.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	item, err := h.svc.CreateTeam(c.Request.Context(), tenant.OrgID(c), req)
	if err != nil {
		response.BadRequest(c, "Unable to create team", nil)
		return
	}
	response.Created(c, "Team created", item)
}

func (h *Handler) ListBranches(c *gin.Context) {
	items, err := h.svc.ListBranches(c.Request.Context(), tenant.OrgID(c))
	if err != nil {
		response.InternalServerError(c, "Unable to list branches")
		return
	}
	response.OK(c, "Branches fetched", items)
}

func (h *Handler) CreateBranch(c *gin.Context) {
	var req dto.CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil || h.validate.Struct(req) != nil {
		response.BadRequest(c, "Invalid payload", nil)
		return
	}
	item, err := h.svc.CreateBranch(c.Request.Context(), tenant.OrgID(c), req)
	if err != nil {
		response.BadRequest(c, "Unable to create branch", nil)
		return
	}
	response.Created(c, "Branch created", item)
}
