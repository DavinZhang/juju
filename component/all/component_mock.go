// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/component/all (interfaces: Component)

// Package all is a generated GoMock package.
package all

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockComponent is a mock of Component interface
type MockComponent struct {
	ctrl     *gomock.Controller
	recorder *MockComponentMockRecorder
}

// MockComponentMockRecorder is the mock recorder for MockComponent
type MockComponentMockRecorder struct {
	mock *MockComponent
}

// NewMockComponent creates a new mock instance
func NewMockComponent(ctrl *gomock.Controller) *MockComponent {
	mock := &MockComponent{ctrl: ctrl}
	mock.recorder = &MockComponentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockComponent) EXPECT() *MockComponentMockRecorder {
	return m.recorder
}

// registerForClient mocks base method
func (m *MockComponent) registerForClient() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "registerForClient")
	ret0, _ := ret[0].(error)
	return ret0
}

// registerForClient indicates an expected call of registerForClient
func (mr *MockComponentMockRecorder) registerForClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "registerForClient", reflect.TypeOf((*MockComponent)(nil).registerForClient))
}

// registerForContainerAgent mocks base method
func (m *MockComponent) registerForContainerAgent() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "registerForContainerAgent")
	ret0, _ := ret[0].(error)
	return ret0
}

// registerForContainerAgent indicates an expected call of registerForContainerAgent
func (mr *MockComponentMockRecorder) registerForContainerAgent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "registerForContainerAgent", reflect.TypeOf((*MockComponent)(nil).registerForContainerAgent))
}

// registerForServer mocks base method
func (m *MockComponent) registerForServer() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "registerForServer")
	ret0, _ := ret[0].(error)
	return ret0
}

// registerForServer indicates an expected call of registerForServer
func (mr *MockComponentMockRecorder) registerForServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "registerForServer", reflect.TypeOf((*MockComponent)(nil).registerForServer))
}
