package repositories

import "2022_1_OnlyGroup_back/app/models"

type PhotosRepository interface {
	Create(userId int) (int, error)
	Save(id int, path string) error
	IsSaved(id int) (bool, error)
	GetAuthor(id int) (int, error)
	SetParams(id int, params models.PhotoParams) error
	GetParams(id int) (models.PhotoParams, error)
	GetPath(id int) (string, error)
	Delete(id int) error
	GetAvatar(userId int) (int, models.PhotoParams, error)
	SetAvatar(id int, params models.PhotoParams, userId int) error
	GetUserPhotos(userId int) (models.UserPhotos, error)
}
