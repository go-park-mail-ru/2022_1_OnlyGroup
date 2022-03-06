package usecases

import "2022_1_OnlyGroup_back/app/models"

type ProfileUseCases interface {
	ProfileGet(cookies string) (profile models.Profile, err error)
	ProfileChange(cookies string, profile models.Profile) (err error)
	ShortProfileGet(cookies string) (profile models.ShortProfile, err error)
	ProfileDelete(cookies string) (err error)

	ProfileCandidateGet(cookies string) (candidateProfile models.Profile, err error)
}
