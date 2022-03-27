package impl

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/tests/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthAuthOk(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	expectedUserModel := models.UserID{ID: testingID}

	sessionsMock.EXPECT().GetIdBySession(testingSecret).Return(testingID, "", nil)
	actualUserModel, err := testingUseCase.UserAuth(testingSecret)

	assert.Equal(t, expectedUserModel, actualUserModel, "models mismatched")
	assert.Equal(t, err, nil)
}

func TestAuthSessionNotFound(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingSecret = "ifiewhufhbbjwdbnfnmwe"

	sessionsMock.EXPECT().GetIdBySession(testingSecret).Return(0, "", handlers.ErrAuthSessionNotFound)
	_, err := testingUseCase.UserAuth(testingSecret)

	assert.Equal(t, err, handlers.ErrAuthSessionNotFound)
}

func TestLoginLoginOk(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}
	expectedUserModel := models.UserID{ID: testingID}

	usersMock.EXPECT().Authorize(testingEmail, testingPassword).Return(testingID, nil)
	sessionsMock.EXPECT().AddSession(testingID, "").Return(testingSecret, nil)
	actualUserModel, actualSecret, err := testingUseCase.UserLogin(testUserModel)

	assert.Equal(t, expectedUserModel, actualUserModel, "models mismatched")
	assert.Equal(t, testingSecret, actualSecret, "secret mismatched")
	assert.Equal(t, err, nil)
}

func TestLoginUserNotFound(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}

	usersMock.EXPECT().Authorize(testingEmail, testingPassword).Return(0, handlers.ErrAuthUserNotFound)
	_, _, err := testingUseCase.UserLogin(testUserModel)

	assert.Equal(t, err, handlers.ErrAuthUserNotFound)
}

func TestLoginSessionNotAdded(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}

	usersMock.EXPECT().Authorize(testingEmail, testingPassword).Return(testingID, nil)
	sessionsMock.EXPECT().AddSession(testingID, "").Return("", handlers.ErrAuthSessionNotFound)
	_, _, err := testingUseCase.UserLogin(testUserModel)

	assert.Equal(t, err, handlers.ErrAuthSessionNotFound)
}

func TestRegisterRegisterOk(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}
	expectedUserModel := models.UserID{ID: testingID}

	usersMock.EXPECT().AddUser(testingEmail, testingPassword).Return(testingID, nil)
	profileMock.EXPECT().AddEmptyProfile(testingID).Return(nil)
	sessionsMock.EXPECT().AddSession(testingID, "").Return(testingSecret, nil)
	actualUserModel, actualSecret, err := testingUseCase.UserRegister(testUserModel)

	assert.Equal(t, expectedUserModel, actualUserModel, "models mismatched")
	assert.Equal(t, testingSecret, actualSecret, "secret mismatched")
	assert.Equal(t, err, nil)
}

func TestRegisterEmailConflict(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}

	usersMock.EXPECT().AddUser(testingEmail, testingPassword).Return(0, handlers.ErrAuthEmailUsed)
	_, _, err := testingUseCase.UserRegister(testUserModel)

	assert.Equal(t, err, handlers.ErrAuthEmailUsed)
}

func TestRegisterProfileError(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}

	usersMock.EXPECT().AddUser(testingEmail, testingPassword).Return(testingID, nil)
	profileMock.EXPECT().AddEmptyProfile(testingID).Return(handlers.ErrProfileNotFiled)
	_, _, err := testingUseCase.UserRegister(testUserModel)

	assert.Equal(t, err, handlers.ErrProfileNotFiled)
}

func TestRegisterSessionError(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}

	usersMock.EXPECT().AddUser(testingEmail, testingPassword).Return(testingID, nil)
	profileMock.EXPECT().AddEmptyProfile(testingID).Return(nil)
	sessionsMock.EXPECT().AddSession(testingID, "").Return("", handlers.ErrAuthSessionNotFound)
	_, _, err := testingUseCase.UserRegister(testUserModel)

	assert.Equal(t, err, handlers.ErrAuthSessionNotFound)
}

func TestLogoutLogoutOk(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingSecret = "ifiewhufhbbjwdbnfnmwe"

	sessionsMock.EXPECT().RemoveSession(testingSecret).Return(nil)
	err := testingUseCase.UserLogout(testingSecret)

	assert.Equal(t, err, nil)
}

func TestChangePswdOk(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingOldPassword = "SomeSecret21"
	const testingNewPassword = "SomeSecret82y3"
	testUserModel := models.UserAuthProfile{Email: testingEmail, OldPassword: testingOldPassword, NewPassword: testingNewPassword}

	sessionsMock.EXPECT().GetIdBySession(testingSecret).Return(testingID, "", nil)
	usersMock.EXPECT().Authorize(testingEmail, testingOldPassword).Return(testingID, nil)
	usersMock.EXPECT().ChangePassword(testingID, testingNewPassword).Return(nil)
	err := testingUseCase.UserChangePassword(testUserModel, testingSecret)

	assert.Equal(t, err, nil)
}

func TestChangePswdWrongSession(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	const badId = 26
	const testingEmail = "test_email@ya.com"
	const testingOldPassword = "SomeSecret21"
	const testingNewPassword = "SomeSecret82y3"
	testUserModel := models.UserAuthProfile{Email: testingEmail, OldPassword: testingOldPassword, NewPassword: testingNewPassword}

	sessionsMock.EXPECT().GetIdBySession(testingSecret).Return(badId, "", nil)
	usersMock.EXPECT().Authorize(testingEmail, testingOldPassword).Return(testingID, nil)
	err := testingUseCase.UserChangePassword(testUserModel, testingSecret)

	assert.Equal(t, err, handlers.ErrAuthWrongPassword)
}

func TestChangePswdSecretError(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	const badId = 26
	const testingEmail = "test_email@ya.com"
	const testingOldPassword = "SomeSecret21"
	const testingNewPassword = "SomeSecret82y3"
	testUserModel := models.UserAuthProfile{Email: testingEmail, OldPassword: testingOldPassword, NewPassword: testingNewPassword}

	sessionsMock.EXPECT().GetIdBySession(testingSecret).Return(badId, "", handlers.ErrAuthSessionNotFound)
	err := testingUseCase.UserChangePassword(testUserModel, testingSecret)

	assert.Equal(t, err, handlers.ErrAuthSessionNotFound)
}

func TestChangePswdAuthError(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	usersMock := mocks.NewMockUsersRepository(mockController)
	sessionsMock := mocks.NewMockSessionsRepository(mockController)
	profileMock := mocks.NewMockProfileRepository(mockController)
	testingUseCase := NewAuthUseCaseImpl(usersMock, sessionsMock, profileMock)

	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	const badId = 26
	const testingEmail = "test_email@ya.com"
	const testingOldPassword = "SomeSecret21"
	const testingNewPassword = "SomeSecret82y3"
	testUserModel := models.UserAuthProfile{Email: testingEmail, OldPassword: testingOldPassword, NewPassword: testingNewPassword}

	sessionsMock.EXPECT().GetIdBySession(testingSecret).Return(testingID, "", nil)
	usersMock.EXPECT().Authorize(testingEmail, testingOldPassword).Return(testingID, handlers.ErrAuthWrongPassword)
	err := testingUseCase.UserChangePassword(testUserModel, testingSecret)

	assert.Equal(t, err, handlers.ErrAuthWrongPassword)
}
