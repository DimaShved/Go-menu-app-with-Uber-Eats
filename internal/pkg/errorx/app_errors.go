package errorx

import (
	"errors"
	"fmt"
)

type AppError struct {
	Code    int
	Message string
	Details string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		if e.Details != "" {
			return fmt.Sprintf("%s: %s (%s)", e.Message, e.Details, e.Err)
		}
		return fmt.Sprintf("%s: %s", e.Message, e.Err)
	}
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Err:     err,
	}
}

func (e *AppError) WithDetails(details string) *AppError {
	newErr := *e
	newErr.Details = details
	return &newErr
}

func NewAppError(code int, msg string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
	}
}

var (
	ErrConfigLoad               = NewAppError(100, "Error loading configuration")
	ErrDatabase                 = NewAppError(101, "Database error") // General DB error
	ErrValidation               = NewAppError(102, "Validation failed")
	ErrServerStart              = NewAppError(103, "Error starting server")
	ErrNotImplementIdentifiable = NewAppError(104, "Entity does not implement Identifiable interface")
	ErrDatabaseQuery            = NewAppError(105, "Error executing database query")
	ErrRecordNotFound           = NewAppError(106, "Record not found") // More specific DB error
)

func IsAppError(err error, target *AppError) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == target.Code && appErr.Message == target.Message
	}
	return false
}
