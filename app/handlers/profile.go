package handlers

import (
	"2022_1_OnlyGroup_back/app/usecases"
	"encoding/json"
	"net/http"
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

	profile, _ := handler.ProfileUseCase.ProfileGet(cook.String())
	response, _ := json.Marshal(profile)
	w.Write(response)
}

func (handler *ProfileHandler) ChangeProfileHandler(w http.ResponseWriter, r *http.Request) {

}

func (handler *ProfileHandler) GetCandidateHandler(w http.ResponseWriter, r *http.Request) {

}
