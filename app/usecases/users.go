package usecases

import "2022_1_OnlyGroup_back/app/models"

type AuthUseCases interface {
	UserAuth(Cookie string) (id models.UserID, err error)
	UserLogin(userInfo models.UserAuthInfo) (id models.UserID, cookie string, err error)
	UserRegister(userInfo models.UserAuthInfo) (id models.UserID, cookie string, err error)
	UserLogout(Cookie string) error
	UserChangePassword(userProfile models.UserAuthProfile, Cookie string) (err error)
}
