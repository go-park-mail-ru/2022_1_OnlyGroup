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
const passwordPattern1 = `[a-z]+`
const passwordPattern2 = `[A-Z]+`
const passwordPattern3 = `[0-9]+`

const authCookie = "session"

type AuthHandler struct {
	AuthUseCase usecases.AuthUseCases
}

func CreateAuthHandler(useCase usecases.AuthUseCases) *AuthHandler {
	return &AuthHandler{useCase}
}

func checkValidUserModel(user models.UserAuthInfo) error {
	//processing email
	match, err := regexp.MatchString(emailPattern, user.Email)
	if err != nil || !match {
		return ErrAuthValidationEmail
	}

	//processing password
	if len(user.Password) > 32 || len(user.Password) < 6 {
		return ErrAuthValidationPassword
	}

	match, err = regexp.MatchString(passwordPattern1, user.Password)
	if err != nil || !match {
		return ErrAuthValidationPassword
	}
	match, err = regexp.MatchString(passwordPattern2, user.Password)
	if err != nil || !match {
		return ErrAuthValidationPassword
	}
	match, err = regexp.MatchString(passwordPattern3, user.Password)
	if err != nil || !match {
		return ErrAuthValidationPassword
	}
	return nil
}

func (handler *AuthHandler) GET(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	cook, err := r.Cookie(authCookie)
	if err == http.ErrNoCookie {
		http.Error(w, WrapError2Json(err.Error()), http.StatusUnauthorized)
		return
	}

	userId, err := handler.AuthUseCase.UserAuth(cook.Value)
	if err != nil {
		errCode := ErrorToHTTPCode(err)
		http.Error(w, WrapError2Json(err.Error()), errCode)
		return
	}
	response, _ := json.Marshal(userId)
	w.Write(response)
}

func (handler *AuthHandler) PUT(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, WrapError2Json(err.Error()), http.StatusInternalServerError)
	}
	user := &models.UserAuthInfo{}
	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, WrapError2Json(err.Error()), http.StatusBadRequest)
		return
	}

	err = checkValidUserModel(*user)
	if err != nil {
		http.Error(w, WrapError2Json(err.Error()), http.StatusBadRequest)
		return
	}

	realUser, cook, err := handler.AuthUseCase.UserLogin(*user)
	if err != nil {
		errCode := ErrorToHTTPCode(err)
		http.Error(w, WrapError2Json(err.Error()), errCode)
		return
	}
	cookie := http.Cookie{Name: authCookie, Value: cook, Expires: time.Now().Add(time.Hour * 24 * 7)}
	http.SetCookie(w, &cookie)

	response, _ := json.Marshal(realUser)
	w.Write(response)
}

func (handler *AuthHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	cook, err := r.Cookie(authCookie)
	if err == http.ErrNoCookie {
		http.Error(w, WrapError2Json(err.Error()), http.StatusUnauthorized)
		return
	}

	err = handler.AuthUseCase.UserLogout(cook.Value)
	if err != nil {
		errCode := ErrorToHTTPCode(err)
		http.Error(w, WrapError2Json(err.Error()), errCode)
	}
	cookie := http.Cookie{Name: authCookie, Value: cook.Value, Expires: time.Now().Add(time.Hour * (-1))}
	http.SetCookie(w, &cookie)
}

func (handler *AuthHandler) POST(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, WrapError2Json(err.Error()), http.StatusInternalServerError)
	}
	user := &models.UserAuthInfo{}
	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, WrapError2Json(err.Error()), http.StatusBadRequest)
		return
	}

	err = checkValidUserModel(*user)
	if err != nil {
		http.Error(w, WrapError2Json(err.Error()), http.StatusBadRequest)
		return
	}

	realUser, cook, err := handler.AuthUseCase.UserRegister(*user)
	if err != nil {
		errCode := ErrorToHTTPCode(err)
		http.Error(w, WrapError2Json(err.Error()), errCode)
		return
	}
	cookie := http.Cookie{Name: authCookie, Value: cook, Expires: time.Now().Add(time.Hour * 24 * 7)}
	http.SetCookie(w, &cookie)

	response, _ := json.Marshal(realUser)
	w.Write(response)
}
