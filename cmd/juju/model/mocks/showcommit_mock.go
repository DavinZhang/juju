// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/cmd/juju/model (interfaces: ShowCommitCommandAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/DavinZhang/juju/core/model"
)

// MockShowCommitCommandAPI is a mock of ShowCommitCommandAPI interface
type MockShowCommitCommandAPI struct {
	ctrl     *gomock.Controller
	recorder *MockShowCommitCommandAPIMockRecorder
}

// MockShowCommitCommandAPIMockRecorder is the mock recorder for MockShowCommitCommandAPI
type MockShowCommitCommandAPIMockRecorder struct {
	mock *MockShowCommitCommandAPI
}

// NewMockShowCommitCommandAPI creates a new mock instance
func NewMockShowCommitCommandAPI(ctrl *gomock.Controller) *MockShowCommitCommandAPI {
	mock := &MockShowCommitCommandAPI{ctrl: ctrl}
	mock.recorder = &MockShowCommitCommandAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockShowCommitCommandAPI) EXPECT() *MockShowCommitCommandAPIMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockShowCommitCommandAPI) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockShowCommitCommandAPIMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockShowCommitCommandAPI)(nil).Close))
}

// ShowCommit mocks base method
func (m *MockShowCommitCommandAPI) ShowCommit(arg0 int) (model.GenerationCommit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShowCommit", arg0)
	ret0, _ := ret[0].(model.GenerationCommit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShowCommit indicates an expected call of ShowCommit
func (mr *MockShowCommitCommandAPIMockRecorder) ShowCommit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShowCommit", reflect.TypeOf((*MockShowCommitCommandAPI)(nil).ShowCommit), arg0)
}
