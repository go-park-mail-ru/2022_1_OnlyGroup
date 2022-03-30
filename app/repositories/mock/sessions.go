package mock

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/pkg/sessionGenerator"
)

type sessionData struct {
	userID         int
	additionalData string
}

type SessionsMock struct {
	sessionTable    map[string]sessionData
	secretGenerator sessionGenerator.SessionGenerator
}

func NewSessionsMock() *SessionsMock {
	return &SessionsMock{sessionTable: map[string]sessionData{}, secretGenerator: sessionGenerator.NewRandomGenerator()}
}

func (tables *SessionsMock) Add(id int, additionalData string) (string, error) {
	secret := tables.secretGenerator.String(hashSize)

	tables.sessionTable[secret] = sessionData{userID: id, additionalData: additionalData}
	return secret, nil
}

func (tables *SessionsMock) Get(secret string) (int, string, error) {
	data, has := tables.sessionTable[secret]
	if !has {
		return 0, "", handlers.ErrAuthSessionNotFound
	}
	return data.userID, data.additionalData, nil
}

func (tables *SessionsMock) Remove(secret string) (err error) {
	_, has := tables.sessionTable[secret]
	if !has {
		return handlers.ErrAuthSessionNotFound
	}
	delete(tables.sessionTable, secret)
	return nil
}
