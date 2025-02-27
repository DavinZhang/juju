// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/worker/caasmodelconfigmanager (interfaces: Facade)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	controller "github.com/DavinZhang/juju/controller"
)

// MockFacade is a mock of Facade interface.
type MockFacade struct {
	ctrl     *gomock.Controller
	recorder *MockFacadeMockRecorder
}

// MockFacadeMockRecorder is the mock recorder for MockFacade.
type MockFacadeMockRecorder struct {
	mock *MockFacade
}

// NewMockFacade creates a new mock instance.
func NewMockFacade(ctrl *gomock.Controller) *MockFacade {
	mock := &MockFacade{ctrl: ctrl}
	mock.recorder = &MockFacadeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFacade) EXPECT() *MockFacadeMockRecorder {
	return m.recorder
}

// ControllerConfig mocks base method.
func (m *MockFacade) ControllerConfig() (controller.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerConfig")
	ret0, _ := ret[0].(controller.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerConfig indicates an expected call of ControllerConfig.
func (mr *MockFacadeMockRecorder) ControllerConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerConfig", reflect.TypeOf((*MockFacade)(nil).ControllerConfig))
}
