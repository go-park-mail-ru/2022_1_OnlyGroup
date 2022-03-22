package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	mockUseCases "2022_1_OnlyGroup_back/app/tests/mocks"
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const url = "http://localhost/user"

func TestAuthAuthOk(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const testSessionSecret = "edfjiwehfbwbwewe"
	var testUserModel = models.UserID{ID: 3}
	var expectedResponse, _ = json.Marshal(testUserModel)
	var expectedCode = http.StatusOK
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserAuth(testSessionSecret).Return(testUserModel, nil)

	testingHandler := CreateAuthHandler(useCaseMock)
	req := httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session="+testSessionSecret)

	w := httptest.NewRecorder()
	testingHandler.GET(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
	assert.Equal(t, w.Body.String(), string(expectedResponse), "body mismatched, expected '%s', got '%s'",
		string(expectedResponse), w.Body.String())
}

func TestAuthAuthNoAuth(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const testSessionSecret = "edfjiwehfbwbwewe"
	var testUserModel = models.UserID{ID: 0}
	var expectedCode = http.StatusUnauthorized
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserAuth(testSessionSecret).Return(testUserModel, ErrAuthSessionNotFound)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session="+testSessionSecret)

	w := httptest.NewRecorder()
	testingHandler.GET(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
}

func TestAuthAuthNoCookie(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var expectedCode = http.StatusUnauthorized
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	testingHandler.GET(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
}

func TestLoginLoginOk(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const testSessionSecret = "edfjiwehfbwbwewe"
	var testRequestModel = models.UserAuthInfo{Email: "test.email@mail.corp.ru", Password: "Test_pass123"}
	var testRequestBody, _ = json.Marshal(testRequestModel)

	var testUserModel = models.UserID{ID: 3}
	var expectedResponse, _ = json.Marshal(testUserModel)
	var expectedCode = http.StatusOK
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserLogin(testRequestModel).Return(testUserModel, testSessionSecret, nil)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("PUT", url, bytes.NewReader(testRequestBody))

	w := httptest.NewRecorder()
	testingHandler.PUT(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
	assert.Equal(t, w.Body.String(), string(expectedResponse), "body mismatched, expected '%s', got '%s'",
		string(expectedResponse), w.Body.String())
	assert.Equal(t, strings.HasPrefix(w.HeaderMap.Get("Set-Cookie"), "session="+testSessionSecret), true,
		"session mismatched, expected %s, got %s", "session="+testSessionSecret, w.HeaderMap.Get("Set-Cookie"))
}

func TestLoginUserNotFound(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var testRequestModel = models.UserAuthInfo{Email: "test_email@test.ru", Password: "Test_pass123"}
	var testRequestBody, _ = json.Marshal(testRequestModel)

	var testUserModel = models.UserID{ID: 0}
	var expectedCode = http.StatusUnauthorized
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserLogin(testRequestModel).Return(testUserModel, "", ErrAuthUserNotFound)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("PUT", url, bytes.NewReader(testRequestBody))

	w := httptest.NewRecorder()
	testingHandler.PUT(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
}

func TestLoginBadRequest(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var testRequestBody = "{bad_json}"

	var expectedCode = http.StatusBadRequest
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("PUT", url, bytes.NewReader([]byte(testRequestBody)))

	w := httptest.NewRecorder()
	testingHandler.PUT(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
}

func TestLoginBadEmail(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var testRequestModel = models.UserAuthInfo{Email: "test_email", Password: "Test_pass123"}
	var testRequestBody, _ = json.Marshal(testRequestModel)

	var expectedCode = http.StatusPreconditionFailed
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("PUT", url, bytes.NewReader(testRequestBody))

	w := httptest.NewRecorder()
	testingHandler.PUT(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
}

func TestLogoutLogoutOk(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	const testSessionSecret = "edfjiwehfbwbwewe"

	var expectedCode = http.StatusOK
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserLogout(testSessionSecret).Return(nil)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("DELETE", url, nil)
	req.Header.Add("Cookie", "session="+testSessionSecret)

	w := httptest.NewRecorder()
	testingHandler.DELETE(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
	assert.Equal(t, strings.HasPrefix(w.HeaderMap.Get("Set-Cookie"), "session="+testSessionSecret), true,
		"session mismatched, expected %s, got %s", "session="+testSessionSecret, w.HeaderMap.Get("Set-Cookie"))
}

func TestLogoutSessionNotFound(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	const testSessionSecret = "edfjiwehfbwbwewe"

	var expectedCode = http.StatusUnauthorized
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserLogout(testSessionSecret).Return(ErrAuthSessionNotFound)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("DELETE", url, nil)
	req.Header.Add("Cookie", "session="+testSessionSecret)

	w := httptest.NewRecorder()
	testingHandler.DELETE(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
}

func TestLogoutNoCookie(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var expectedCode = http.StatusUnauthorized
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("DELETE", url, nil)

	w := httptest.NewRecorder()
	testingHandler.DELETE(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
}

func TestLogupLogupOk(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const testSessionSecret = "edfjiwehfbwbwewe"
	var testRequestModel = models.UserAuthInfo{Email: "test_email@test.ru", Password: "Test_pass123"}
	var testRequestBody, _ = json.Marshal(testRequestModel)

	var testUserModel = models.UserID{ID: 3}
	var expectedResponse, _ = json.Marshal(testUserModel)
	var expectedCode = http.StatusOK
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserRegister(testRequestModel).Return(testUserModel, testSessionSecret, nil)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("POST", url, bytes.NewReader(testRequestBody))

	w := httptest.NewRecorder()
	testingHandler.POST(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
	assert.Equal(t, w.Body.String(), string(expectedResponse), "body mismatched, expected '%s', got '%s'",
		string(expectedResponse), w.Body.String())
	assert.Equal(t, strings.HasPrefix(w.HeaderMap.Get("Set-Cookie"), "session="+testSessionSecret), true,
		"session mismatched, expected %s, got %s", "session="+testSessionSecret, w.HeaderMap.Get("Set-Cookie"))
}

func TestLogupUserConflict(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var testRequestModel = models.UserAuthInfo{Email: "test_email@test.ru", Password: "Test_pass123"}
	var testRequestBody, _ = json.Marshal(testRequestModel)

	var testUserModel = models.UserID{ID: 0}
	var expectedCode = http.StatusConflict
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserRegister(testRequestModel).Return(testUserModel, "", ErrAuthEmailUsed)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("POST", url, bytes.NewReader(testRequestBody))

	w := httptest.NewRecorder()
	testingHandler.POST(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
}

func TestLogupBadEmail(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var testRequestModel = models.UserAuthInfo{Email: "test_email", Password: "Test_pass123"}
	var testRequestBody, _ = json.Marshal(testRequestModel)

	var expectedCode = http.StatusPreconditionFailed
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("POST", url, bytes.NewReader(testRequestBody))

	w := httptest.NewRecorder()
	testingHandler.POST(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
}

func TestLogupBadRequest(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var testRequestBody = "{bad_json}"

	var expectedCode = http.StatusBadRequest
	useCaseMock := mockUseCases.NewMockAuthUseCases(mockController)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("POST", url, bytes.NewReader([]byte(testRequestBody)))

	w := httptest.NewRecorder()
	testingHandler.POST(w, req)

	assert.Equal(t, expectedCode, w.Code, "status code error, expected %d, got %d", expectedCode, w.Code)
}

func TestUserModelValidationTableDriven(t *testing.T) {
	var tests = []struct {
		email          string
		password       string
		expectedResult error
	}{{"test.email@corp.mail.ru", "TEst_pass123", nil},
		{"test.email@corp.mail.ru", "Len1", ErrAuthValidationPassword},
		{"test.email@corp.mail.ru", "len_max_test_kenNCVJNECJNNJCBdY487374367dgydghVGvdgdvfgevfyefhdvbfhevfvd", ErrAuthValidationPassword},
		{"test.email@corp.mail.ru", "no_upper_character38247834", ErrAuthValidationPassword},
		{"test.email@corp.mail.ru", "NO_LOWER_CHARACTER837463", ErrAuthValidationPassword},
		{"test.email@corp.mail.ru", "NO_number", ErrAuthValidationPassword},
		{"some_bad_email", "TEst_pass123", ErrAuthValidationEmail},
		{"some_bad@email", "TEst_pass123", ErrAuthValidationEmail},
		{"some_bad@email.", "TEst_pass123", ErrAuthValidationEmail},
		{"@bad.email.ru", "TEst_pass123", ErrAuthValidationEmail},
	}

	for _, testCase := range tests {
		testModel := models.UserAuthInfo{Email: testCase.email, Password: testCase.password}
		assert.Equal(t, testCase.expectedResult, checkValidUserModel(testModel), "email:'%s', pass:'%s'",
			testCase.email, testCase.password)
	}
}
