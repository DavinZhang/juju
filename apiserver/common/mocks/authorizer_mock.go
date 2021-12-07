// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/apiserver/common (interfaces: Authorizer)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	names "github.com/juju/names/v4"
	reflect "reflect"
)

// MockAuthorizer is a mock of Authorizer interface
type MockAuthorizer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizerMockRecorder
}

// MockAuthorizerMockRecorder is the mock recorder for MockAuthorizer
type MockAuthorizerMockRecorder struct {
	mock *MockAuthorizer
}

// NewMockAuthorizer creates a new mock instance
func NewMockAuthorizer(ctrl *gomock.Controller) *MockAuthorizer {
	mock := &MockAuthorizer{ctrl: ctrl}
	mock.recorder = &MockAuthorizerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthorizer) EXPECT() *MockAuthorizerMockRecorder {
	return m.recorder
}

// AuthController mocks base method
func (m *MockAuthorizer) AuthController() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthController")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthController indicates an expected call of AuthController
func (mr *MockAuthorizerMockRecorder) AuthController() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthController", reflect.TypeOf((*MockAuthorizer)(nil).AuthController))
}

// AuthMachineAgent mocks base method
func (m *MockAuthorizer) AuthMachineAgent() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthMachineAgent")
	ret0, _ := ret[0].(bool)
	return ret0
}

// AuthMachineAgent indicates an expected call of AuthMachineAgent
func (mr *MockAuthorizerMockRecorder) AuthMachineAgent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthMachineAgent", reflect.TypeOf((*MockAuthorizer)(nil).AuthMachineAgent))
}

// GetAuthTag mocks base method
func (m *MockAuthorizer) GetAuthTag() names.Tag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthTag")
	ret0, _ := ret[0].(names.Tag)
	return ret0
}

// GetAuthTag indicates an expected call of GetAuthTag
func (mr *MockAuthorizerMockRecorder) GetAuthTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthTag", reflect.TypeOf((*MockAuthorizer)(nil).GetAuthTag))
}
