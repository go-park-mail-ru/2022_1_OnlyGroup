package redis

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/pkg/sessionGenerator"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"testing"
)

var redisError = errors.New("test redis err")

func TestAddSessionTableDriven(t *testing.T) {
	tests := []struct {
		testName        string
		mockPrepare     func(mock redismock.ClientMock)
		id              int
		generatedSecret string
		additionalData  string
		expectedError   error
	}{
		{
			"All ok",
			func(mock redismock.ClientMock) {
				mock.ExpectHSet("test_5", "5_fiiewifjwiefjwe", "daw").SetVal(1)
			},
			5,
			"fiiewifjwiefjwe",
			"daw",
			nil,
		},
		{
			"Redis internal error",
			func(mock redismock.ClientMock) {
				mock.ExpectHSet("test_5", "5_fiiewifjwiefjwe", "daw").SetErr(redisError)
			},
			5,
			"fiiewifjwiefjwe",
			"daw",
			handlers.ErrBaseApp,
		},
		{
			"Redis adding error",
			func(mock redismock.ClientMock) {
				mock.ExpectHSet("test_5", "5_fiiewifjwiefjwe", "daw").SetVal(0)
			},
			5,
			"fiiewifjwiefjwe",
			"daw",
			handlers.ErrBaseApp,
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			client, mock := redismock.NewClientMock()
			test.mockPrepare(mock)
			testingRepo := NewRedisSessionRepository(client, "test", sessionGenerator.NewRandomGenerator())
			_, err := testingRepo.addSessionInternal(test.id, test.additionalData, test.generatedSecret)
			err1 := mock.ExpectationsWereMet()
			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrapped error mismatched, expected: '%v', got '%v'", test.expectedError, err)
			}
			assert.Equal(t, nil, err1)
		})
	}
}

func TestGetSessionTableDriven(t *testing.T) {
	tests := []struct {
		testName               string
		mockPrepare            func(mock redismock.ClientMock)
		secret                 string
		expectedError          error
		expectedAdditionalData string
		expectedId             int
	}{
		{
			"All ok",
			func(mock redismock.ClientMock) {
				mock.ExpectHGet("test_5", "5_fiiewifjwiefjwe").SetVal("daw")
			},
			"5_fiiewifjwiefjwe",
			nil,
			"daw",
			5,
		},
		{
			"Session not found",
			func(mock redismock.ClientMock) {
				mock.ExpectHGet("test_5", "5_fiiewifjwiefjwe").SetErr(redis.Nil)
			},
			"5_fiiewifjwiefjwe",
			handlers.ErrAuthSessionNotFound,
			"",
			0,
		},
		{
			"Redis internal error",
			func(mock redismock.ClientMock) {
				mock.ExpectHGet("test_5", "5_fiiewifjwiefjwe").SetErr(redisError)
			},
			"5_fiiewifjwiefjwe",
			handlers.ErrBaseApp,
			"",
			0,
		},
		{
			"Bad session",
			func(mock redismock.ClientMock) {},
			"fiiewifjwiefjwe",
			handlers.ErrAuthSessionNotFound,
			"",
			0,
		},
		{
			"Bad session",
			func(mock redismock.ClientMock) {},
			"afadfs_fiiewifjwiefjwe",
			handlers.ErrAuthSessionNotFound,
			"",
			0,
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			client, mock := redismock.NewClientMock()
			test.mockPrepare(mock)
			testingRepo := NewRedisSessionRepository(client, "test", sessionGenerator.NewRandomGenerator())
			id, data, err := testingRepo.Get(test.secret)
			err1 := mock.ExpectationsWereMet()
			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrapped error mismatched, expected: '%v', got '%v'", test.expectedError, err)
			}
			assert.Equal(t, nil, err1)
			assert.Equal(t, test.expectedId, id)
			assert.Equal(t, test.expectedAdditionalData, data)
		})
	}
}

func TestRemoveSessionTableDriven(t *testing.T) {
	tests := []struct {
		testName      string
		mockPrepare   func(mock redismock.ClientMock)
		secret        string
		expectedError error
	}{
		{
			"All ok",
			func(mock redismock.ClientMock) {
				mock.ExpectHDel("test_5", "5_fiiewifjwiefjwe").SetVal(1)
			},
			"5_fiiewifjwiefjwe",
			nil,
		},
		{
			"Session not found",
			func(mock redismock.ClientMock) {
				mock.ExpectHDel("test_5", "5_fiiewifjwiefjwe").SetVal(0)
			},
			"5_fiiewifjwiefjwe",
			handlers.ErrAuthSessionNotFound,
		},
		{
			"Bad session",
			func(mock redismock.ClientMock) {},
			"fiiewifjwiefjwe",
			handlers.ErrAuthSessionNotFound,
		},
		{
			"Redis internal error",
			func(mock redismock.ClientMock) {
				mock.ExpectHDel("test_5", "5_fiiewifjwiefjwe").SetErr(redisError)
			},
			"5_fiiewifjwiefjwe",
			handlers.ErrBaseApp,
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			client, mock := redismock.NewClientMock()
			test.mockPrepare(mock)
			testingRepo := NewRedisSessionRepository(client, "test", sessionGenerator.NewRandomGenerator())
			err := testingRepo.Remove(test.secret)
			err1 := mock.ExpectationsWereMet()
			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrapped error mismatched, expected: '%v', got '%v'", test.expectedError, err)
			}
			assert.Equal(t, nil, err1)
		})
	}
}
