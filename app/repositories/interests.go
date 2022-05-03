package repositories

import "2022_1_OnlyGroup_back/app/models"

type InterestsRepository interface {
	GetInterests() ([]models.Interest, error)
	CheckInterests([]models.Interest) error
}
