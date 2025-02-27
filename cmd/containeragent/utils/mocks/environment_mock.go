// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/cmd/containeragent/utils (interfaces: Environment)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockEnvironment is a mock of Environment interface
type MockEnvironment struct {
	ctrl     *gomock.Controller
	recorder *MockEnvironmentMockRecorder
}

// MockEnvironmentMockRecorder is the mock recorder for MockEnvironment
type MockEnvironmentMockRecorder struct {
	mock *MockEnvironment
}

// NewMockEnvironment creates a new mock instance
func NewMockEnvironment(ctrl *gomock.Controller) *MockEnvironment {
	mock := &MockEnvironment{ctrl: ctrl}
	mock.recorder = &MockEnvironmentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEnvironment) EXPECT() *MockEnvironmentMockRecorder {
	return m.recorder
}

// ExpandEnv mocks base method
func (m *MockEnvironment) ExpandEnv(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExpandEnv", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// ExpandEnv indicates an expected call of ExpandEnv
func (mr *MockEnvironmentMockRecorder) ExpandEnv(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExpandEnv", reflect.TypeOf((*MockEnvironment)(nil).ExpandEnv), arg0)
}

// Getenv mocks base method
func (m *MockEnvironment) Getenv(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Getenv", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Getenv indicates an expected call of Getenv
func (mr *MockEnvironmentMockRecorder) Getenv(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Getenv", reflect.TypeOf((*MockEnvironment)(nil).Getenv), arg0)
}

// Setenv mocks base method
func (m *MockEnvironment) Setenv(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setenv", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setenv indicates an expected call of Setenv
func (mr *MockEnvironmentMockRecorder) Setenv(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setenv", reflect.TypeOf((*MockEnvironment)(nil).Setenv), arg0, arg1)
}

// Unsetenv mocks base method
func (m *MockEnvironment) Unsetenv(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unsetenv", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsetenv indicates an expected call of Unsetenv
func (mr *MockEnvironmentMockRecorder) Unsetenv(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsetenv", reflect.TypeOf((*MockEnvironment)(nil).Unsetenv), arg0)
}
