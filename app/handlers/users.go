package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/usecases"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"time"
)

const emailPattern = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
const authCookie = "session"

type AuthHandler struct {
	AuthUseCase usecases.AuthUseCases
}

func CreateAuthHandler(useCase usecases.AuthUseCases) *AuthHandler {
	return &AuthHandler{useCase}
}

func checkValidUserModel(user models.UserAuthInfo) bool {
	match, err := regexp.MatchString(emailPattern, user.Email)
	if err != nil || !match {
		return false
	}
	return true
}

func (handler *AuthHandler) GET(w http.ResponseWriter, r *http.Request) {
	cook, err := r.Cookie(authCookie)
	if err == http.ErrNoCookie {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userId, err := handler.AuthUseCase.UserAuth(cook.Value)
	if err != nil {
		errCode := ErrorToHTTPCode(err)
		http.Error(w, http.StatusText(errCode), errCode)
		return
	}
	response, _ := json.Marshal(userId)
	w.Write(response)
}

func (handler *AuthHandler) PUT(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	user := &models.UserAuthInfo{}
	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !checkValidUserModel(*user) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	realUser, cook, err := handler.AuthUseCase.UserLogin(*user)
	if err != nil {
		errCode := ErrorToHTTPCode(err)
		http.Error(w, http.StatusText(errCode), errCode)
		return
	}
	cookie := http.Cookie{Name: authCookie, Value: cook, Expires: time.Now().Add(time.Hour * 24 * 7)}
	http.SetCookie(w, &cookie)

	response, _ := json.Marshal(realUser)
	w.Write(response)
}

func (handler *AuthHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	cook, err := r.Cookie(authCookie)
	if err == http.ErrNoCookie {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err = handler.AuthUseCase.UserLogout(cook.Value)
	if err != nil {
		errCode := ErrorToHTTPCode(err)
		http.Error(w, http.StatusText(errCode), errCode)
	}
	cookie := http.Cookie{Name: authCookie, Value: cook.Value, Expires: time.Now().Add(time.Hour * (-1))}
	http.SetCookie(w, &cookie)
}

func (handler *AuthHandler) POST(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	user := &models.UserAuthInfo{}
	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !checkValidUserModel(*user) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	realUser, cook, err := handler.AuthUseCase.UserRegister(*user)
	if err != nil {
		errCode := ErrorToHTTPCode(err)
		http.Error(w, http.StatusText(errCode), errCode)
		return
	}
	cookie := http.Cookie{Name: authCookie, Value: cook, Expires: time.Now().Add(time.Hour * 24 * 7)}
	http.SetCookie(w, &cookie)

	response, _ := json.Marshal(realUser)
	w.Write(response)
}
