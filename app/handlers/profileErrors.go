package handlers

import (
	"net/http"
)

var (
	ErrBadUserID       = AppError{"bad user id from url", http.StatusBadRequest, nil, ""}
	ErrProfileNotFound = AppError{"profile not found", http.StatusNotFound, nil, ""}
	ErrProfileNotFiled = AppError{"profile not filed", http.StatusNotAcceptable, nil, ""}
	ErrMockIsEmpty     = AppError{"profile not filed", http.StatusNotAcceptable, nil, ""}
)
