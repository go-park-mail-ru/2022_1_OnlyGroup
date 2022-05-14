package http

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/usecases"
	"encoding/json"
	"gopkg.in/validator.v2"
	"io"
	"net/http"
)

type FiltersHandler struct {
	profileUseCase usecases.ProfileUseCases
}

func CreateFiltersHandler(useCase usecases.ProfileUseCases) *FiltersHandler {
	return &FiltersHandler{useCase}
}

func (handler *FiltersHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		appErr := ErrBaseApp.LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	filters, err := handler.profileUseCase.GetFilters(cookieId)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	response, err := json.Marshal(filters)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	w.Write(response)
}

func (handler *FiltersHandler) Put(w http.ResponseWriter, r *http.Request) {
	msg, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		appErr := ErrBaseApp.LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	model := models.Filters{}

	err = json.Unmarshal(msg, &model)

	if err != nil {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.Code)
		return
	}
	err = validator.Validate(model)
	if err == ErrBaseApp {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	if err != nil {
		http.Error(w, ErrValidateProfile.String(), ErrValidateProfile.Code)
		return
	}

	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		appErr := ErrBaseApp.LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	err = handler.profileUseCase.ChangeFilters(cookieId, model)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

}
