package handlers

import (
	"2022_1_OnlyGroup_back/app/usecases"
	"2022_1_OnlyGroup_back/pkg/csrf"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const requestIdContextKey = "requestId"
const userIdContextKey = "userId"

type Middlewares interface {
	AccessLogMiddleware(next http.Handler) http.Handler
	PanicMiddleware(next http.Handler) http.Handler
	CheckAuthMiddleware(next http.Handler) http.Handler
	CorsMiddleware(next http.Handler) http.Handler
	CSRFMiddleware(next http.Handler) http.Handler
}

type MiddlewaresImpl struct {
	AuthUseCase usecases.AuthUseCases
	JwtConf     csrf.CsrfGenerator
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
				http.Error(w, ErrBaseApp.String(), ErrBaseApp.Code)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (imlp MiddlewaresImpl) CheckAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(authCookie)
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, ErrAuthRequired.String(), ErrAuthRequired.Code)
			return
		}
		userIdModel, err := imlp.AuthUseCase.UserAuth(cookie.Value)
		if err != nil {
			appErr := AppErrorFromError(err)
			http.Error(w, appErr.String(), appErr.Code)
			return
		}

		newContext := context.WithValue(r.Context(), userIdContextKey, userIdModel.ID)
		rNew := r.Clone(newContext)
		next.ServeHTTP(w, rNew)
	})
}

func (imlp MiddlewaresImpl) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(w, r)
	})
}

func (imlp MiddlewaresImpl) CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrfToken := r.Header.Get("X-CSRF-TOKEN")
		if len(csrfToken) == 0 {
			http.Error(w, ErrBadCSRF.String(), ErrBadCSRF.Code)
			return
		}
		ctx := r.Context()
		cookieId, ok := ctx.Value(userIdContextKey).(int)
		if !ok {
			appErr := ErrBaseApp.LogServerError(r.Context().Value(requestIdContextKey))
			http.Error(w, appErr.String(), appErr.Code)
			return
		}
		cookie, err := r.Cookie(authCookie)
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, ErrAuthRequired.String(), ErrAuthRequired.Code)
			return
		}
		err = imlp.JwtConf.Check(cookie.Value, cookieId, r.URL.String(), csrfToken)
		if err != nil {
			appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
			http.Error(w, appErr.String(), appErr.Code)
			return
		}
		next.ServeHTTP(w, r)
	})
}
