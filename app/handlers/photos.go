package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/usecases"
	"2022_1_OnlyGroup_back/pkg/fileService"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
)

type PhotosHandler struct {
	PhotosUseCase     usecases.PhotosUseCase
	PhotosFileService fileService.PhotosStorage
	supportedTypes    map[string]bool
}

func CreatePhotosHandler(useCase usecases.PhotosUseCase, fileService fileService.PhotosStorage) *PhotosHandler {
	typesMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".webp": true,
	}
	return &PhotosHandler{PhotosFileService: fileService, PhotosUseCase: useCase, supportedTypes: typesMap}
}

func (handler *PhotosHandler) getMimeFromPath(path string) string {
	lastDot := strings.LastIndexByte(path, '.')
	if lastDot == -1 {
		panic("last dot not found in path: " + path)
	}
	ext := path[lastDot:]
	return mime.TypeByExtension(ext)
}

func (handler *PhotosHandler) getPathFromMime(path string, mimeString string) (string, error) {
	extension, err := mime.ExtensionsByType(mimeString)
	if err != nil {
		return "", ErrUnsupportedType
	}
	if extension == nil {
		return "", ErrUnsupportedType
	}
	b, has := handler.supportedTypes[extension[0]]
	if !has || !b {
		return "", ErrUnsupportedType
	}
	return path + extension[0], nil
}

func (handler *PhotosHandler) getPhotoPath(photoId int) string {
	return strconv.Itoa(photoId)
}

func (handler PhotosHandler) GETPhoto(w http.ResponseWriter, r *http.Request) {
	photoId, err := getIdFromUrl(r)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	ctx := r.Context()
	userId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		http.Error(w, ErrBaseApp.String(), ErrBaseApp.Code)
		return
	}

	path, err := handler.PhotosUseCase.Read(photoId, userId)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	readFile, size, err := handler.PhotosFileService.Read(path)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	defer readFile.Close()

	contentType := handler.getMimeFromPath(path)
	w.Header().Add("Content-Type", contentType)
	w.Header().Add("Content-Length", strconv.Itoa(int(size)))

	_, err = io.Copy(w, readFile)
	if err != nil {
		AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
	}
}

func (handler PhotosHandler) GETParams(w http.ResponseWriter, r *http.Request) {
	photoId, err := getIdFromUrl(r)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	params, err := handler.PhotosUseCase.GetParams(photoId)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	marshal, _ := json.Marshal(params)
	w.Write(marshal)
}

func (handler PhotosHandler) POST(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		http.Error(w, ErrBaseApp.String(), ErrBaseApp.Code)
		return
	}
	createdPhoto, err := handler.PhotosUseCase.Create(userId)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	response, _ := json.Marshal(createdPhoto)
	w.Write(response)
}

func (handler PhotosHandler) POSTPhoto(w http.ResponseWriter, r *http.Request) {
	photoId, err := getIdFromUrl(r)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	ctx := r.Context()
	userId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		http.Error(w, ErrBaseApp.String(), ErrBaseApp.Code)
		return
	}

	err = handler.PhotosUseCase.CanSave(photoId, userId)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	mimeString := r.Header.Get("Content-Type")
	realSizeString := r.Header.Get("Content-Length")
	realSize, err := strconv.Atoi(realSizeString)
	if err != nil {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.Code)
		return
	}

	path, err := handler.getPathFromMime(handler.getPhotoPath(photoId), mimeString)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	writeFile, err := handler.PhotosFileService.Write(path)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	defer writeFile.Close()
	defer r.Body.Close()

	reallyWritten, err := io.Copy(writeFile, r.Body)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	if int(reallyWritten) != realSize {
		err = handler.PhotosFileService.Remove(path)
		if err != nil {
			AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		}
		http.Error(w, ErrContentLengthMismatched.String(), ErrContentLengthMismatched.Code)
		return
	}

	err = handler.PhotosUseCase.Save(photoId, path)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler PhotosHandler) PUTParams(w http.ResponseWriter, r *http.Request) {
	photoId, err := getIdFromUrl(r)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	photoParams := &models.PhotoParams{}
	err = json.Unmarshal(body, photoParams)
	if err != nil {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.Code)
		return
	}

	ctx := r.Context()
	userId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		http.Error(w, ErrBaseApp.String(), ErrBaseApp.Code)
		return
	}

	err = handler.PhotosUseCase.SetParams(photoId, userId, *photoParams)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler PhotosHandler) DELETE(w http.ResponseWriter, r *http.Request) {
	photoId, err := getIdFromUrl(r)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	ctx := r.Context()
	userId, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		http.Error(w, ErrBaseApp.String(), ErrBaseApp.Code)
		return
	}

	err = handler.PhotosUseCase.Delete(photoId, userId)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler PhotosHandler) GETAll(w http.ResponseWriter, r *http.Request) {
	userId, err := getIdFromUrl(r)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	photos, err := handler.PhotosUseCase.GetUserPhotos(userId)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	response, _ := json.Marshal(photos)
	w.Write(response)
}

func (handler PhotosHandler) GETAvatar(w http.ResponseWriter, r *http.Request) {
	userId, err := getIdFromUrl(r)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	avatar, err := handler.PhotosUseCase.GetUserAvatar(userId)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}

	response, _ := json.Marshal(avatar)
	w.Write(response)
}

func (handler PhotosHandler) PUTAvatar(w http.ResponseWriter, r *http.Request) {
	userId, err := getIdFromUrl(r)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	ctx := r.Context()
	userIdCookie, ok := ctx.Value(userIdContextKey).(int)
	if !ok {
		http.Error(w, ErrBaseApp.String(), ErrBaseApp.Code)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	avatar := models.UserAvatar{}
	err = json.Unmarshal(body, &avatar)
	if err != nil {
		http.Error(w, ErrBadRequest.String(), ErrBadRequest.Code)
		return
	}

	err = handler.PhotosUseCase.SetUserAvatar(avatar, userId, userIdCookie)
	if err != nil {
		appErr := AppErrorFromError(err).LogServerError(r.Context().Value(requestIdContextKey))
		http.Error(w, appErr.String(), appErr.Code)
		return
	}
	w.WriteHeader(http.StatusOK)
}
