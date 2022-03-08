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

func checkData(profile *models.Profile) bool {
	var check bool
	if len(profile.Birthday) > 0 {
		check, _ = regexp.MatchString(patternDate, profile.Birthday)
		if !check {
			return false
		}
	}

	if profile.FirstName != "" {
		check, _ = regexp.MatchString(patternStr, profile.FirstName)
		if !check {
			return false
		}
	}

	if profile.LastName != "" {
		check, _ = regexp.MatchString(patternStr, profile.LastName)
		if !check {
			return false
		}
	}
	if profile.Gender != "" {
		check, _ = regexp.MatchString(patternStr, profile.Gender)
		if !check {
			return false
		}
	}
	if profile.City != "" {
		check, _ = regexp.MatchString(patternStr, profile.City)
		if !check {
			return false
		}
	}
	if profile.AboutUser != "" {
		check, _ = regexp.MatchString(patternStr, profile.AboutUser)
		if !check {
			return false
		}
	}
	if len(profile.Interests) != 0 {
		for _, value := range profile.Interests {
			check, _ = regexp.MatchString(patternStr, value)
			if !check {
				return false
			}
		}
	}
	return true
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
	cook, err := r.Cookie(authCookie)
	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	idFromUrl := mux.Vars(r)["id"]
	checkIdFromUrl, _ := regexp.MatchString("^[0-9]+$", idFromUrl)
	if !checkIdFromUrl {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idFromUrl)
	if errors.Is(err, strconv.ErrSyntax) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	profile, err := handler.ProfileUseCase.Get(cook.Value, id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // вот это нужно проверить
		return
	}

	response, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.Write(response)

}

func (handler *ProfileHandler) GetShortProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	cook, err := r.Cookie(authCookie)
	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	idFromUrl := mux.Vars(r)["id"]
	checkIdFromUrl, _ := regexp.MatchString("^[0-9]+$", idFromUrl)
	if !checkIdFromUrl {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idFromUrl)
	if errors.Is(err, strconv.ErrSyntax) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	profile, err := handler.ProfileUseCase.GetShort(cook.Value, id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // вот это нужно проверить
		return
	}
	response, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (handler *ProfileHandler) ChangeProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	cook, err := r.Cookie(authCookie)
	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	msg, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}
	model := &models.Profile{}

	err = json.Unmarshal(msg, model)
	if err != nil || !checkData(model) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = handler.ProfileUseCase.Change(cook.Value, *model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // вот это нужно проверить
		return
	}
}

func (handler *ProfileHandler) GetCandidateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	cook, err := r.Cookie(authCookie)
	if errors.Is(err, http.ErrNoCookie) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	profile, err := handler.ProfileUseCase.GetCandidates(cook.Value)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // вот это нужно проверить
		return
	}
	response, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.Write(response)
}
