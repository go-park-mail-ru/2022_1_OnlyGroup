package mock

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	"2022_1_OnlyGroup_back/pkg/randomGenerator"
	impl3 "2022_1_OnlyGroup_back/pkg/randomGenerator/impl"
)

type sessionData struct {
	userID         int
	additionalData string
}

type SessionsMock struct {
	sessionTable    map[string]sessionData
	secretGenerator randomGenerator.RandomGenerator
}

func NewSessionsMock() *SessionsMock {
	return &SessionsMock{sessionTable: map[string]sessionData{}, secretGenerator: impl3.NewMathRandomGenerator()}
}

func (tables *SessionsMock) Add(id int, additionalData string) (string, error) {
	secret, err := tables.secretGenerator.String(hashSize)
	if err != nil {
		return "", err
	}

	tables.sessionTable[secret] = sessionData{userID: id, additionalData: additionalData}
	return secret, nil
}

func (tables *SessionsMock) Get(secret string) (int, string, error) {
	data, has := tables.sessionTable[secret]
	if !has {
		return 0, "", http.ErrAuthSessionNotFound
	}
	return data.userID, data.additionalData, nil
}

func (tables *SessionsMock) Remove(secret string) (err error) {
	_, has := tables.sessionTable[secret]
	if !has {
		return http.ErrAuthSessionNotFound
	}
	delete(tables.sessionTable, secret)
	return nil
}
