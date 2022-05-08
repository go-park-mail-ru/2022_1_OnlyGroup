//// Code generated by MockGen. DO NOT EDIT.
//// Source: ./app/usecase/profile.go
//
//// Package mock is a generated GoMock package.
package mocks

import (
	models "2022_1_OnlyGroup_back/app/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockProfileUseCases is a mock of ProfileUseCases interface
type MockProfileUseCases struct {
	ctrl     *gomock.Controller
	recorder *MockProfileUseCasesMockRecorder
}

// MockProfileUseCasesMockRecorder is the mock recorder for MockProfileUseCases
type MockProfileUseCasesMockRecorder struct {
	mock *MockProfileUseCases
}

// NewMockProfileUseCases creates a new mock instance
func NewMockProfileUseCases(ctrl *gomock.Controller) *MockProfileUseCases {
	mock := &MockProfileUseCases{ctrl: ctrl}
	mock.recorder = &MockProfileUseCasesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProfileUseCases) EXPECT() *MockProfileUseCasesMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockProfileUseCases) Get(cookies string, candidateId int) (models.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", cookies, candidateId)
	ret0, _ := ret[0].(models.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockProfileUseCasesMockRecorder) Get(cookies, candidateId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockProfileUseCases)(nil).Get), cookies, candidateId)
}

// Change mocks base method
func (m *MockProfileUseCases) Change(cookies string, profile models.Profile) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Change", cookies, profile)
	ret0, _ := ret[0].(error)
	return ret0
}

// Change indicates an expected call of Change
func (mr *MockProfileUseCasesMockRecorder) Change(cookies, profile interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Change", reflect.TypeOf((*MockProfileUseCases)(nil).Change), cookies, profile)
}

// GetShort mocks base method
func (m *MockProfileUseCases) GetShort(cookies string, profileId int) (models.ShortProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShort", cookies, profileId)
	ret0, _ := ret[0].(models.ShortProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShort indicates an expected call of GetShort
func (mr *MockProfileUseCasesMockRecorder) GetShort(cookies, profileId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShort", reflect.TypeOf((*MockProfileUseCases)(nil).GetShort), cookies, profileId)
}

// Delete mocks base method
func (m *MockProfileUseCases) Delete(cookies string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", cookies)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockProfileUseCasesMockRecorder) Delete(cookies interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProfileUseCases)(nil).Delete), cookies)
}

// GetCandidates mocks base method
func (m *MockProfileUseCases) GetCandidates(cookies string) (models.VectorCandidate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCandidates", cookies)
	ret0, _ := ret[0].(models.VectorCandidate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCandidates indicates an expected call of GetCandidates
func (mr *MockProfileUseCasesMockRecorder) GetCandidates(cookies interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCandidates", reflect.TypeOf((*MockProfileUseCases)(nil).GetCandidates), cookies)
}
