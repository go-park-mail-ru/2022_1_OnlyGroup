package impl

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/repositories"
)

type authUseCaseImpl struct {
	authRepo repositories.AuthRepository
}

func NewAuthUseCaseImpl(authRepo repositories.AuthRepository) *authUseCaseImpl {
	return &authUseCaseImpl{authRepo: authRepo}
}

func (useCase *authUseCaseImpl) UserAuth(Cookie string) (id models.UserID, err error) {
	realId, err := useCase.authRepo.GetIdBySession(Cookie)
	if err != nil {
		return
	}
	return models.UserID{ID: realId}, nil
}

func (useCase *authUseCaseImpl) UserLogin(userInfo models.UserAuthInfo) (id models.UserID, cookie string, err error) {
	realId, err := useCase.authRepo.Authorize(userInfo.Email, userInfo.Password)
	if err != nil {
		return
	}
	cookie, err = useCase.authRepo.AddSession(realId)
	return
}

func (useCase *authUseCaseImpl) UserRegister(userInfo models.UserAuthInfo) (id models.UserID, cookie string, err error) {
	realId, err := useCase.authRepo.AddUser(userInfo.Email, userInfo.Password)
	if err != nil {
		return
	}
	cookie, err = useCase.authRepo.AddSession(realId)
	return
}

func (useCase *authUseCaseImpl) UserLogout(Cookie string) error {
	return useCase.authRepo.RemoveSession(Cookie)
}

func (useCase *authUseCaseImpl) UserChangePassword(userProfile models.UserAuthProfile, Cookie string) (err error) {
	realIdSession, err := useCase.authRepo.GetIdBySession(Cookie)
	if err != nil {
		return
	}

	realIdAuth, err := useCase.authRepo.Authorize(userProfile.Email, userProfile.OldPassword)
	if err != nil {
		return
	}

	if realIdAuth != realIdSession {
		err = handlers.ErrAuthWrongPassword
		return
	}
	err = useCase.authRepo.ChangePassword(realIdAuth, userProfile.NewPassword)
	return
}
