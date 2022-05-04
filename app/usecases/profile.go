package usecases

import "2022_1_OnlyGroup_back/app/models"

type ProfileUseCases interface {
	Get(cookieId int, candidateId int) (profile models.Profile, err error)
	Change(profileId int, profile models.Profile) (err error)
	GetShort(cookieId int, profileId int) (profile models.ShortProfile, err error)
	Delete(profileId int) (err error)

	GetCandidates(profileId int) (candidateProfiles models.VectorCandidate, err error)

	GetInterest() ([]models.Interest, error)
	GetDynamicInterests(string) ([]models.Interest, error)
	CheckInterests([]models.Interest) error

	GetFilters(userId int) (models.Filters, error)
	ChangeFilters(userId int, filters models.Filters) error
}
