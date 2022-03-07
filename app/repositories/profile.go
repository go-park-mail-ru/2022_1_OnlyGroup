package repositories

import "2022_1_OnlyGroup_back/app/models"

type ProfileRepository interface {
	GetUserProfile(profileId int) (profile models.Profile, err error)
	GetUserShortProfile(profileId int) (shortProfile models.ShortProfile, err error)
	ChangeProfile(profileId int, profile models.Profile) (err error)
	DeleteProfile(profileId int) (err error)
	AddProfile(profile models.Profile) (err error)
	CheckProfileFiled(profileId int) (err error)
	AddEmptyProfile(profileId int) (err error)

	FindCandidateProfile(profileId int) (candidateProfiles models.VectorCandidate, err error)
}
