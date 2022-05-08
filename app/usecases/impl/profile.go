package impl

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/repositories"
	"math"
	"strconv"
	"time"
)

type profileUseCaseImpl struct {
	profileRepo repositories.ProfileRepository
}

func NewProfileUseCaseImpl(profileRepo repositories.ProfileRepository) *profileUseCaseImpl {
	return &profileUseCaseImpl{profileRepo: profileRepo}
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
		return profile, http.ErrBaseApp.Wrap(err, "failed convert age")
	}
	profile.Birthday = nil

	return
}

func (useCase *profileUseCaseImpl) Change(profileId int, profile models.Profile) (err error) {
	if profileId != profile.UserId {
		return http.ErrProfileForbiddenChange
	}
	err = useCase.profileRepo.CheckInterests(profile.Interests)
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

func (useCase *profileUseCaseImpl) GetInterest() ([]models.Interest, error) {
	var interests []models.Interest
	interests, err := useCase.profileRepo.GetInterests()
	if err != nil {
		return nil, err
	}
	return interests, nil
}

func (useCase *profileUseCaseImpl) GetDynamicInterests(interest string) ([]models.Interest, error) {
	findInterests, err := useCase.profileRepo.GetDynamicInterest(interest)
	if err != nil {
		return nil, err
	}
	return findInterests, nil
}

func (useCase *profileUseCaseImpl) CheckInterests([]models.Interest) error {
	return nil
}

func (useCase *profileUseCaseImpl) GetFilters(userId int) (models.Filters, error) {
	filters, err := useCase.profileRepo.GetFilters(userId)
	if err != nil {
		return models.Filters{}, err
	}
	return filters, nil
}

func (useCase *profileUseCaseImpl) ChangeFilters(userId int, filters models.Filters) error {
	err := useCase.profileRepo.ChangeFilters(userId, filters)
	return err
}

func (useCase *profileUseCaseImpl) SetAction(userid int, likes models.Likes) (err error) {
	if userid == likes.Id {
		return http.ErrBadRequest
	}
	err = useCase.profileRepo.SetAction(userid, likes)
	if err != nil {
		return err
	}
	return
}

func (useCase *profileUseCaseImpl) GetMatched(userId int) (likesVector models.LikesMatched, err error) {
	likesVector, err = useCase.profileRepo.GetMatched(userId)
	if err != nil {
		return
	}
	return
}

func (useCase *profileUseCaseImpl) AddEmpty(profileId int) (err error) {
	err = useCase.profileRepo.AddEmpty(profileId)
	return err
}
