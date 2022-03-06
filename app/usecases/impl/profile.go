package impl

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/repositories"
)

type profileUseCaseImpl struct {
	profileRepo repositories.ProfileRepository
	authRepo    repositories.AuthRepository
}

func NewProfileUseCaseImpl(profileRepo repositories.ProfileRepository, authRepo repositories.AuthRepository) *profileUseCaseImpl {
	return &profileUseCaseImpl{profileRepo: profileRepo, authRepo: authRepo}
}

func (useCase *profileUseCaseImpl) ProfileGet(cookies string) (profile models.Profile, err error) {
	profileId, err := useCase.authRepo.GetIdBySession(cookies)
	if err != nil {
		return
	}
	profile, err = useCase.profileRepo.GetUserProfile(profileId)
	if err != nil {
		return
	}

	return
}
func (useCase *profileUseCaseImpl) ProfileChange(cookies string, profile models.Profile) (err error) {
	profileId, err := useCase.authRepo.GetIdBySession(cookies)
	if err != nil {
		return
	}
	err = useCase.profileRepo.ChangeProfile(profileId, profile)
	return
}

func (useCase *profileUseCaseImpl) ShortProfileGet(cookies string) (profile models.ShortProfile, err error) {
	profileId, err := useCase.authRepo.GetIdBySession(cookies)
	if err != nil {
		return
	}
	profile, err = useCase.profileRepo.GetUserShortProfile(profileId)
	return
}
func (useCase *profileUseCaseImpl) ProfileDelete(cookies string) (err error) {
	profileId, err := useCase.authRepo.GetIdBySession(cookies)
	if err != nil {
		return
	}
	err = useCase.profileRepo.DeleteProfile(profileId)
	return
}

func (useCase *profileUseCaseImpl) ProfileCandidateGet(cookies string) (candidateProfile models.Profile, err error) {
	profileId, err := useCase.authRepo.GetIdBySession(cookies)
	if err != nil {
		return
	}
	candidateProfile, err = useCase.profileRepo.FindCandidateProfile(profileId)
	return
}
