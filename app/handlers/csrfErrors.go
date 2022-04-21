package handlers

import "net/http"

var (
	ErrBadCSRF = AppError{Msg: "Bad token", Code: http.StatusUnauthorized}
)
