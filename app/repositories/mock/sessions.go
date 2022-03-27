package mock

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"math/rand"
)

type sessionData struct {
	userID         int
	additionalData string
}

type SessionsMock struct {
	sessionTable map[string]sessionData
}

func NewSessionsMock() *SessionsMock {
	return &SessionsMock{sessionTable: map[string]sessionData{}}
}

func generateSecret(size int) string {
	result := ""
	for i := 0; i < size; i++ {
		result += string(secretRunes[rand.Intn(len(secretRunes))])
	}
	return result
}

func (tables *SessionsMock) AddSession(id int, additionalData string) (string, error) {
	secret := generateSecret(hashSize)

	tables.sessionTable[secret] = sessionData{userID: id, additionalData: additionalData}
	return secret, nil
}

func (tables *SessionsMock) GetIdBySession(secret string) (int, string, error) {
	data, has := tables.sessionTable[secret]
	if !has {
		return 0, "", handlers.ErrAuthSessionNotFound
	}
	return data.userID, data.additionalData, nil
}

func (tables *SessionsMock) RemoveSession(secret string) (err error) {
	_, has := tables.sessionTable[secret]
	if !has {
		return handlers.ErrAuthSessionNotFound
	}
	delete(tables.sessionTable, secret)
	return nil
}
