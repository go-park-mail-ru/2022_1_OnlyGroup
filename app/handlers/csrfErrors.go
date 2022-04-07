package handlers

var (
	ErrBadCSRF    = AppError{Msg: "Bad token", Code: 419}
	ErrBadSession = AppError{Msg: "Different session data"}
)
