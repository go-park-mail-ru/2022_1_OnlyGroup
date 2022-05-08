package mock

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	impl3 "2022_1_OnlyGroup_back/pkg/randomGenerator/impl"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

const numTest = 6
const defaultMockUser = 5

func TestGenerateSecret(t *testing.T) {
	generator := impl3.NewMathRandomGenerator()
	for i := 0; i < numTest; i++ {
		s, err := generator.String(i)
		assert.Equal(t, nil, err)
		assert.Equal(t, s, i)
	}
}

func TestNewAuthMock(t *testing.T) {
	assert := assert.New(t)
	TestMock := NewUsersMock()
	for i := 1; i < numTest; i++ {
		len, err := TestMock.AddUser(faker.Email(), faker.Password())
		assert.Equal(len, defaultMockUser+i)
		assert.Equal(err, nil)
	}
}

func TestAuthorizeMock(t *testing.T) {
	assert := assert.New(t)
	TestMock := NewUsersMock()

	for i := 1; i < numTest; i++ {
		len, err := TestMock.Authorize(faker.Email(), faker.Password())
		assert.Equal(len, 0)
		assert.Equal(err, http.ErrAuthUserNotFound)
	}
}

func TestChangePasswordMock(t *testing.T) {
	assert := assert.New(t)
	TestMock := NewUsersMock()

	for i := 1; i < numTest; i++ {
		err := TestMock.ChangePassword(i, faker.Password())
		assert.Equal(err, nil)
	}
}
