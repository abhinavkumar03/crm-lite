package handler

import (
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/dto"

	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/validation"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func New(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {

	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.BadRequest(
			c,
			"Invalid request body",
			nil,
		)

		return
	}

	if err := validation.ValidateStruct(&req); err != nil {

		response.BadRequest(
			c,
			"Validation failed",
			validation.FormatErrors(err),
		)

		return
	}

	user, err := h.service.Register(
		c.Request.Context(),
		req,
	)

	if err != nil {

		response.BadRequest(
			c,
			err.Error(),
			nil,
		)

		return
	}

	response.Created(
		c,
		"User registered successfully",
		user,
	)
}

func (h *AuthHandler) Login(
	c *gin.Context,
) {

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.BadRequest(
			c,
			"Invalid request body",
			nil,
		)

		return
	}

	if err := validation.ValidateStruct(&req); err != nil {

		response.BadRequest(
			c,
			"Validation failed",
			validation.FormatErrors(err),
		)

		return
	}

	loginResponse, err := h.service.Login(
		c.Request.Context(),
		req,
	)

	if err != nil {

		response.BadRequest(
			c,
			err.Error(),
			nil,
		)

		return
	}

	response.OK(
		c,
		"Login successful",
		loginResponse,
	)
}

func (h *AuthHandler) Profile(c *gin.Context) {

	userID := c.GetString("userID")

	user, err := h.service.GetProfile(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		response.NotFound(
			c,
			"User not found",
		)
		return
	}

	response.OK(
		c,
		"Profile fetched successfully",
		user,
	)
}
