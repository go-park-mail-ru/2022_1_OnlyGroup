package handlers

import (
	"2022_1_OnlyGroup_back/app/usecases"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const requestIdContextKey = "requestId"
const userIdContextKey = "userId"

//const authCookie = "session"

type Middlewares interface {
	AccessLogMiddleware(next http.Handler) http.Handler
	PanicMiddleware(next http.Handler) http.Handler
	CheckAuthMiddleware(next http.Handler) http.Handler
}

type MiddlewaresImpl struct {
	AuthUseCase usecases.AuthUseCases
}

func (impl MiddlewaresImpl) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := r.Header.Get("X-Request-ID")
		if len(reqId) == 0 {
			reqId = uuid.NewString()
		}
		newContext := context.WithValue(r.Context(), requestIdContextKey, reqId)
		rNew := r.Clone(newContext)

		startTime := time.Now()
		next.ServeHTTP(w, rNew)
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
				reqId := r.Context().Value(requestIdContextKey)
				if reqId == nil {
					reqId = "[no RequestId]"
				}

				logrus.WithFields(logrus.Fields{
					"mode": "panic_log",
					"time": time.Now().String(),
				}).Errorf("[%s] Panic! %v %s %s %s", reqId, err, r.RemoteAddr, r.Method, r.RequestURI)
				appErr := appError{code: http.StatusInternalServerError, Msg: fmt.Sprint(err)}
				http.Error(w, appErr.String(), appErr.Code())
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (imlp MiddlewaresImpl) CheckAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(authCookie)
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, ErrAuthRequired.String(), ErrAuthRequired.Code())
			return
		}
		userIdModel, err := imlp.AuthUseCase.UserAuth(cookie.Value)
		if errors.Is(err, ErrAuthSessionNotFound) {
			http.Error(w, ErrAuthSessionNotFound.String(), ErrAuthSessionNotFound.Code())
			return
		}

		newContext := context.WithValue(r.Context(), userIdContextKey, userIdModel.ID)
		rNew := r.Clone(newContext)
		next.ServeHTTP(w, rNew)
	})
}
