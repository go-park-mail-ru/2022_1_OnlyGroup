package mock

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"github.com/bxcodec/faker/v3"
)

const hashSize = 64

const secretRunes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

type userData struct {
	email    string
	password string
}

type UsersMock struct {
	userTable []userData
}

func NewUsersMock() *UsersMock {
	data := []userData{{email: "petrenko@mail.ru", password: "Qwerty1234"}}
	for i := 1; i < 6; i++ {
		data = append(data, userData{email: faker.Email(), password: faker.Password()})
	}
	return &UsersMock{userTable: data}
}

func (tables *UsersMock) AddUser(email string, password string) (id int, err error) {
	for _, item := range tables.userTable {
		if item.email == email {
			return 0, handlers.ErrAuthEmailUsed
		}
	}
	tables.userTable = append(tables.userTable, userData{email: email, password: password})
	return len(tables.userTable) - 1, nil
}

func (tables *UsersMock) Authorize(email string, password string) (id int, err error) {
	for index, item := range tables.userTable {
		if item.email == email {
			if item.password == password {
				return index, nil
			}
			return 0, handlers.ErrAuthWrongPassword
		}
	}
	return 0, handlers.ErrAuthUserNotFound
}

func (tables *UsersMock) ChangePassword(id int, newPassword string) (err error) {
	if id > len(tables.userTable)-1 {
		return handlers.ErrAuthUserNotFound
	}
	tables.userTable[id].password = newPassword
	return nil
}
