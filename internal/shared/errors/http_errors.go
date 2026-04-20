package errors

// ErrorResponse is the standard JSON body for all error responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{Error: msg}
}
