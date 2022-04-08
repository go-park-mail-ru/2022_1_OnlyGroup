package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/usecases"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/validator.v2"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
)

const patternInt = "^[0-9]+$"

func setValidators() {
	validator.SetValidationFunc("interests", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.Kind() != reflect.Slice {
			return validator.ErrUnsupported
		}
		if v.IsNil() {
			return nil
		}
		nVal := val.([]string)
		for _, value := range nVal {
			if len(value) > models.InterestSize {
				return validator.ErrLen
			}
		}
		return nil
	})
	validator.SetValidationFunc("birthday", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.Kind() != reflect.String {
			return validator.ErrUnsupported
		}
		if v.IsZero() {
			return nil
		}
		nVal := val.(string)
		if len(nVal) > models.BirthdaySize {
			return validator.ErrLen
		}
		check, err := regexp.MatchString(models.BirthdayRexexp, nVal)
		if err != nil {
			return ErrBaseApp
		}
		if !check {
			return validator.ErrRegexp
		}
		return nil
	})
	validator.SetValidationFunc("password", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.Kind() != reflect.String {
			return validator.ErrUnsupported
		}
		if v.IsZero() {
			return nil
		}
		nVal := val.(string)

		if len(nVal) > models.PasswordMaxLength || len(nVal) < models.PasswordMinLength {
			return validator.ErrLen
		}
		match, err := regexp.MatchString(models.PasswordPatternLowerCase, nVal)
		if err != nil || !match {
			return validator.ErrRegexp
		}
		match, err = regexp.MatchString(models.PasswordPatternUpperCase, nVal)
		if err != nil || !match {
			return validator.ErrRegexp
		}
		match, err = regexp.MatchString(models.PasswordPatternNumber, nVal)
		if err != nil || !match {
			return validator.ErrRegexp
		}
		return nil

		return nil
	})
}

func sanitizeProfileModel(profile *models.Profile) {
	sanitizer := bluemonday.UGCPolicy()
	for idx, value := range profile.Interests {
		profile.Interests[idx] = sanitizer.Sanitize(value)
	}
	profile.Birthday = sanitizer.Sanitize(profile.Birthday)
	profile.FirstName = sanitizer.Sanitize(profile.FirstName)
	profile.AboutUser = sanitizer.Sanitize(profile.AboutUser)
	profile.LastName = sanitizer.Sanitize(profile.LastName)
}

func sanitizeShortProfileModel(profile *models.ShortProfile) {
	sanitizer := bluemonday.UGCPolicy()
	profile.City = sanitizer.Sanitize(profile.City)
	profile.FirstName = sanitizer.Sanitize(profile.FirstName)
	profile.LastName = sanitizer.Sanitize(profile.LastName)
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
	setValidators()
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
	sanitizeProfileModel(&profile)
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
	sanitizeShortProfileModel(&model)
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
	model := models.Profile{}

	err = json.Unmarshal(msg, model)
	if err != nil {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.Code)
		return
	}
	sanitizeProfileModel(&model)
	if err = validator.Validate(model); err != nil {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.Code)
	}
	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	err = handler.ProfileUseCase.Change(cookieId, model)
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
