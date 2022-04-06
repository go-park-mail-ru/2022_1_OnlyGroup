package impl

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/repositories"
)

type profileUseCaseImpl struct {
	profileRepo repositories.ProfileRepository
}

func NewProfileUseCaseImpl(profileRepo repositories.ProfileRepository) *profileUseCaseImpl {
	return &profileUseCaseImpl{profileRepo: profileRepo}
}

func (useCase *profileUseCaseImpl) Get(cookieProfileId int, profileId int) (profile models.Profile, err error) {
	if profileId == cookieProfileId {
		profile, err = useCase.profileRepo.Get(cookieProfileId)
		if err != nil {
			return profile, err
		}
		return
	}
	profile, err = useCase.profileRepo.Get(profileId)
	if err != nil {
		return profile, err
	}
	return
}

func (useCase *profileUseCaseImpl) Change(profileId int, profile models.Profile) (err error) {
	if profileId != profile.UserId {
		return handlers.ErrProfileForbiddenChange
	}
	err = useCase.profileRepo.Change(profileId, profile)
	return
}

func (useCase *profileUseCaseImpl) Delete(profileId int) (err error) {
	err = useCase.profileRepo.Delete(profileId)
	return
}

func (useCase *profileUseCaseImpl) GetShort(cookieId int, profileId int) (profile models.ShortProfile, err error) {
	if cookieId == profileId {
		profile, err = useCase.profileRepo.GetShort(cookieId)
		return
	}
	profile, err = useCase.profileRepo.GetShort(profileId)
	return
}

func (useCase *profileUseCaseImpl) GetCandidates(profileId int) (candidateProfiles models.VectorCandidate, err error) {

	candidateProfiles, err = useCase.profileRepo.FindCandidate(profileId)
	return
}
