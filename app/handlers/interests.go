package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/usecases"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
)

func sanitizeInterests(interest []models.Interest) {
	sanitizer := bluemonday.UGCPolicy()
	for idx, value := range interest {
		interest[idx].Title = sanitizer.Sanitize(value.Title)
	}
}

func getStringFromUrl(r *http.Request) (string, error) {
	str := mux.Vars(r)["str"]
	return str, nil
}

type InterestsHandler struct {
	InterestsUseCase usecases.InterestsUseCase
}

func CreateInterestsHandler(useCase usecases.InterestsUseCase) *InterestsHandler {
	return &InterestsHandler{useCase}
}

func (handler *InterestsHandler) Get(w http.ResponseWriter, r *http.Request) {
	var interests []models.Interest
	interests, err := handler.InterestsUseCase.Get()
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	sanitizeInterests(interests)
	response, err := json.Marshal(interests)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	w.Write(response)
}

func (handler *InterestsHandler) GetDynamic(w http.ResponseWriter, r *http.Request) {
	var interests []models.Interest
	str, err := getStringFromUrl(r)
	interests, err = handler.InterestsUseCase.GetDynamic(str)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	sanitizeInterests(interests)
	response, err := json.Marshal(interests)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	w.Write(response)
}
