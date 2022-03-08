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

func (useCase *profileUseCaseImpl) Get(cookies string, profileId int) (profile models.Profile, err error) {
	profileIdCheck, err := useCase.authRepo.GetIdBySession(cookies)
	if err != nil {
		return
	}

	if profileId == profileIdCheck {
		profile, err = useCase.profileRepo.GetProfile(profileIdCheck)
		return
	}
	profile, err = useCase.profileRepo.GetProfile(profileId)
	if err != nil {
		return
	}
	return
}

func (useCase *profileUseCaseImpl) Change(cookies string, profile models.Profile) (err error) {
	profileId, err := useCase.authRepo.GetIdBySession(cookies)
	if err != nil {
		return
	}
	err = useCase.profileRepo.ChangeProfile(profileId, profile)
	return
}

func (useCase *profileUseCaseImpl) Delete(cookies string) (err error) {
	profileId, err := useCase.authRepo.GetIdBySession(cookies)
	if err != nil {
		return
	}
	err = useCase.profileRepo.DeleteProfile(profileId)
	return
}

func (useCase *profileUseCaseImpl) GetShort(cookies string, profileId int) (profile models.ShortProfile, err error) {
	checkProfileId, err := useCase.authRepo.GetIdBySession(cookies)
	if err != nil {
		return
	}
	if checkProfileId != profileId {
		return
	}
	profile, err = useCase.profileRepo.GetShortProfile(profileId)
	return
}

func (useCase *profileUseCaseImpl) GetCandidates(cookies string) (candidateProfiles models.VectorCandidate, err error) {
	profileId, err := useCase.authRepo.GetIdBySession(cookies)
	if err != nil {
		return
	}
	candidateProfiles, err = useCase.profileRepo.FindCandidateProfile(profileId)
	return
}
