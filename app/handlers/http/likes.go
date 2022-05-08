package http

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/usecases"
	"encoding/json"

	"gopkg.in/validator.v2"
	"io"
	"net/http"
)

type LikesHandler struct {
	profileUseCase usecases.ProfileUseCases
}

func CreateLikesHandler(likesUseCase usecases.ProfileUseCases) *LikesHandler {
	return &LikesHandler{profileUseCase: likesUseCase}
}

func (handler *LikesHandler) Set(w http.ResponseWriter, r *http.Request) {
	msg, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		appErr := ErrBaseApp.LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		appErr := ErrBaseApp.LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	model := &models.Likes{}
	err = json.Unmarshal(msg, model)

	if err != nil {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.Code)
		return
	}
	if err = validator.Validate(model); err != nil {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.Code)
	}
	err = handler.profileUseCase.SetAction(cookieId, *model)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
}

func (handler *LikesHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		appErr := ErrBaseApp.LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	model, err := handler.profileUseCase.GetMatched(cookieId)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	response, err := json.Marshal(model)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	w.Write(response)
	return
}
