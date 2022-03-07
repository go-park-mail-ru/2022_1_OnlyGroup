package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/usecases"
	"encoding/json"
	"errors"
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
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	cook, err := r.Cookie(authCookie)
	if err == http.ErrNoCookie {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	idFromUrl := mux.Vars(r)["idFromUrl"]
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
	if err == http.ErrNoCookie {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	idFromUrl := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idFromUrl)
	if errors.Is(err, strconv.ErrSyntax) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	profile, err := handler.ProfileUseCase.ShortProfileGet(cook.Value, id)
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

func (handler *ProfileHandler) ChangeProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
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
	err = handler.ProfileUseCase.Change(cook.Value, model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // вот это нужно проверить
		return
	}
}

func (handler *ProfileHandler) GetCandidateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	cook, err := r.Cookie(authCookie)
	if err == http.ErrNoCookie {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	profile, err := handler.ProfileUseCase.ProfilesCandidateGet(cook.Value)
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
