package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/usecases"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

const patternStr = "^[a-zA-Z]+$"
const patternDate = "\\d{2}.\\d{2}.\\d{4}"
const patternInt = "^[0-9]+$"
const textSize = 256
const nameSize = 32

func checkProfileData(profile *models.Profile) bool {
	var check bool
	var err error
	if len(profile.Birthday) > 0 {
		check, err = regexp.MatchString(patternDate, profile.Birthday)
		if !check || err != nil {
			return false
		}
	}
	if profile.FirstName != "" || len(profile.FirstName) > nameSize {
		check, err = regexp.MatchString(patternStr, profile.FirstName)
		if !check || err != nil {
			return false
		}
	}
	if profile.LastName != "" || len(profile.LastName) > nameSize {
		check, err = regexp.MatchString(patternStr, profile.LastName)
		if !check || err != nil {
			return false
		}
	}
	if profile.Gender > 1 || profile.Gender < 0 {
		return false
	}
	if profile.City != "" {
		check, err = regexp.MatchString(patternStr, profile.City)
		if !check || err != nil {
			return false
		}
	}
	if len(profile.AboutUser) > textSize {
		return false
	}
	for _, value := range profile.Interests {
		if len(value) > textSize {
			return false
		}
	}
	if profile.Height > 300 || profile.Height < 0 {
		return false
	}
	return true
}

func getIdFromUrl(r *http.Request) (int, error) {
	idFromUrl := mux.Vars(r)["id"]
	checkIdFromUrl, _ := regexp.MatchString(patternInt, idFromUrl)
	if !checkIdFromUrl {
		return 0, ErrBadUserID
	}
	id, err := strconv.Atoi(idFromUrl)
	if errors.Is(err, strconv.ErrSyntax) {
		return 0, ErrBadUserID
	}
	return id, nil
}

type ProfileHandler struct {
	ProfileUseCase usecases.ProfileUseCases
}

func CreateProfileHandler(useCase usecases.ProfileUseCases) *ProfileHandler {
	return &ProfileHandler{ProfileUseCase: useCase}
}

func (handler *ProfileHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromUrl(r)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
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
	profile, err := handler.ProfileUseCase.Get(cookieId, id)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	response, err := json.Marshal(profile)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	w.Write(response)
}

func (handler *ProfileHandler) GetShortProfileHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromUrl(r)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
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
	model, err := handler.ProfileUseCase.GetShort(cookieId, id)
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
}

func (handler *ProfileHandler) ChangeProfileHandler(w http.ResponseWriter, r *http.Request) {
	msg, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		appErr := ErrBaseApp.LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	model := &models.Profile{}

	err = json.Unmarshal(msg, model)

	if err != nil || !checkProfileData(model) {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.Code)
		return
	}
	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	err = handler.ProfileUseCase.Change(cookieId, *model)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
}

func (handler *ProfileHandler) GetCandidateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		appErr := ErrBaseApp.LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	vectorCandidates, err := handler.ProfileUseCase.GetCandidates(cookieId)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	response, err := json.Marshal(vectorCandidates)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	w.Write(response)
}
