package impl

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	"2022_1_OnlyGroup_back/app/models"
	"2022_1_OnlyGroup_back/app/tests/mocks"
	"2022_1_OnlyGroup_back/app/usecases"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type testSuite struct {
	suite.Suite
	usersMock      *mocks.MockUsersRepository
	sessionsMock   *mocks.MockSessionsRepository
	profileMock    *mocks.MockProfileRepository
	mockController *gomock.Controller
	testingUseCase usecases.AuthUseCases
}

func (suite *testSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.usersMock = mocks.NewMockUsersRepository(suite.mockController)
	suite.sessionsMock = mocks.NewMockSessionsRepository(suite.mockController)
	suite.profileMock = mocks.NewMockProfileRepository(suite.mockController)
	suite.testingUseCase = NewAuthUseCaseImpl(suite.usersMock, suite.sessionsMock, suite.profileMock)
}

func (suite *testSuite) AfterTest() {
	suite.mockController.Finish()
}

func (suite *testSuite) TestAuthAuthOk() {
	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	expectedUserModel := models.UserID{ID: testingID}

	suite.sessionsMock.EXPECT().Get(testingSecret).Return(testingID, "", nil)
	actualUserModel, err := suite.testingUseCase.UserAuth(testingSecret)

	assert.Equal(suite.T(), expectedUserModel, actualUserModel, "models mismatched")
	assert.Equal(suite.T(), err, nil)
}

func (suite *testSuite) TestAuthSessionNotFound() {
	const testingSecret = "ifiewhufhbbjwdbnfnmwe"

	suite.sessionsMock.EXPECT().Get(testingSecret).Return(0, "", http.ErrAuthSessionNotFound)
	_, err := suite.testingUseCase.UserAuth(testingSecret)

	assert.Equal(suite.T(), err, http.ErrAuthSessionNotFound)
}

func (suite *testSuite) TestLoginLoginOk() {
	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}
	expectedUserModel := models.UserID{ID: testingID}

	suite.usersMock.EXPECT().Authorize(testingEmail, testingPassword).Return(testingID, nil)
	suite.sessionsMock.EXPECT().Add(testingID, "").Return(testingSecret, nil)
	actualUserModel, actualSecret, err := suite.testingUseCase.UserLogin(testUserModel)

	assert.Equal(suite.T(), expectedUserModel, actualUserModel, "models mismatched")
	assert.Equal(suite.T(), testingSecret, actualSecret, "secret mismatched")
	assert.Equal(suite.T(), err, nil)
}

func (suite *testSuite) TestLoginUserNotFound() {
	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}

	suite.usersMock.EXPECT().Authorize(testingEmail, testingPassword).Return(0, http.ErrAuthUserNotFound)
	_, _, err := suite.testingUseCase.UserLogin(testUserModel)

	assert.Equal(suite.T(), err, http.ErrAuthUserNotFound)
}

func (suite *testSuite) TestLoginSessionNotAdded() {
	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}

	suite.usersMock.EXPECT().Authorize(testingEmail, testingPassword).Return(testingID, nil)
	suite.sessionsMock.EXPECT().Add(testingID, "").Return("", http.ErrAuthSessionNotFound)
	_, _, err := suite.testingUseCase.UserLogin(testUserModel)

	assert.Equal(suite.T(), err, http.ErrAuthSessionNotFound)
}

func (suite *testSuite) TestRegisterRegisterOk() {
	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}
	expectedUserModel := models.UserID{ID: testingID}

	suite.usersMock.EXPECT().AddUser(testingEmail, testingPassword).Return(testingID, nil)
	suite.profileMock.EXPECT().AddEmptyProfile(testingID).Return(nil)
	suite.sessionsMock.EXPECT().Add(testingID, "").Return(testingSecret, nil)
	actualUserModel, actualSecret, err := suite.testingUseCase.UserRegister(testUserModel)

	assert.Equal(suite.T(), expectedUserModel, actualUserModel, "models mismatched")
	assert.Equal(suite.T(), testingSecret, actualSecret, "secret mismatched")
	assert.Equal(suite.T(), err, nil)
}

func (suite *testSuite) TestRegisterEmailConflict() {
	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}

	suite.usersMock.EXPECT().AddUser(testingEmail, testingPassword).Return(0, http.ErrAuthEmailUsed)
	_, _, err := suite.testingUseCase.UserRegister(testUserModel)

	assert.Equal(suite.T(), err, http.ErrAuthEmailUsed)
}

