package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Middlewares interface {
	AccessLogMiddleware(next http.Handler) http.Handler
}

type MiddlewaresImpl struct {
	//AuthUseCase usecases.AuthUseCases
}

func (impl MiddlewaresImpl) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := r.Header.Get("X-Request-ID")
		if len(reqId) == 0 {
			reqId = uuid.NewString()
		}
		a := r.Context()
		a.
			r.WithContext(r.Context().With)

		startTime := time.Now()
		next.ServeHTTP(w, r)
		logrus.WithFields(logrus.Fields{
			"mode": "access_log",
			"time": startTime.String(),
		}).Infof("[%s] %s %s %s %s", reqId, r.RemoteAddr, r.Method, r.RequestURI, time.Since(startTime))
	})
}

func (impl MiddlewaresImpl) PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				fmt.Println("recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
