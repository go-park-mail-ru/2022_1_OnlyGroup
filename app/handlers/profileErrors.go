package handlers

import (
	"net/http"
)

var (
	ErrBadUserID       = appError{"bad user id from url", http.StatusBadRequest}
	ErrProfileNotFound = appError{"profile not found", http.StatusNotFound}
	ErrProfileNotFiled = appError{"profile not filed", http.StatusNotAcceptable}
	ErrMockIsEmpty     = appError{"profile not filed", http.StatusNotAcceptable}
)
