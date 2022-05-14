package impl

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/repositories"
)

type photosUseCase struct {
	repo repositories.PhotosRepository
}

func NewPhotosUseCase(photosRepo repositories.PhotosRepository) *photosUseCase {
	return &photosUseCase{repo: photosRepo}
}

func (useCase *photosUseCase) Create(userId int) (models.PhotoID, error) {
	created, err := useCase.repo.Create(userId)
	if err != nil {
		return models.PhotoID{}, err
	}

	return models.PhotoID{ID: created}, nil
}

func (useCase *photosUseCase) CanSave(id int, userId int) error {
	author, err := useCase.repo.GetAuthor(id)
	if err != nil {
		return err
	}
	if author != userId {
		return http.ErrPhotoChangeForbidden
	}
	saved, err := useCase.repo.IsSaved(id)
	if err != nil {
		return err
	}
	if saved {
		return http.ErrPhotoChangeForbidden
	}
	return nil
}

func (useCase *photosUseCase) Save(id int, path string) error {
	return useCase.repo.Save(id, path)
}

func (useCase *photosUseCase) Read(id int, userId int) (string, error) {
	return useCase.repo.GetPathIfFilled(id)
}

func (useCase *photosUseCase) SetParams(id int, userId int, params models.PhotoParams) error {
	author, err := useCase.repo.GetAuthor(id)
	if err != nil {
		return err
	}
	if author != userId {
		return http.ErrPhotoChangeForbidden
	}
	return useCase.repo.SetParams(id, params)
}

func (useCase *photosUseCase) GetParams(id int) (models.PhotoParams, error) {
	return useCase.repo.GetParams(id)
}

func (useCase *photosUseCase) Delete(id int, userId int) error {
	author, err := useCase.repo.GetAuthor(id)
	if err != nil {
		return err
	}
	if author != userId {
		return http.ErrPhotoChangeForbidden
	}
	return useCase.repo.Delete(id)
}

func (useCase *photosUseCase) GetUserPhotos(userId int) (models.UserPhotos, error) {
	return useCase.repo.GetUserPhotos(userId)
}

func (useCase *photosUseCase) GetUserAvatar(userId int) (models.UserAvatar, error) {
	avatarId, params, err := useCase.repo.GetAvatar(userId)
	if err != nil {
		return models.UserAvatar{}, err
	}
	return models.UserAvatar{Avatar: avatarId, Params: params}, nil
}

func (useCase *photosUseCase) SetUserAvatar(avatar models.UserAvatar, userId int, userIdCookie int) error {
	if userId != userIdCookie {
		return http.ErrPhotoChangeForbidden
	}
	return useCase.repo.SetAvatar(avatar.Avatar, avatar.Params, userId)
}
