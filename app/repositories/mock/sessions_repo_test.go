package mock

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddSessionMock(t *testing.T) {
	assert := assert.New(t)
	TestMock := NewSessionsMock()

	for i := 1; i < numTest; i++ {
		_, err := TestMock.AddSession(i, "")
		assert.Equal(err, nil)
	}
}

func TestGetIdBySessionMock(t *testing.T) {
	assert := assert.New(t)
	TestMock := NewSessionsMock()
	var secretArray []string
	for i := 0; i < numTest; i++ {
		str, err := TestMock.AddSession(i, "")
		secretArray = append(secretArray, str)
		assert.Equal(err, nil)
	}
	for i := 0; i < numTest; i++ {
		str, _, err := TestMock.GetIdBySession(secretArray[i])
		assert.Equal(str, i)
		assert.Equal(err, nil)
	}
}

func TestRemoveSessionMock(t *testing.T) {
	assert := assert.New(t)
	TestMock := NewSessionsMock()
	var secretArray []string
	for i := 0; i < numTest; i++ {
		str, err := TestMock.AddSession(i, "")
		secretArray = append(secretArray, str)
		assert.Equal(err, nil)
	}
	for i := 0; i < numTest; i++ {
		err := TestMock.RemoveSession(secretArray[i])
		assert.Equal(err, nil)
	}
}
