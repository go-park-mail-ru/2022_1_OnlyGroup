package mock

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

const numTest = 6
const defaultMockUser = 6

func TestGenerateSecret(t *testing.T) {
	assert := assert.New(t)
	for i := 0; i < numTest; i++ {
		assert.Equal(len(generateSecret(i)), i)
	}
}

func TestNewAuthMock(t *testing.T) {
	assert := assert.New(t)
	TestMock := NewAuthMock()
	for i := 1; i < numTest; i++ {
		len, err := TestMock.AddUser(faker.Email(), faker.Password())
		assert.Equal(len, defaultMockUser+i)
		assert.Equal(err, nil)
	}
}

func TestAuthorizeMock(t *testing.T) {
	assert := assert.New(t)
	TestMock := NewAuthMock()

	for i := 1; i < numTest; i++ {
		len, err := TestMock.Authorize(faker.Email(), faker.Password())
		assert.Equal(len, 0)
		assert.Equal(err, handlers.ErrAuthUserNotFound)
	}
}

func TestChangePasswordMock(t *testing.T) {
	assert := assert.New(t)
	TestMock := NewAuthMock()

	for i := 1; i < numTest; i++ {
		err := TestMock.ChangePassword(i, faker.Password())
		assert.Equal(err, nil)
	}
}

func TestAddSessionMock(t *testing.T) {
	assert := assert.New(t)
	TestMock := NewAuthMock()

	for i := 1; i < numTest; i++ {
		_, err := TestMock.AddSession(i)
		assert.Equal(err, nil)
	}
}

func TestGetIdBySessionMock(t *testing.T) {
	assert := assert.New(t)
	TestMock := NewAuthMock()
	var secretArray []string
	for i := 0; i < numTest; i++ {
		str, err := TestMock.AddSession(i)
		secretArray = append(secretArray, str)
		assert.Equal(err, nil)
	}
	for i := 0; i < numTest; i++ {
		str, err := TestMock.GetIdBySession(secretArray[i])
		assert.Equal(str, i)
		assert.Equal(err, nil)
	}
}

func TestRemoveSessionMock(t *testing.T) {
	assert := assert.New(t)
	TestMock := NewAuthMock()
	var secretArray []string
	for i := 0; i < numTest; i++ {
		str, err := TestMock.AddSession(i)
		secretArray = append(secretArray, str)
		assert.Equal(err, nil)
	}
	for i := 0; i < numTest; i++ {
		err := TestMock.RemoveSession(secretArray[i])
		assert.Equal(err, nil)
	}
}