func (suite *testSuite) TestRegisterProfileError() {
	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}

	suite.usersMock.EXPECT().AddUser(testingEmail, testingPassword).Return(testingID, nil)
	suite.profileMock.EXPECT().AddEmptyProfile(testingID).Return(http.ErrProfileNotFiled)
	_, _, err := suite.testingUseCase.UserRegister(testUserModel)

	assert.Equal(suite.T(), err, http.ErrProfileNotFiled)
}

func (suite *testSuite) TestRegisterSessionError() {
	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingPassword = "SomeSecret21"
	testUserModel := models.UserAuthInfo{Email: testingEmail, Password: testingPassword}

	suite.usersMock.EXPECT().AddUser(testingEmail, testingPassword).Return(testingID, nil)
	suite.profileMock.EXPECT().AddEmptyProfile(testingID).Return(nil)
	suite.sessionsMock.EXPECT().Add(testingID, "").Return("", http.ErrAuthSessionNotFound)
	_, _, err := suite.testingUseCase.UserRegister(testUserModel)

	assert.Equal(suite.T(), err, http.ErrAuthSessionNotFound)
}

func (suite *testSuite) TestLogoutLogoutOk() {
	const testingSecret = "ifiewhufhbbjwdbnfnmwe"

	suite.sessionsMock.EXPECT().Remove(testingSecret).Return(nil)
	err := suite.testingUseCase.UserLogout(testingSecret)

	assert.Equal(suite.T(), err, nil)
}

func (suite *testSuite) TestChangePswdOk() {
	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingOldPassword = "SomeSecret21"
	const testingNewPassword = "SomeSecret82y3"
	testUserModel := models.UserAuthProfile{Email: testingEmail, OldPassword: testingOldPassword, NewPassword: testingNewPassword}

	suite.sessionsMock.EXPECT().Get(testingSecret).Return(testingID, "", nil)
	suite.usersMock.EXPECT().Authorize(testingEmail, testingOldPassword).Return(testingID, nil)
	suite.usersMock.EXPECT().ChangePassword(testingID, testingNewPassword).Return(nil)
	err := suite.testingUseCase.UserChangePassword(testUserModel, testingSecret)

	assert.Equal(suite.T(), err, nil)
}

func (suite *testSuite) TestChangePswdWrongSession() {
	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	const badId = 26
	const testingEmail = "test_email@ya.com"
	const testingOldPassword = "SomeSecret21"
	const testingNewPassword = "SomeSecret82y3"
	testUserModel := models.UserAuthProfile{Email: testingEmail, OldPassword: testingOldPassword, NewPassword: testingNewPassword}

	suite.sessionsMock.EXPECT().Get(testingSecret).Return(badId, "", nil)
	suite.usersMock.EXPECT().Authorize(testingEmail, testingOldPassword).Return(testingID, nil)
	err := suite.testingUseCase.UserChangePassword(testUserModel, testingSecret)

	assert.Equal(suite.T(), err, http.ErrAuthWrongPassword)
}

func (suite *testSuite) TestChangePswdSecretError() {
	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const badId = 26
	const testingEmail = "test_email@ya.com"
	const testingOldPassword = "SomeSecret21"
	const testingNewPassword = "SomeSecret82y3"
	testUserModel := models.UserAuthProfile{Email: testingEmail, OldPassword: testingOldPassword, NewPassword: testingNewPassword}

	suite.sessionsMock.EXPECT().Get(testingSecret).Return(badId, "", http.ErrAuthSessionNotFound)
	err := suite.testingUseCase.UserChangePassword(testUserModel, testingSecret)

	assert.Equal(suite.T(), err, http.ErrAuthSessionNotFound)
}

func (suite *testSuite) TestChangePswdAuthError() {
	const testingSecret = "ifiewhufhbbjwdbnfnmwe"
	const testingID = 24
	const testingEmail = "test_email@ya.com"
	const testingOldPassword = "SomeSecret21"
	const testingNewPassword = "SomeSecret82y3"
	testUserModel := models.UserAuthProfile{Email: testingEmail, OldPassword: testingOldPassword, NewPassword: testingNewPassword}

	suite.sessionsMock.EXPECT().Get(testingSecret).Return(testingID, "", nil)
	suite.usersMock.EXPECT().Authorize(testingEmail, testingOldPassword).Return(testingID, http.ErrAuthWrongPassword)
	err := suite.testingUseCase.UserChangePassword(testUserModel, testingSecret)

	assert.Equal(suite.T(), err, http.ErrAuthWrongPassword)
}

func TestAll(t *testing.T) {
	suite.Run(t, new(testSuite))
}
