package mock

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/pkg/errors"
	"github.com/bxcodec/faker/v3"
)

type ProfileMock struct {
	profileRepo []models.Profile
}

func NewProfileMock() *ProfileMock {
	mock := ProfileMock{}
	for i := 0; i < 6; i++ {
		mock.AddProfile(models.Profile{FirstName: faker.FirstName(), LastName: faker.LastName(), Birthday: faker.Date(), City: "Moscow", Interests: []string{faker.Word(), faker.Word()}, AboutUser: faker.Sentence(), UserId: i, Gender: faker.Gender()})
	}
	return &mock
}

func (tables *ProfileMock) GetUserProfile(profileId int) (profile models.Profile, err error) {
	for _, item := range tables.profileRepo {
		if item.UserId == profileId {
			profile = item
			return profile, nil
		}
	}
	return profile, errors.ErrProfileNotFound
}

func (tables *ProfileMock) GetUserShortProfile(profileId int) (shortProfile models.ShortProfile, err error) {
	for _, item := range tables.profileRepo {
		if item.UserId == profileId {
			shortProfile.FirstName = item.FirstName
			shortProfile.LastName = item.LastName
			shortProfile.City = item.City
			return shortProfile, nil
		}
	}
	return shortProfile, errors.ErrProfileNotFound
}

var count = 0

func (tables *ProfileMock) ChangeProfile(profileId int, profile models.Profile) (err error) {
	for id, item := range tables.profileRepo {
		if item.UserId == profileId {
			tables.profileRepo[id].Interests = profile.Interests
			tables.profileRepo[id].FirstName = profile.FirstName
			tables.profileRepo[id].LastName = profile.LastName
			tables.profileRepo[id].Birthday = profile.Birthday
			tables.profileRepo[id].City = profile.City
			tables.profileRepo[id].AboutUser = profile.AboutUser
			//tables.profileRepo[id].UserId = profile.UserId
			tables.profileRepo[id].Gender = profile.Gender
			tables.profileRepo[id].Gender = profile.Gender
			return nil
		}
	}
	return errors.ErrProfileNotFound
}

func (tables *ProfileMock) DeleteProfile(profileId int) (err error) {
	if len(tables.profileRepo) == 0 {
		return errors.ErrMockIsEmpty
	}
	for count, item := range tables.profileRepo {
		if item.UserId == profileId {
			tables.profileRepo = append(tables.profileRepo[:count], tables.profileRepo[count+1:]...)
			return nil
		}
	}
	return errors.ErrProfileNotFound
}

func (tables *ProfileMock) AddProfile(profile models.Profile) (err error) {
	tables.profileRepo = append(tables.profileRepo, profile)
	return nil
}

func (tables *ProfileMock) AddEmptyProfile(profileId int) (err error) {
	profile := models.Profile{FirstName: "", LastName: "", Birthday: "", City: "", Interests: []string{}, AboutUser: "", UserId: profileId, Gender: ""}
	tables.profileRepo = append(tables.profileRepo, profile)
	return nil
}

func (tables *ProfileMock) FindCandidateProfile(profileId int) (candidateProfiles models.VectorCandidate, err error) {
	if len(tables.profileRepo) == 0 {
		return candidateProfiles, errors.ErrProfileNotFound
	}
	for i := 0; i < 3; i++ {
		if count == len(tables.profileRepo) {
			count = 0
		}
		candidateProfiles.Candidates = append(candidateProfiles.Candidates, tables.profileRepo[count].UserId)
		count += 1
	}
	return candidateProfiles, nil
}

func (tables *ProfileMock) CheckProfileFiled(profileId int) (err error) {
	for _, item := range tables.profileRepo {
		if item.UserId == profileId {
			if item.Gender == "" ||
				item.City == "" ||
				item.LastName == "" ||
				item.AboutUser == "" ||
				item.Birthday == "" ||
				len(item.Interests) == 0 ||
				item.FirstName == "" {
				return errors.ErrProfileNotFiled
			} else {
				return nil
			}
		}
	}
	return errors.ErrProfileNotFound
}
