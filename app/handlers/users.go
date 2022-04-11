package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/usecases"
	"encoding/json"
	"fmt"
	"gopkg.in/validator.v2"
	"io"
	"net/http"
	"time"
)

const authCookie = "session"

func identityValidatorErr(err error) error {
	if err != nil {
		errs, ok := err.(validator.ErrorMap)
		if !ok {
			return ErrBaseApp
		}
		for _, value := range errs["Email"] {
			if value != nil {
				return ErrAuthValidationEmail
			}
		}
		for _, value := range errs["Password"] {
			if value != nil {
				return ErrAuthValidationPassword
			}
		}
	}
	return nil
}

type AuthHandler struct {
	AuthUseCase usecases.AuthUseCases
}

func CreateAuthHandler(useCase usecases.AuthUseCases) *AuthHandler {
	return &AuthHandler{useCase}
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

	err = validator.Validate(user)
	if err == ErrBaseApp {
		fmt.Println(ErrBaseApp.Is(err))
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	if err := identityValidatorErr(err); err != nil {
		err := AppErrorFromError(err)
		http.Error(w, err.String(), err.Code)
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

	err = validator.Validate(user)
	if err == ErrBaseApp {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	if err := identityValidatorErr(err); err != nil {
		err := AppErrorFromError(err)
		http.Error(w, err.String(), err.Code)
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
