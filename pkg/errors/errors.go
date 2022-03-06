package errors

import (
	"errors"
	"net/http"
)

var (
	ErrAuthEmailUsed     = errors.New("email already registered")
	ErrAuthWrongPassword = errors.New("wrong password")
	ErrAuthUserNotFound  = errors.New("user not registered")

	ErrAuthSessionNotFound = errors.New("session not found")

	ErrProfileNotFound = errors.New("profile not found")
	ErrProfileNotFiled = errors.New("profile not filed")
	ErrMockIsEmpty     = errors.New("profile not filed")
)

var errorsCodes = map[error]int{
	ErrAuthSessionNotFound: http.StatusUnauthorized,
	ErrAuthEmailUsed:       http.StatusConflict,
	ErrAuthWrongPassword:   http.StatusUnauthorized,
	ErrAuthUserNotFound:    http.StatusUnauthorized,
}

func ErrorToHTTPCode(err error) int {
	return errorsCodes[err]
}
