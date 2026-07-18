package handler

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/parser"
	"github.com/abhinavkumar03/crm-lite/backend/internal/importer/service"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/response"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tenant"
)

const (
	paramModuleID = "id"
	paramImportID = "importId"

	// maxUploadBytes caps the accepted upload to keep a single import bounded.
	maxUploadBytes = 10 << 20 // 10 MiB
)

type ImportHandler struct {
	service *service.Service
}

func New(service *service.Service) *ImportHandler {
	return &ImportHandler{service: service}
}

func userID(c *gin.Context) string { return c.GetString("userID") }

// Analyze parses an uploaded file (without persisting it) and returns the
// columns, a preview sample and a suggested column-to-field mapping.
func (h *ImportHandler) Analyze(c *gin.Context) {
	filename, data, ok := readUpload(c)
	if !ok {
		return
	}

	result, err := h.service.Analyze(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), filename, data)
	if err != nil {
		h.writeError(c, err, "Unable to analyze file")
		return
	}
	response.OK(c, "File analyzed successfully", result)
}

// Create stages the uploaded file with the chosen mapping and enqueues it for
// asynchronous processing.
func (h *ImportHandler) Create(c *gin.Context) {
	filename, data, ok := readUpload(c)
	if !ok {
		return
	}

	mapping, err := decodeStringMap(c.PostForm("mapping"))
	if err != nil {
		response.BadRequest(c, "Invalid mapping: expected a JSON object of column -> field", nil)
		return
	}

	options, err := decodeAnyMap(c.PostForm("options"))
	if err != nil {
		response.BadRequest(c, "Invalid options: expected a JSON object", nil)
		return
	}

	job, err := h.service.Create(
		c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), userID(c), filename, data, mapping, options,
	)
	if err != nil {
		h.writeError(c, err, "Unable to start import")
		return
	}
	response.Created(c, "Import started successfully", job)
}

func (h *ImportHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	q := dto.ListQuery{
		Page:     page,
		PageSize: pageSize,
		Status:   c.Query("status"),
	}

	result, err := h.service.List(c.Request.Context(), tenant.OrgID(c), c.Param(paramModuleID), q)
	if err != nil {
		h.writeError(c, err, "Unable to fetch imports")
		return
	}
	response.OK(c, "Imports fetched successfully", result)
}

func (h *ImportHandler) Get(c *gin.Context) {
	job, err := h.service.Get(c.Request.Context(), tenant.OrgID(c), c.Param(paramImportID))
	if err != nil {
		h.writeError(c, err, "Unable to fetch import")
		return
	}
	response.OK(c, "Import fetched successfully", job)
}

// readUpload extracts the multipart "file" part, enforcing the size cap. It
// writes the error response itself and returns ok=false on failure.
func readUpload(c *gin.Context) (string, []byte, bool) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "A file upload is required (multipart field 'file')", nil)
		return "", nil, false
	}
	if fileHeader.Size > maxUploadBytes {
		response.BadRequest(c, "File is too large (max 10 MiB)", nil)
		return "", nil, false
	}

	f, err := fileHeader.Open()
	if err != nil {
		response.InternalServerError(c, "Unable to read uploaded file")
		return "", nil, false
	}
	defer f.Close()

	data, err := io.ReadAll(io.LimitReader(f, maxUploadBytes+1))
	if err != nil {
		response.InternalServerError(c, "Unable to read uploaded file")
		return "", nil, false
	}
	return fileHeader.Filename, data, true
}

func decodeStringMap(raw string) (map[string]string, error) {
	if raw == "" {
		return map[string]string{}, nil
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, err
	}
	return m, nil
}

func decodeAnyMap(raw string) (map[string]any, error) {
	if raw == "" {
		return map[string]any{}, nil
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (h *ImportHandler) writeError(c *gin.Context, err error, fallback string) {
	switch {
	case errors.Is(err, service.ErrModuleNotFound):
		response.NotFound(c, "Module not found")
	case errors.Is(err, service.ErrNotDynamic):
		response.BadRequest(c, "This module does not support imports", nil)
	case errors.Is(err, service.ErrNotFound):
		response.NotFound(c, "Import job not found")
	case errors.Is(err, service.ErrNoMapping):
		response.BadRequest(c, "Map at least one column to a field before importing", nil)
	case errors.Is(err, service.ErrTooManyRows):
		response.BadRequest(c, "File exceeds the maximum of 5000 rows", nil)
	case errors.Is(err, parser.ErrUnsupported):
		response.BadRequest(c, "Unsupported file type (use .csv or .xlsx)", nil)
	case errors.Is(err, parser.ErrEmptyFile):
		response.BadRequest(c, "The file has no data rows", nil)
	case errors.Is(err, parser.ErrNoHeaders):
		response.BadRequest(c, "The file has no header row", nil)
	default:
		response.InternalServerError(c, fallback)
	}
}
