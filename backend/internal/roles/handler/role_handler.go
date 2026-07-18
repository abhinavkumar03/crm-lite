package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/rbac"
	"github.com/abhinavkumar03/crm-lite/backend/internal/roles/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/roles/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

type RoleHandler struct {
	service *service.Service
}

func New(service *service.Service) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) ListPermissions(c *gin.Context) {
	items, err := h.service.ListPermissions(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Unable to fetch permissions")
		return
	}
	response.OK(c, "Permissions fetched successfully", items)
}

func (h *RoleHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context(), tenant.OrgID(c))
	if err != nil {
		response.InternalServerError(c, "Unable to fetch roles")
		return
	}
	response.OK(c, "Roles fetched successfully", items)
}

func (h *RoleHandler) Get(c *gin.Context) {
	detail, err := h.service.Get(c.Request.Context(), tenant.OrgID(c), c.Param("id"))
	if err != nil {
		h.writeError(c, err, "Unable to fetch role")
		return
	}
	response.OK(c, "Role fetched successfully", detail)
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	detail, err := h.service.Create(c.Request.Context(), tenant.OrgID(c), req)
	if err != nil {
		h.writeError(c, err, "Unable to create role")
		return
	}
	response.Created(c, "Role created successfully", detail)
}

func (h *RoleHandler) Update(c *gin.Context) {
	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	detail, err := h.service.Update(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		h.writeError(c, err, "Unable to update role")
		return
	}
	response.OK(c, "Role updated successfully", detail)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), tenant.OrgID(c), c.Param("id")); err != nil {
		h.writeError(c, err, "Unable to delete role")
		return
	}
	response.OK(c, "Role deleted successfully", nil)
}

func (h *RoleHandler) SetPermissions(c *gin.Context) {
	var req dto.SetPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	detail, err := h.service.SetPermissions(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		h.writeError(c, err, "Unable to update permissions")
		return
	}
	response.OK(c, "Permissions updated successfully", detail)
}

func (h *RoleHandler) SetModuleAccess(c *gin.Context) {
	var req dto.SetModuleAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	detail, err := h.service.SetModuleAccess(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		h.writeError(c, err, "Unable to update module access")
		return
	}
	response.OK(c, "Module access updated successfully", detail)
}

func (h *RoleHandler) SetFieldAccess(c *gin.Context) {
	var req dto.SetFieldAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", nil)
		return
	}
	if err := validation.ValidateStruct(&req); err != nil {
		response.BadRequest(c, "Validation failed", validation.FormatErrors(err))
		return
	}

	detail, err := h.service.SetFieldAccess(c.Request.Context(), tenant.OrgID(c), c.Param("id"), req)
	if err != nil {
		h.writeError(c, err, "Unable to update field access")
		return
	}
	response.OK(c, "Field access updated successfully", detail)
}

// Me returns the caller's effective role, permissions and ACL.
func (h *RoleHandler) Me(c *gin.Context) {
	me, err := h.service.Me(
		c.Request.Context(),
		tenant.RoleID(c),
		tenant.RoleSlug(c),
		rbac.Permissions(c),
	)
	if err != nil {
		response.InternalServerError(c, "Unable to fetch access context")
		return
	}
	response.OK(c, "Access context fetched successfully", me)
}

func (h *RoleHandler) writeError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, "Role not found")
	case errors.Is(err, service.ErrSlugTaken):
		response.BadRequest(c, "Role slug already exists", nil)
	case errors.Is(err, service.ErrInvalidSlug):
		response.BadRequest(c, err.Error(), nil)
	case errors.Is(err, service.ErrSystemRole):
		response.BadRequest(c, "System roles cannot be deleted", nil)
	case errors.Is(err, service.ErrHasMembers):
		response.BadRequest(c, "Reassign members before deleting this role", nil)
	default:
		response.InternalServerError(c, fallback)
	}
}
