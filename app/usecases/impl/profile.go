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

func (useCase *profileUseCaseImpl) Get(cookieProfileId int, profileId int) (profile models.Profile, err error) {
	if profileId == cookieProfileId {
		profile, err = useCase.profileRepo.GetProfile(cookieProfileId)
		if err != nil {
			return profile, err
		}
		return
	}
	profile, err = useCase.profileRepo.GetProfile(profileId)
	if err != nil {
		return profile, err
	}
	return
}

func (useCase *profileUseCaseImpl) Change(profileId int, profile models.Profile) (err error) {
	err = useCase.profileRepo.ChangeProfile(profileId, profile)
	return
}

func (useCase *profileUseCaseImpl) Delete(profileId int) (err error) {
	err = useCase.profileRepo.DeleteProfile(profileId)
	return
}

func (useCase *profileUseCaseImpl) GetShort(cookieId int, profileId int) (profile models.ShortProfile, err error) {
	if cookieId != profileId {
		return
	}
	profile, err = useCase.profileRepo.GetShortProfile(profileId)
	return
}

func (useCase *profileUseCaseImpl) GetCandidates(profileId int) (candidateProfiles models.VectorCandidate, err error) {

	candidateProfiles, err = useCase.profileRepo.FindCandidateProfile(profileId)
	return
}
