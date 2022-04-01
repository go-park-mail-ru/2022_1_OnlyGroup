package impl

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/repositories"
	"fmt"
)

type likesUseCaseImpl struct {
	likesRepo repositories.LikesRepository
}

func NewLikesUseCaseImpl(likesRepo repositories.LikesRepository) *likesUseCaseImpl {
	return &likesUseCaseImpl{likesRepo: likesRepo}
}

func (useCase *likesUseCaseImpl) SetAction(userid int, likes models.Likes) (err error) {
	if userid == likes.Id {
		return fmt.Errorf("like to own profile failed: %w", handlers.ErrBadRequest)
	}
	err = useCase.likesRepo.SetAction(userid, likes)
	if err != nil {
		return err
	}
	return
}

func (useCase *likesUseCaseImpl) GetMatched(userId int) (likesVector models.LikesMatched, err error) {
	likesVector, err = useCase.likesRepo.GetMatched(userId)
	if err != nil {
		return
	}
	return
}
