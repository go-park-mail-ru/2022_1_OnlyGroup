package handlers

import (
	"net/http"
)

var (
	ErrBadID           = AppError{"bad id from url", http.StatusBadRequest, nil, ""}
	ErrProfileNotFound        = AppError{"profile not found", http.StatusNotFound, nil, ""}
	ErrProfileNotFiled        = AppError{"profile not filed", http.StatusNotAcceptable, nil, ""}
	ErrProfileForbiddenChange = AppError{"forbidden change profile", http.StatusForbidden, nil, ""}
	ErrMockIsEmpty            = AppError{"profile not filed", http.StatusNotAcceptable, nil, ""}
	ErrValidateProfile        = AppError{"profile fields notvalid,", http.StatusPreconditionFailed, nil, ""}
)
