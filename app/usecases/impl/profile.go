package impl

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/repositories"
	"math"
	"strconv"
	"time"
)

type profileUseCaseImpl struct {
	profileRepo   repositories.ProfileRepository
	interestsRepo repositories.InterestsRepository
}

func NewProfileUseCaseImpl(profileRepo repositories.ProfileRepository, interestsRepo repositories.InterestsRepository) *profileUseCaseImpl {
	return &profileUseCaseImpl{profileRepo: profileRepo, interestsRepo: interestsRepo}
}

func (useCase *profileUseCaseImpl) Get(cookieProfileId int, profileId int) (profile models.Profile, err error) {
	profile, err = useCase.profileRepo.Get(profileId)
	if err != nil {
		return profile, err
	}
	if profileId == cookieProfileId {
		return
	}
	if profile.Birthday == nil {
		return
	}

	age := int(math.Floor(time.Now().Sub(*profile.Birthday).Hours() / 24 / 365))
	profile.Age = strconv.Itoa(age)
	if err != nil {
		return profile, handlers.ErrBaseApp.Wrap(err, "failed convert age")
	}
	profile.Birthday = nil

	return
}

func (useCase *profileUseCaseImpl) Change(profileId int, profile models.Profile) (err error) {
	if profileId != profile.UserId {
		return handlers.ErrProfileForbiddenChange
	}
	err = useCase.interestsRepo.CheckInterests(profile.Interests)
	if err != nil {
		return err
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
