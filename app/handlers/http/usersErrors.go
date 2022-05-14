package http

import (
	"net/http"
)

var (
	ErrAuthRequired           = AppError{"authentication required", http.StatusUnauthorized, nil, ""}
	ErrAuthEmailUsed          = AppError{"email already registered", http.StatusConflict, nil, ""}
	ErrAuthWrongPassword      = AppError{"wrong password", http.StatusUnauthorized, nil, ""}
	ErrAuthUserNotFound       = AppError{"user not registered", http.StatusUnauthorized, nil, ""}
	ErrAuthSessionNotFound    = AppError{"session not found", http.StatusUnauthorized, nil, ""}
	ErrAuthValidationEmail    = AppError{"email validation failed", http.StatusPreconditionFailed, nil, ""}
	ErrAuthValidationPassword = AppError{"password validation failed", http.StatusPreconditionFailed, nil, ""}
)
