package mock

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/pkg/errors"
)

type ProfileMock struct {
	profileRepo []models.Profile
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
	for id, item := range tables.profileRepo {
		if item.UserId == profileId {
			tables.profileRepo[id].Interests = profile.Interests
			tables.profileRepo[id].FirstName = profile.FirstName
			tables.profileRepo[id].LastName = profile.LastName
			tables.profileRepo[id].Birthday = profile.Birthday
			tables.profileRepo[id].City = profile.City
			tables.profileRepo[id].AboutUser = profile.AboutUser
			tables.profileRepo[id].UserId = profile.UserId
			tables.profileRepo[id].Gender = profile.Gender
			tables.profileRepo[id].Gender = profile.Gender
			return nil
		}
	}
	return errors.ErrAuthUserNotFound
}

func (tables *ProfileMock) DeleteProfile(profileId int) (err error) {
	for count, item := range tables.profileRepo {
		if item.UserId == profileId {
			tables.profileRepo = append(tables.profileRepo[:count], tables.profileRepo[count+1:]...)
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
