package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CSRFHandler struct {
	JwtToken      JwtToken
	TokenLifeTime int
}

func CreateCSRFHandler(jwt JwtToken, lifeTime int) *CSRFHandler {
	return &CSRFHandler{jwt, lifeTime}
}

func (impl *CSRFHandler) GetCSRF(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println(time.Now().Add(time.Hour * time.Duration(impl.TokenLifeTime)))
	token, err := impl.JwtToken.Create(cookie.Value, cookieId, model.URL, time.Now().Add(time.Hour*time.Duration(impl.TokenLifeTime)).Unix())
	if err != nil {
		appErr := ErrBaseApp.LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	w.Header().Add("X-CSRF-TOKEN", token)
}
