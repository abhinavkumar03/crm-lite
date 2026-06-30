package handler

import (
	"net/http"

	"github.com/abhinavkumar03/crm-lite/backend/internal/media/service"
	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	service *service.Service
}

func New(
	service *service.Service,
) *MediaHandler {

	return &MediaHandler{

		service: service,
	}
}

func (h *MediaHandler) Upload(
	c *gin.Context,
) {

	fileHeader, err := c.FormFile("file")

	if err != nil {

		c.JSON(

			http.StatusBadRequest,

			gin.H{

				"success": false,

				"message": "file is required",
			},
		)

		return
	}

	file, err := fileHeader.Open()

	if err != nil {

		c.JSON(

			http.StatusBadRequest,

			gin.H{

				"success": false,

				"message": err.Error(),
			},
		)

		return
	}

	defer file.Close()

	result, err := h.service.Upload(

		c.Request.Context(),

		file,
	)

	if err != nil {

		c.JSON(

			http.StatusInternalServerError,

			gin.H{

				"success": false,

				"message": err.Error(),
			},
		)

		return
	}

	c.JSON(

		http.StatusOK,

		gin.H{

			"success": true,

			"data": result,
		},
	)
}
