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
const passwordPatternLowerCase = `[a-z]+`
const passwordPatternUpperCase = `[A-Z]+`
const passwordPatternNumber = `[0-9]+`
const passwordMinLength = 6
const passwordMaxLength = 32
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
	if len(user.Password) > passwordMaxLength || len(user.Password) < passwordMinLength {
		return ErrAuthValidationPassword
	}

	match, err = regexp.MatchString(passwordPatternLowerCase, user.Password)
	if err != nil || !match {
		return ErrAuthValidationPassword
	}
	match, err = regexp.MatchString(passwordPatternUpperCase, user.Password)
	if err != nil || !match {
		return ErrAuthValidationPassword
	}
	match, err = regexp.MatchString(passwordPatternNumber, user.Password)
	if err != nil || !match {
		return ErrAuthValidationPassword
	}
	return nil
}

func (handler *AuthHandler) GET(w http.ResponseWriter, r *http.Request) {
	cook, err := r.Cookie(authCookie)
	if err == http.ErrNoCookie {
		http.Error(w, ErrAuthRequired.String(), ErrAuthRequired.Code)
		return
	}

	userId, err := handler.AuthUseCase.UserAuth(cook.Value)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	response, _ := json.Marshal(userId)
	w.Write(response)
}

func (handler *AuthHandler) PUT(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
	}
	user := &models.UserAuthInfo{}
	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.Code)
		return
	}

	err = checkValidUserModel(*user)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	realUser, cook, err := handler.AuthUseCase.UserLogin(*user)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
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
		http.Error(w, ErrAuthRequired.String(), ErrAuthRequired.Code)
		return
	}

	err = handler.AuthUseCase.UserLogout(cook.Value)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	cookie := http.Cookie{Name: authCookie, Value: cook.Value, Expires: time.Now().Add(time.Hour * (-1))}
	http.SetCookie(w, &cookie)
}

func (handler *AuthHandler) POST(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	user := &models.UserAuthInfo{}
	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.Code)
		return
	}

	err = checkValidUserModel(*user)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	realUser, cook, err := handler.AuthUseCase.UserRegister(*user)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	cookie := http.Cookie{Name: authCookie, Value: cook, Expires: time.Now().Add(time.Hour * 24 * 7)}
	http.SetCookie(w, &cookie)

	response, _ := json.Marshal(realUser)
	w.Write(response)
}
