package impl

import (
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/repositories"
	"2022_1_OnlyGroup_back/pkg/errors"
)

type authUseCaseImpl struct {
	authRepo    repositories.AuthRepository
	profileRepo repositories.ProfileRepository
}

func NewAuthUseCaseImpl(authRepo repositories.AuthRepository, profileRepo repositories.ProfileRepository) *authUseCaseImpl {
	return &authUseCaseImpl{authRepo: authRepo, profileRepo: profileRepo}
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
	id.ID = realId
	return
}

func (useCase *authUseCaseImpl) UserRegister(userInfo models.UserAuthInfo) (id models.UserID, cookie string, err error) {
	realId, err := useCase.authRepo.AddUser(userInfo.Email, userInfo.Password)
	if err != nil {
		return
	}
	err = useCase.profileRepo.AddEmptyProfile(realId)
	if err != nil {
		return
	}
	cookie, err = useCase.authRepo.AddSession(realId)
	id.ID = realId
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
		err = errors.ErrAuthWrongPassword
		return
	}
	err = useCase.authRepo.ChangePassword(realIdAuth, userProfile.NewPassword)
	return
}
