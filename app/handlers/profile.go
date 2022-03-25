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
const patternDate = "[0-9]+"
const patternId = "^[0-9]+$"

func checkData(profile *models.Profile) bool {
	var check bool
	var err error
	if len(profile.Birthday) > 0 {
		check, err = regexp.MatchString(patternDate, profile.Birthday)
		if !check || err != nil {
			return false
		}
	}
	if profile.FirstName != "" {
		check, err = regexp.MatchString(patternStr, profile.FirstName)
		if !check || err != nil {
			return false
		}
	}
	if profile.LastName != "" {
		check, err = regexp.MatchString(patternStr, profile.LastName)
		if !check || err != nil {
			return false
		}
	}
	if profile.Gender != "" {
		check, err = regexp.MatchString(patternStr, profile.Gender)
		if !check || err != nil {
			return false
		}
	}
	if profile.City != "" {
		check, err = regexp.MatchString(patternStr, profile.City)
		if !check || err != nil {
			return false
		}
	}
	if profile.AboutUser != "" {
		check, err = regexp.MatchString(patternStr, profile.AboutUser)
		if !check || err != nil {
			return false
		}
	}
	for _, value := range profile.Interests {
		check, _ = regexp.MatchString(patternStr, value)
		if !check {
			return false
		}
	}
	return true
}

func getIdFromUrl(r *http.Request) (int, error) {
	idFromUrl := mux.Vars(r)["id"]
	checkIdFromUrl, _ := regexp.MatchString(patternId, idFromUrl)
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
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	id, err := getIdFromUrl(r)
	if err != nil {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}
	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		return
	}
	profile, err := handler.ProfileUseCase.Get(cookieId, id)
	if err != nil {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}

	response, err := json.Marshal(profile)
	if err != nil {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}
	w.Write(response)
}

func (handler *ProfileHandler) GetShortProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	id, err := getIdFromUrl(r)
	if err != nil {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}
	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		return
	}

	profile, err := handler.ProfileUseCase.GetShort(cookieId, id)
	if err != nil {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}
	response, err := json.Marshal(profile)
	if err != nil {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}
	w.Write(response)
}

func (handler *ProfileHandler) ChangeProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	
	msg, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}
	model := &models.Profile{}

	err = json.Unmarshal(msg, model)
	if err != nil || !checkData(model) {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}
	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		return
	}

	err = handler.ProfileUseCase.Change(cookieId, *model)
	if err != nil {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}
}

func (handler *ProfileHandler) GetCandidateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		return
	}

	profile, err := handler.ProfileUseCase.GetCandidates(cookieId)
	if err != nil {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}
	response, err := json.Marshal(profile)
	if err != nil {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}
	w.Write(response)
}
