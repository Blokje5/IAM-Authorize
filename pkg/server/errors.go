package server

import "encoding/json"

const (
	ErrConflict = "/errors/conflict"
	ErrInternalServer = "errors/internal-server"
)

// ErrorResponse wraps the error neatly in a RFC7807 problem
type ErrorResponse struct {
	ErrorType string `json:"type"`
	Title string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

func (e *ErrorResponse) Error() string {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return string(data)
}

// NewErrorResponse response returns a new ErrorResponse
func NewErrorResponse(errorType, title, detail string) *ErrorResponse {
	// TODO add path resolution based on server host name
	errorPath := "http://localhost:8080" + errorType
	return &ErrorResponse{
		ErrorType: errorPath,
		Title: title,
		Detail: detail,
	}
}

// NewConflictError returns a new HTTP 409 Conflict error with the given message
func NewConflictError(title, detail string) *ErrorResponse {
	return NewErrorResponse(ErrConflict, title, detail)
}

// NewInternalServerError returns a new HTTP 500 Internal Server error with the given message
func NewInternalServerError(title, detail string) *ErrorResponse {
	return NewErrorResponse(ErrInternalServer, title, detail)
}