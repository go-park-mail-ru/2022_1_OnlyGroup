package usecases

import (
	"2022_1_OnlyGroup_back/app/models"
)

type PhotosUseCase interface {
	Create(userId int) (models.PhotoID, error)
	CanSave(id int, userId int) error
	Save(id int, path string) error
	Read(id int, userId int) (string, error)
	SetParams(id int, userId int, params models.PhotoParams) error
	GetParams(id int) (models.PhotoParams, error)
	Delete(id int, userId int) error
	GetUserPhotos(userId int) (models.UserPhotos, error)
	GetUserAvatar(userId int) (models.UserAvatar, error)
	SetUserAvatar(avatar models.UserAvatar, userId int, userIdCookie int) error
}
