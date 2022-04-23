package repositories

import (
	"2022_1_OnlyGroup_back/app/models"
)

type ProfileRepository interface {
	Get(profileId int) (profile models.Profile, err error)
	GetShort(profileId int) (shortProfile models.ShortProfile, err error)
	Change(profileId int, profile models.Profile) (err error)
	Delete(profileId int) (err error)
	Add(profile models.Profile) (err error)
	CheckFiled(profileId int) (err error)
	AddEmpty(profileId int) (err error)

	FindCandidate(profileId int) (candidateProfiles models.VectorCandidate, err error)
}
