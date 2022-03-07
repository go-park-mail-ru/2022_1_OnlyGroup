package usecases

import "2022_1_OnlyGroup_back/app/models"

type ProfileUseCases interface {
	ProfileGet(cookies string, candidateId int) (profile models.Profile, err error)
	ProfileChange(cookies string, profile models.Profile) (err error)
	ShortProfileGet(cookies string, profileId int) (profile models.ShortProfile, err error)
	ProfileDelete(cookies string) (err error)

	ProfilesCandidateGet(cookies string) (candidateProfiles models.VectorCandidate, err error)
}
