package http

import "net/http"

var (
	ErrPhotoNotFound           = AppError{"photo not found", http.StatusNotFound, nil, ""}
	ErrPhotoChangeForbidden    = AppError{"change photo forbidden", http.StatusForbidden, nil, ""}
	ErrUnsupportedType         = AppError{"unsupported photo type", http.StatusUnsupportedMediaType, nil, ""}
	ErrContentLengthMismatched = AppError{"content lenght header mismatched", http.StatusBadRequest, nil, ""}
)
