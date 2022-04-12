package handlers

var (
	ErrBadCSRF = AppError{Msg: "Bad token", Code: 419}
)
