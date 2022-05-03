package usecases

import "2022_1_OnlyGroup_back/app/models"

type InterestsUseCase interface {
	Get() ([]models.Interest, error)
	GetDynamic(string) ([]models.Interest, error)
	Check([]models.Interest) error
}
