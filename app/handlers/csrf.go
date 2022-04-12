package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/pkg/csrf"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type CSRFHandler struct {
	JwtToken csrf.CsrfGenerator
}

func CreateCSRFHandler(jwt csrf.CsrfGenerator) *CSRFHandler {
	return &CSRFHandler{jwt}
}

func (impl *CSRFHandler) PostCSRF(w http.ResponseWriter, r *http.Request) {
	msg, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	model := &models.CSRF{}
	err = json.Unmarshal(msg, model)
	if err != nil {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.Code)
		return
	}
	ctx := r.Context()
	cookie, err := r.Cookie(authCookie)
	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, ErrAuthRequired.String(), ErrAuthRequired.Code)
		return
	}
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		appErr := ErrBaseApp.LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	token, err := impl.JwtToken.Create(cookie.Value, cookieId, model.URL)
	if err != nil {
		appErr := ErrBaseApp.LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	w.Header().Add("X-CSRF-TOKEN", token)
	w.WriteHeader(http.StatusOK)
}
