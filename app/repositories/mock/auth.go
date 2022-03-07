package mock

import (
	"2022_1_OnlyGroup_back/pkg/errors"
	"math/rand"
)

const hashSize = 64

const secretRunes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

type userData struct {
	email    string
	password string
}

type AuthMock struct {
	userTable    []userData
	sessionTable map[string]int
}

func generateSecret(size int) string {
	result := ""
	for i := 0; i < size; i++ {
		result += string(secretRunes[rand.Intn(len(secretRunes))])
	}
	return result
}

func NewAuthMock() *AuthMock {
	return &AuthMock{sessionTable: make(map[string]int)}
}

func (tables *AuthMock) AddUser(email string, password string) (id int, err error) {
	for _, item := range tables.userTable {
		if item.email == email {
			return 0, errors.ErrAuthEmailUsed
		}
	}
	tables.userTable = append(tables.userTable, userData{email: email, password: password})
	return len(tables.userTable), nil
}

func (tables *AuthMock) Authorize(email string, password string) (id int, err error) {
	for index, item := range tables.userTable {
		if item.email == email {
			if item.password == password {
				return index, nil
			}
			return 0, errors.ErrAuthWrongPassword
		}
	}
	return 0, errors.ErrAuthUserNotFound
}

func (tables *AuthMock) ChangePassword(id int, newPassword string) (err error) {
	if id > len(tables.userTable)-1 {
		return errors.ErrAuthUserNotFound
	}
	tables.userTable[id].password = newPassword
	return nil
}

func (tables *AuthMock) AddSession(id int) (secret string, err error) {
	secret = generateSecret(hashSize)

	tables.sessionTable[secret] = id
	return
}

func (tables *AuthMock) GetIdBySession(secret string) (id int, err error) {
	id, has := tables.sessionTable[secret]
	if !has {
		return 0, errors.ErrAuthSessionNotFound
	}
	return id, nil
}

func (tables *AuthMock) RemoveSession(secret string) (err error) {
	_, has := tables.sessionTable[secret]
	if !has {
		return errors.ErrAuthSessionNotFound
	}
	delete(tables.sessionTable, secret)
	return nil
}
