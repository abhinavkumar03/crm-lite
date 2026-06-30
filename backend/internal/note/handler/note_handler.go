package handler

import (
	"net/http"

	"github.com/abhinavkumar03/crm-lite/backend/internal/note/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/note/service"
	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	service *service.Service
}

func New(
	service *service.Service,
) *NoteHandler {
	return &NoteHandler{
		service: service,
	}
}

func (h *NoteHandler) CreateLeadNote(
	c *gin.Context,
) {

	var req dto.CreateNoteRequest

	if err := c.ShouldBindJSON(
		&req,
	); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)

		return

	}

	req.EntityType = "LEAD"

	req.EntityID = c.Param("leadId")

	userID := c.GetString("userID")

	err := h.service.Create(
		c.Request.Context(),
		userID,
		req,
	)

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

	c.JSON(
		http.StatusCreated,
		gin.H{
			"success": true,
			"message": "Note created successfully",
		},
	)

}

func (h *NoteHandler) ListLeadNotes(
	c *gin.Context,
) {

	userID := c.GetString(
		"userID",
	)

	notes, err := h.service.List(

		c.Request.Context(),

		userID,

		"LEAD",

		c.Param("leadId"),
	)

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

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data":    notes,
		},
	)

}
func (h *NoteHandler) UpdateNote(
	c *gin.Context,
) {

	var req dto.UpdateNoteRequest

	if err := c.ShouldBindJSON(
		&req,
	); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)

		return

	}

	userID := c.GetString(
		"userID",
	)

	err := h.service.Update(

		c.Request.Context(),

		userID,

		c.Param("noteId"),

		req,
	)

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

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "Note updated successfully",
		},
	)

}

func (h *NoteHandler) DeleteNote(
	c *gin.Context,
) {

	userID := c.GetString(
		"userID",
	)

	err := h.service.Delete(

		c.Request.Context(),

		userID,

		c.Param("noteId"),
	)

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

	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"message": "Note deleted successfully",
		},
	)

}
