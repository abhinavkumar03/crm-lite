package handler

import "github.com/gin-gonic/gin"

type AuthHandler struct {
}

func New() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Register(c *gin.Context) {
}

func (h *AuthHandler) Login(c *gin.Context) {
}

func (h *AuthHandler) Profile(c *gin.Context) {
}
