package handlers

import (
	"encoding/json"
	"net/http"
)

type appError struct {
	ErrorMsg string
	code     int
}

var ErrBaseApp = appError{"internal server error", http.StatusInternalServerError}
var ErrBadRequest = appError{"bad request", http.StatusBadRequest}

func appErrorFromError(inputError error) appError {
	appErr, ok := inputError.(appError)
	if !ok {
		return ErrBaseApp
	}
	return appErr
}

func (err appError) Error() string {
	return err.ErrorMsg
}

func (err appError) String() string {
	errBuffer, er := json.Marshal(err)
	if er != nil {
		panic(er)
	}
	return string(errBuffer)
}

func (err appError) Code() int {
	return err.code
}
