package mock

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/pkg/errors"
)

type ProfileMock struct {
	profileRepo []models.Profile
}

func remove(slice []models.Profile, s int) {
	slice = append(slice[:s], slice[s+1:]...)
}

func NewProfileMock() *ProfileMock {
	return &ProfileMock{}
}

func (tables *ProfileMock) GetUserProfile(profileId int) (profile models.Profile, err error) {
	for _, item := range tables.profileRepo {
		if item.UserId == profileId {
			profile = item
			return profile, nil
		}
	}
	return profile, errors.ErrAuthUserNotFound
}

func (tables *ProfileMock) ChangeProfile(profileId int, profile models.Profile) (err error) {
	for _, item := range tables.profileRepo {
		if item.UserId == profileId {
			item.Interests = profile.Interests
			item.FirstName = profile.FirstName
			item.LastName = profile.LastName
			item.Birthday = profile.Birthday
			item.City = profile.City
			item.AboutUser = profile.AboutUser
			item.UserId = profile.UserId
			item.Gender = profile.Gender
			item.Gender = profile.Gender
			return nil
		}
	}
	return errors.ErrAuthUserNotFound
}

func (tables *ProfileMock) DeleteProfile(profileId int) (err error) {
	for count, item := range tables.profileRepo {
		if item.UserId == profileId {
			remove(tables.profileRepo, count)
			return nil
		}
	}
	return errors.ErrAuthUserNotFound
}

func (tables *ProfileMock) AddProfile(profile models.Profile) (err error) {
	tables.profileRepo = append(tables.profileRepo, profile)
	return nil

}

func (tables *ProfileMock) FindCandidateProfile(profileId int) (profile *models.Profile, err error) {

	return profile, nil
}
