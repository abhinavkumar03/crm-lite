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

func (h *NoteHandler) createEntityNote(
	c *gin.Context,
	entityType string,
	entityID string,
) {

	var req dto.CreateNoteRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"success": false,
				"message": err.Error(),
			},
		)

		return
	}

	req.EntityType = entityType
	req.EntityID = entityID

	userID := c.GetString("userID")

	if err := h.service.Create(
		c.Request.Context(),
		userID,
		req,
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

	c.JSON(
		http.StatusCreated,
		gin.H{
			"success": true,
			"message": "Note created successfully",
		},
	)
}

func (h *NoteHandler) listEntityNotes(
	c *gin.Context,
	entityType string,
	entityID string,
) {

	userID := c.GetString("userID")

	notes, err := h.service.List(
		c.Request.Context(),
		userID,
		entityType,
		entityID,
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

func (h *NoteHandler) CreateLeadNote(
	c *gin.Context,
) {
	h.createEntityNote(
		c,
		"LEAD",
		c.Param("leadId"),
	)
}

func (h *NoteHandler) ListLeadNotes(
	c *gin.Context,
) {
	h.listEntityNotes(
		c,
		"LEAD",
		c.Param("leadId"),
	)
}

func (h *NoteHandler) CreateTaskNote(
	c *gin.Context,
) {
	h.createEntityNote(
		c,
		"TASK",
		c.Param("taskId"),
	)
}

func (h *NoteHandler) ListTaskNotes(
	c *gin.Context,
) {
	h.listEntityNotes(
		c,
		"TASK",
		c.Param("taskId"),
	)
}

func (h *NoteHandler) CreateContactNote(
	c *gin.Context,
) {
	h.createEntityNote(
		c,
		"CONTACT",
		c.Param("contactId"),
	)
}

func (h *NoteHandler) ListContactNotes(
	c *gin.Context,
) {
	h.listEntityNotes(
		c,
		"CONTACT",
		c.Param("contactId"),
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
