package response

// APIResponse is the standard response returned by every endpoint.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// ValidationError represents a validation error for a single field.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
