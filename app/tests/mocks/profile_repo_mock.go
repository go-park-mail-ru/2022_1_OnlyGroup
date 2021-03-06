// Code generated by MockGen. DO NOT EDIT.
// Source: app/repositories/profile.go

// Package mock_repositories is a generated GoMock package.
package mocks

import (
	models "2022_1_OnlyGroup_back/app/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockProfileRepository is a mock of ProfileRepository interface.
type MockProfileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProfileRepositoryMockRecorder
}

// MockProfileRepositoryMockRecorder is the mock recorder for MockProfileRepository.
type MockProfileRepositoryMockRecorder struct {
	mock *MockProfileRepository
}

// NewMockProfileRepository creates a new mock instance.
func NewMockProfileRepository(ctrl *gomock.Controller) *MockProfileRepository {
	mock := &MockProfileRepository{ctrl: ctrl}
	mock.recorder = &MockProfileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProfileRepository) EXPECT() *MockProfileRepositoryMockRecorder {
	return m.recorder
}

// AddEmptyProfile mocks base method.
func (m *MockProfileRepository) AddEmpty(profileId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEmpty", profileId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEmptyProfile indicates an expected call of AddEmptyProfile.
func (mr *MockProfileRepositoryMockRecorder) AddEmptyProfile(profileId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEmpty", reflect.TypeOf((*MockProfileRepository)(nil).AddEmpty), profileId)
}

// AddProfile mocks base method.
func (m *MockProfileRepository) Add(profile models.Profile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", profile)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddProfile indicates an expected call of AddProfile.
func (mr *MockProfileRepositoryMockRecorder) AddProfile(profile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockProfileRepository)(nil).Add), profile)
}

// ChangeProfile mocks base method.
func (m *MockProfileRepository) Change(profileId int, profile models.Profile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Change", profileId, profile)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeProfile indicates an expected call of ChangeProfile.
func (mr *MockProfileRepositoryMockRecorder) ChangeProfile(profileId, profile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Change", reflect.TypeOf((*MockProfileRepository)(nil).Change), profileId, profile)
}

// CheckProfileFiled mocks base method.
func (m *MockProfileRepository) CheckFiled(profileId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckFiled", profileId)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckProfileFiled indicates an expected call of CheckProfileFiled.
func (mr *MockProfileRepositoryMockRecorder) CheckProfileFiled(profileId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckFiled", reflect.TypeOf((*MockProfileRepository)(nil).CheckFiled), profileId)
}

// DeleteProfile mocks base method.
func (m *MockProfileRepository) Delete(profileId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", profileId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProfile indicates an expected call of DeleteProfile.
func (mr *MockProfileRepositoryMockRecorder) DeleteProfile(profileId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProfileRepository)(nil).Delete), profileId)
}

// FindCandidateProfile mocks base method.
func (m *MockProfileRepository) FindCandidate(profileId int) (models.VectorCandidate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindCandidate", profileId)
	ret0, _ := ret[0].(models.VectorCandidate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindCandidateProfile indicates an expected call of FindCandidateProfile.
func (mr *MockProfileRepositoryMockRecorder) FindCandidateProfile(profileId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindCandidate", reflect.TypeOf((*MockProfileRepository)(nil).FindCandidate), profileId)
}

// GetProfile mocks base method.
func (m *MockProfileRepository) Get(profileId int) (models.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", profileId)
	ret0, _ := ret[0].(models.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfile indicates an expected call of GetProfile.
func (mr *MockProfileRepositoryMockRecorder) GetProfile(profileId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockProfileRepository)(nil).Get), profileId)
}

// GetShortProfile mocks base method.
func (m *MockProfileRepository) GetShort(profileId int) (models.ShortProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShort", profileId)
	ret0, _ := ret[0].(models.ShortProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortProfile indicates an expected call of GetShortProfile.
func (mr *MockProfileRepositoryMockRecorder) GetShortProfile(profileId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShort", reflect.TypeOf((*MockProfileRepository)(nil).GetShort), profileId)
}
