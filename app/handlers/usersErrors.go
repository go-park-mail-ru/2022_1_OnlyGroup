package handlers

import (
	"net/http"
)

var (
	ErrAuthRequired           = appError{"authentication required", http.StatusUnauthorized}
	ErrAuthEmailUsed          = appError{"email already registered", http.StatusConflict}
	ErrAuthWrongPassword      = appError{"wrong password", http.StatusUnauthorized}
	ErrAuthUserNotFound       = appError{"user not registered", http.StatusUnauthorized}
	ErrAuthSessionNotFound    = appError{"session not found", http.StatusUnauthorized}
	ErrAuthValidationEmail    = appError{"email validation failed", http.StatusPreconditionFailed}
	ErrAuthValidationPassword = appError{"password validation failed", http.StatusPreconditionFailed}
)
