package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrAuthEmailUsed     = errors.New("email already registered")
	ErrAuthWrongPassword = errors.New("wrong password")
	ErrAuthUserNotFound  = errors.New("user not registered")

	ErrAuthSessionNotFound = errors.New("session not found")
)

var errorsCodes = map[error]int{
	ErrAuthSessionNotFound: http.StatusUnauthorized,
	ErrAuthEmailUsed:       http.StatusConflict,
	ErrAuthWrongPassword:   http.StatusUnauthorized,
	ErrAuthUserNotFound:    http.StatusUnauthorized,
}

func ErrorToHTTPCode(err error) int {
	code, has := errorsCodes[err]
	if !has {
		code = http.StatusInternalServerError
	}
	return code
}

func WrapError2Json(input string) string {
	errorStruct := struct {
		ErrorMsg string
	}{input}
	wrappedError, _ := json.Marshal(errorStruct)
	return string(wrappedError)
}
