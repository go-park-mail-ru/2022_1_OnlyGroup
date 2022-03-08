package handlers

import (
	"2022_1_OnlyGroup_back/app/models"
	mock_usecases "2022_1_OnlyGroup_back/app/tests/mocks"
	"2022_1_OnlyGroup_back/pkg/errors"
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
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
	useCaseMock := mock_usecases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserAuth(testSessionSecret).Return(testUserModel, nil)

	testingHandler := CreateAuthHandler(useCaseMock)
	req := httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session="+testSessionSecret)

	w := httptest.NewRecorder()
	testingHandler.AuthUserHandler(w, req)

	if w.Code != expectedCode {
		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
	}
	if w.Body.String() != string(expectedResponse) {
		t.Fatalf("body mismatched, expected '%s', got '%s'", string(expectedResponse), w.Body.String())
	}
}

func TestAuthAuthNoAuth(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const testSessionSecret = "edfjiwehfbwbwewe"
	var testUserModel = models.UserID{ID: 0}
	var expectedCode = http.StatusUnauthorized
	useCaseMock := mock_usecases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserAuth(testSessionSecret).Return(testUserModel, errors.ErrAuthSessionNotFound)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session="+testSessionSecret)

	w := httptest.NewRecorder()
	testingHandler.AuthUserHandler(w, req)

	if w.Code != expectedCode {
		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
	}
}

func TestAuthAuthNoCookie(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var expectedCode = http.StatusUnauthorized
	useCaseMock := mock_usecases.NewMockAuthUseCases(mockController)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("GET", url, nil)

	w := httptest.NewRecorder()
	testingHandler.AuthUserHandler(w, req)

	if w.Code != expectedCode {
		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
	}
}

func TestLoginLoginOk(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const testSessionSecret = "edfjiwehfbwbwewe"
	var testRequestModel = models.UserAuthInfo{Email: "test_email", Password: "test_pass"}
	var testRequestBody, _ = json.Marshal(testRequestModel)

	var testUserModel = models.UserID{ID: 3}
	var expectedResponse, _ = json.Marshal(testUserModel)
	var expectedCode = http.StatusOK
	useCaseMock := mock_usecases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserLogin(testRequestModel).Return(testUserModel, testSessionSecret, nil)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("PUT", url, bytes.NewReader(testRequestBody))

	w := httptest.NewRecorder()
	testingHandler.LoginUserHandler(w, req)

	if w.Code != expectedCode {
		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
	}
	if w.Body.String() != string(expectedResponse) {
		t.Fatalf("body mismatched, expected '%s', got '%s'", string(expectedResponse), w.Body.String())
	}
	if !strings.HasPrefix(w.HeaderMap.Get("Set-Cookie"), "session="+testSessionSecret) {
		t.Fatalf("session mismatched, expected %s, got %s", "session="+testSessionSecret, w.HeaderMap.Get("Set-Cookie"))
	}
}

func TestLoginUserNotFound(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var testRequestModel = models.UserAuthInfo{Email: "test_email", Password: "test_pass"}
	var testRequestBody, _ = json.Marshal(testRequestModel)

	var testUserModel = models.UserID{ID: 0}
	var expectedCode = http.StatusUnauthorized
	useCaseMock := mock_usecases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserLogin(testRequestModel).Return(testUserModel, "", errors.ErrAuthUserNotFound)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("PUT", url, bytes.NewReader(testRequestBody))

	w := httptest.NewRecorder()
	testingHandler.LoginUserHandler(w, req)

	if w.Code != expectedCode {
		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
	}
}

func TestLoginBadRequest(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var testRequestBody = "{bad_json}"

	var expectedCode = http.StatusBadRequest
	useCaseMock := mock_usecases.NewMockAuthUseCases(mockController)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("PUT", url, bytes.NewReader([]byte(testRequestBody)))

	w := httptest.NewRecorder()
	testingHandler.LoginUserHandler(w, req)

	if w.Code != expectedCode {
		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
	}
}

