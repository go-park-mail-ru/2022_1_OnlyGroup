package tests

//const url = "/profiles/"
//
//func TestGetProfileOk(t *testing.T) {
//	mockController := gomock.NewController(t)
//	defer mockController.Finish()
//
//	const testSessionSecret = "edfjiwehfbwbwewe"
//
//	var testProfileModel = models.Profile{FirstName: faker.FirstName(), LastName: faker.LastName(), Birthday: faker.Date(), City: "Moscow", Interests: []string{faker.Word(), faker.Word()}, AboutUser: faker.Sentence(), UserId: 0, Gender: faker.Gender()}
//	var expectedResponse, _ = json.Marshal(testProfileModel)
//	var expectedCode = http.StatusOK
//	useCaseMock := mock_usecases.NewMockProfileUseCases(mockController)
//
//	useCaseMock.EXPECT().Get(testSessionSecret, 0).Return(testProfileModel, nil)
//
//	testingHandler := handlers.ProfileHandler{ProfileUseCase: useCaseMock}
//	req := httptest.NewRequest("GET", url, nil)
//
//	req.Header.Add("Cookie", "session="+testSessionSecret)
//	w := httptest.NewRecorder()
//
//	testingHandler.GetProfileHandler(w, req)
//
//	if w.Code != expectedCode {
//		t.Fatalf("status code error, expected %d, got %d", expectedCode, w.Code)
//	}
//	if w.Body.String() != string(expectedResponse) {
//		t.Fatalf("body mismatched, expected '%s', got '%s'", string(expectedResponse), w.Body.String())
//	}
//}
