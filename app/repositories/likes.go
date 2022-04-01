package repositories

import "2022_1_OnlyGroup_back/app/models"

type LikesRepository interface {
	SetAction(profileId int, likes models.Likes) (err error)
	GetMatched(profileId int) (likesVector models.LikesMatched, err error)
}
