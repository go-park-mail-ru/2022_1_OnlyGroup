package usecases

import "2022_1_OnlyGroup_back/app/models"

type AuthUseCases interface {
	Create(Cookie string) (id models.UserID, err error)
	Login(userInfo models.UserAuthInfo) (id models.UserID, cookie string, err error)
	Register(userInfo models.UserAuthInfo) (id models.UserID, cookie string, err error)
	Logout(Cookie string) error
	ChangePassword(userProfile models.UserAuthProfile, Cookie string) (err error)
}
