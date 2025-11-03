package dto

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    MetaData    `json:"meta"`
}

type MetaData struct {
	Page        int   `json:"page"`
	Limit       int   `json:"limit"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
	HasNext     bool  `json:"has_next"`
	HasPrevious bool  `json:"has_previous"`
}

func NewSuccessResponse(message string, data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(message, error string) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Message: message,
		Error:   error,
	}
}
