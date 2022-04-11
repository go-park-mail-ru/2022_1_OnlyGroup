package usecases

import "2022_1_OnlyGroup_back/app/models"

type LikesUseCases interface {
	SetAction(userid int, likes models.Likes) (err error)
	GetMatched(userId int) (likesVector models.LikesMatched, err error)
}
