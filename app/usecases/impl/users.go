package impl

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/repositories"
)

type authUseCaseImpl struct {
	usersRepo    repositories.UsersRepository
	sessionsRepo repositories.SessionsRepository
	profileRepo  repositories.ProfileRepository
}

func NewAuthUseCaseImpl(usersRepo repositories.UsersRepository, sessionsRepo repositories.SessionsRepository, profileRepo repositories.ProfileRepository) *authUseCaseImpl {
	return &authUseCaseImpl{usersRepo: usersRepo, sessionsRepo: sessionsRepo, profileRepo: profileRepo}
}

func (useCase *authUseCaseImpl) UserAuth(Cookie string) (id models.UserID, err error) {
	realId, _, err := useCase.sessionsRepo.Get(Cookie)
	if err != nil {
		return
	}
	return models.UserID{ID: realId}, nil
}

func (useCase *authUseCaseImpl) UserLogin(userInfo models.UserAuthInfo) (id models.UserID, cookie string, err error) {
	realId, err := useCase.usersRepo.Authorize(userInfo.Email, userInfo.Password)
	if err != nil {
		return
	}
	cookie, err = useCase.sessionsRepo.Add(realId, "")
	id.ID = realId
	return
}

func (useCase *authUseCaseImpl) UserRegister(userInfo models.UserAuthInfo) (id models.UserID, cookie string, err error) {
	realId, err := useCase.usersRepo.AddUser(userInfo.Email, userInfo.Password)
	if err != nil {
		return
	}
	err = useCase.profileRepo.AddEmpty(realId)
	if err != nil {
		return
	}
	cookie, err = useCase.sessionsRepo.Add(realId, "")
	id.ID = realId
	return
}

func (useCase *authUseCaseImpl) UserLogout(Cookie string) error {
	return useCase.sessionsRepo.Remove(Cookie)
}

func (useCase *authUseCaseImpl) UserChangePassword(userProfile models.UserAuthProfile, Cookie string) (err error) {
	realIdSession, _, err := useCase.sessionsRepo.Get(Cookie)
	if err != nil {
		return
	}

	realIdAuth, err := useCase.usersRepo.Authorize(userProfile.Email, userProfile.OldPassword)
	if err != nil {
		return
	}

	if realIdAuth != realIdSession {
		err = http.ErrAuthWrongPassword
		return
	}

	err = useCase.usersRepo.ChangePassword(realIdAuth, userProfile.NewPassword)
	return
}
