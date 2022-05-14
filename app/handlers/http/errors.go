package http

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type AppError struct {
	Msg         string
	Code        int
	Base        error  `json:"-"`
	Description string `json:"-"`
}

var ErrBaseApp = AppError{"internal server error", http.StatusInternalServerError, nil, ""}
var ErrBadRequest = AppError{"bad request", http.StatusBadRequest, nil, ""}
var ErrServiceUnavailable = AppError{"profile service Unavailable,", http.StatusServiceUnavailable, nil, ""}

func AppErrorFromError(inputError error) AppError {
	appErr, ok := inputError.(AppError)
	if !ok {
		return ErrBaseApp.Wrap(inputError, "")
	}
	return appErr
}

func (err AppError) IsInternalError() bool {
	if err.Code/100 == 5 {
		return true
	}
	return false
}

func (err AppError) Wrap(baseErr error, desc string) AppError {
	err.Base = baseErr
	err.Description = desc
	return err
}

func (err AppError) Is(target error) bool {
	targetAppErr, ok := target.(AppError)
	if !ok {
		return target == err.Base
	}
	return targetAppErr.Code == err.Code && targetAppErr.Msg == err.Msg
}

func (err AppError) LogServerError(reqId interface{}) AppError {
	if err.IsInternalError() {
		logrus.WithFields(logrus.Fields{
			"mode": "server_error_log",
		}).Errorf("[%s] %d %s %v", reqId, err.Code, err.Description, err.Base)
	}
	return err
}

func (err AppError) Error() string {
	return err.Msg
}

func (err AppError) String() string {
	errBuffer, er := json.Marshal(err)
	if er != nil {
		panic(er)
	}
	return string(errBuffer)
}

func (err AppError) WrapGRPC() error {
	switch err.Code {
	case ErrBaseApp.Code:
		return status.Errorf(codes.Internal, err.Msg)
	case ErrBadRequest.Code:
		return status.Errorf(codes.InvalidArgument, err.Msg)
	}
	return nil
}

func AppErrorFromGRPC(inputError error) AppError {
	err, ok := status.FromError(inputError)
	if !ok {
		return ErrBaseApp.Wrap(inputError, "")
	}
	switch err.Code() {
	case codes.InvalidArgument:
		return ErrBadRequest.Wrap(inputError, err.Message())
	case codes.Unavailable:
		return ErrServiceUnavailable.Wrap(inputError, err.Message())
	}
	return ErrBaseApp.Wrap(inputError, err.Message())
}
