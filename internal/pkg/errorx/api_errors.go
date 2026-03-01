package errorx

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Code       int    `json:"code"`
	HTTPStatus int    `json:"status"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API Error %d (%d): %s - %s", e.Code, e.HTTPStatus, e.Message, e.Details)
	}
	return fmt.Sprintf("API Error %d (%d): %s", e.Code, e.HTTPStatus, e.Message)
}

func NewAPIError(code, httpStatus int, msg string, details string) *APIError {
	return &APIError{
		Code:       code,
		HTTPStatus: httpStatus,
		Message:    msg,
		Details:    details,
	}
}

func (e *APIError) WithDetails(details string) *APIError {
	return &APIError{
		Code:       e.Code,
		HTTPStatus: e.HTTPStatus,
		Message:    e.Message,
		Details:    details,
	}
}

var (
	ErrNotFound     = NewAPIError(1100, http.StatusNotFound, "Resource not found", "")
	ErrInternal     = NewAPIError(1101, http.StatusInternalServerError, "Internal server error", "")
	ErrInvalidInput = NewAPIError(1102, http.StatusBadRequest, "Invalid request data", "")
)
