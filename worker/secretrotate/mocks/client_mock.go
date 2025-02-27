// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/worker/secretrotate (interfaces: SecretManagerFacade)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	watcher "github.com/DavinZhang/juju/core/watcher"
)

// MockSecretManagerFacade is a mock of SecretManagerFacade interface.
type MockSecretManagerFacade struct {
	ctrl     *gomock.Controller
	recorder *MockSecretManagerFacadeMockRecorder
}

// MockSecretManagerFacadeMockRecorder is the mock recorder for MockSecretManagerFacade.
type MockSecretManagerFacadeMockRecorder struct {
	mock *MockSecretManagerFacade
}

// NewMockSecretManagerFacade creates a new mock instance.
func NewMockSecretManagerFacade(ctrl *gomock.Controller) *MockSecretManagerFacade {
	mock := &MockSecretManagerFacade{ctrl: ctrl}
	mock.recorder = &MockSecretManagerFacadeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretManagerFacade) EXPECT() *MockSecretManagerFacadeMockRecorder {
	return m.recorder
}

// WatchSecretsRotationChanges mocks base method.
func (m *MockSecretManagerFacade) WatchSecretsRotationChanges(arg0 string) (watcher.SecretRotationWatcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchSecretsRotationChanges", arg0)
	ret0, _ := ret[0].(watcher.SecretRotationWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchSecretsRotationChanges indicates an expected call of WatchSecretsRotationChanges.
func (mr *MockSecretManagerFacadeMockRecorder) WatchSecretsRotationChanges(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchSecretsRotationChanges", reflect.TypeOf((*MockSecretManagerFacade)(nil).WatchSecretsRotationChanges), arg0)
}
