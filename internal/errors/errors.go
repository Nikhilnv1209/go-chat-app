package errors

import "fmt"

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"` // HTTP status code
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Predefined Errors
var (
	ErrInvalidCredentials = &AppError{Code: "AUTH_INVALID_CREDENTIALS", Message: "Email or password is incorrect", Status: 401}
	ErrEmailExists        = &AppError{Code: "AUTH_EMAIL_EXISTS", Message: "Email already registered", Status: 409}
	ErrUnauthorized       = &AppError{Code: "AUTH_UNAUTHORIZED", Message: "Authentication required", Status: 401}
	ErrForbidden          = &AppError{Code: "AUTH_FORBIDDEN", Message: "You don't have permission", Status: 403}
	ErrNotFound           = &AppError{Code: "RESOURCE_NOT_FOUND", Message: "User not found", Status: 404}
	ErrValidation         = &AppError{Code: "VALIDATION_ERROR", Message: "Invalid input", Status: 400}
	ErrInternalServer     = &AppError{Code: "INTERNAL_SERVER_ERROR", Message: "An unexpected error occurred", Status: 500}
)