func TestLogoutLogoutOk(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	const testSessionSecret = "edfjiwehfbwbwewe"

	var expectedCode = http.StatusOK
	useCaseMock := mock_usecases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserLogout(testSessionSecret).Return(nil)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("DELETE", url, nil)
	req.Header.Add("Cookie", "session="+testSessionSecret)

	w := httptest.NewRecorder()
	testingHandler.LogoutUserHandler(w, req)

	if w.Code != expectedCode {
		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
	}
}

func TestLogoutSessionNotFound(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	const testSessionSecret = "edfjiwehfbwbwewe"

	var expectedCode = http.StatusUnauthorized
	useCaseMock := mock_usecases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserLogout(testSessionSecret).Return(errors.ErrAuthSessionNotFound)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("DELETE", url, nil)
	req.Header.Add("Cookie", "session="+testSessionSecret)

	w := httptest.NewRecorder()
	testingHandler.LogoutUserHandler(w, req)

	if w.Code != expectedCode {
		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
	}
}

func TestLogoutNoCookie(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var expectedCode = http.StatusUnauthorized
	useCaseMock := mock_usecases.NewMockAuthUseCases(mockController)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("DELETE", url, nil)

	w := httptest.NewRecorder()
	testingHandler.LogoutUserHandler(w, req)

	if w.Code != expectedCode {
		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
	}
}

func TestLogupLogupOk(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	const testSessionSecret = "edfjiwehfbwbwewe"
	var testRequestModel = models.UserAuthInfo{Email: "test_email", Password: "test_pass"}
	var testRequestBody, _ = json.Marshal(testRequestModel)

	var testUserModel = models.UserID{ID: 3}
	var expectedResponse, _ = json.Marshal(testUserModel)
	var expectedCode = http.StatusOK
	useCaseMock := mock_usecases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserRegister(testRequestModel).Return(testUserModel, testSessionSecret, nil)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("POST", url, bytes.NewReader(testRequestBody))

	w := httptest.NewRecorder()
	testingHandler.LogupUserHandler(w, req)

	if w.Code != expectedCode {
		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
	}
	if w.Body.String() != string(expectedResponse) {
		t.Fatalf("body mismatched, expected '%s', got '%s'", string(expectedResponse), w.Body.String())
	}
	if !strings.HasPrefix(w.HeaderMap.Get("Set-Cookie"), "session="+testSessionSecret) {
		t.Fatalf("session mismatched, expected %s, got %s", "session="+testSessionSecret, w.HeaderMap.Get("Set-Cookie"))
	}
}

func TestLogupUserConflict(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var testRequestModel = models.UserAuthInfo{Email: "test_email", Password: "test_pass"}
	var testRequestBody, _ = json.Marshal(testRequestModel)

	var testUserModel = models.UserID{ID: 0}
	var expectedCode = http.StatusConflict
	useCaseMock := mock_usecases.NewMockAuthUseCases(mockController)

	useCaseMock.EXPECT().UserRegister(testRequestModel).Return(testUserModel, "", errors.ErrAuthEmailUsed)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("POST", url, bytes.NewReader(testRequestBody))

	w := httptest.NewRecorder()
	testingHandler.LogupUserHandler(w, req)

	if w.Code != expectedCode {
		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
	}
}

func TestLogupBadRequest(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	var testRequestBody = "{bad_json}"

	var expectedCode = http.StatusBadRequest
	useCaseMock := mock_usecases.NewMockAuthUseCases(mockController)

	testingHandler := AuthHandler{AuthUseCase: useCaseMock}
	req := httptest.NewRequest("POST", url, bytes.NewReader([]byte(testRequestBody)))

	w := httptest.NewRecorder()
	testingHandler.LogupUserHandler(w, req)

	if w.Code != expectedCode {
		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
	}
}
