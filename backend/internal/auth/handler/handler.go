package handler

import "github.com/gin-gonic/gin"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(
	router *gin.RouterGroup,
) {

	// Routes will be added in Phase 5.
}
