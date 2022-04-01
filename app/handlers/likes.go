package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/usecases"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

func checkLikesData(likes *models.Likes) bool {
	var check bool
	var err error
	check, err = regexp.MatchString(patternInt, strconv.Itoa(likes.Id))
	if !check || err != nil {
		return false
	}
	check, err = regexp.MatchString(patternInt, strconv.Itoa(likes.Action))
	if !check || err != nil {
		return false
	}
	return true
}

type LikesHandler struct {
	likesUseCase usecases.LikesUseCases
}

func CreateLikesHandler(likesUseCase usecases.LikesUseCases) *LikesHandler {
	return &LikesHandler{likesUseCase: likesUseCase}
}

func (handler *LikesHandler) Set(w http.ResponseWriter, r *http.Request) {
	msg, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}
	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		return
	}
	model := &models.Likes{}
	err = json.Unmarshal(msg, model)
	if err != nil || !checkLikesData(model) {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.code)
		return
	}
	err = handler.likesUseCase.SetAction(cookieId, *model)
	if err != nil {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}

}

func (handler *LikesHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookieId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		return
	}
	model, err := handler.likesUseCase.GetMatched(cookieId)
	if err != nil {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}
	response, err := json.Marshal(model)
	if err != nil {
		appErr := appErrorFromError(err)
		http.Error(w, appErr.String(), appErr.code)
		return
	}
	w.Write(response)
	return
}
