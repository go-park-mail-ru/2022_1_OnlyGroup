package handlers

const statusAuthenticationTimeout = 419

var (
	ErrBadCSRF = AppError{Msg: "Bad token", Code: statusAuthenticationTimeout}
)
