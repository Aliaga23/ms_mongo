package dto

// APIResponse - Respuesta estándar de la API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse - Respuesta de error
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// PaginatedResponse - Respuesta paginada
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    MetaData    `json:"meta"`
}

// MetaData - Metadatos para paginación
type MetaData struct {
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

// NewSuccessResponse crea una respuesta exitosa
func NewSuccessResponse(message string, data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse crea una respuesta de error
func NewErrorResponse(message, error string) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Message: message,
		Error:   error,
	}
}