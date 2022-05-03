package impl

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/repositories"
)

type interestsUseCaseImpl struct {
	interestsRepo repositories.InterestsRepository
}

func NewInterestsUseCaseImpl(interestsRepo repositories.InterestsRepository) *interestsUseCaseImpl {
	return &interestsUseCaseImpl{interestsRepo: interestsRepo}
}

func (useCase *interestsUseCaseImpl) Get() ([]models.Interest, error) {
	var interests []models.Interest
	interests, err := useCase.interestsRepo.GetInterests()
	if err != nil {
		return nil, err
	}
	return interests, nil
}

func (useCase *interestsUseCaseImpl) Check([]models.Interest) error {

	return nil
}
