package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/usecases"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ProfileHandler struct {
	ProfileUseCase usecases.ProfileUseCases
}

func CreateProfileHandler(useCase usecases.ProfileUseCases) *ProfileHandler {
	return &ProfileHandler{ProfileUseCase: useCase}
}

func (handler *ProfileHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	cook, err := r.Cookie(authCookie)
	if err == http.ErrNoCookie {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	id, _ := mux.Vars(r)["id"]
	convertID, err := strconv.Atoi(id)
	if err == http.ErrNoCookie {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	profile, err := handler.ProfileUseCase.ProfileGet(cook.String(), convertID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // вот это нужно проверить
		return
	}
	response, _ := json.Marshal(profile)
	w.Write(response)
}

func (handler *ProfileHandler) GetShortProfileHandler(w http.ResponseWriter, r *http.Request) {
	cook, err := r.Cookie(authCookie)
	if err == http.ErrNoCookie {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	id, _ := mux.Vars(r)["id"]
	convertID, err := strconv.Atoi(id)
	if err == http.ErrNoCookie {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	profile, err := handler.ProfileUseCase.ShortProfileGet(cook.String(), convertID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // вот это нужно проверить
		return
	}
	response, _ := json.Marshal(profile)
	w.Write(response)
}

func (handler *ProfileHandler) ChangeProfileHandler(w http.ResponseWriter, r *http.Request) {
	cook, err := r.Cookie(authCookie)
	if err == http.ErrNoCookie {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	var model models.Profile
	var msg []byte
	_, err = r.Body.Read(msg)
	if err != nil {
		return
	}
	err = json.Unmarshal(msg, model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	err = handler.ProfileUseCase.ProfileChange(cook.String(), model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // вот это нужно проверить
		return
	}
}

func (handler *ProfileHandler) GetCandidateHandler(w http.ResponseWriter, r *http.Request) {
	cook, err := r.Cookie(authCookie)
	if err == http.ErrNoCookie {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	profile, err := handler.ProfileUseCase.ProfilesCandidateGet(cook.String())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // вот это нужно проверить
		return
	}
	response, _ := json.Marshal(profile)
	w.Write(response)
}
