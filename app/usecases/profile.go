package usecases

import "2022_1_OnlyGroup_back/app/models"

type ProfileUseCases interface {
	Get(cookieId int, candidateId int) (profile models.Profile, err error)
	Change(profileId int, profile models.Profile) (err error)
	GetShort(cookieId int, profileId int) (profile models.ShortProfile, err error)
	Delete(profileId int) (err error)

	GetCandidates(profileId int) (candidateProfiles models.VectorCandidate, err error)
}
