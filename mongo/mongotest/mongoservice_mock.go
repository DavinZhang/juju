// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/mongo (interfaces: MongoSnapService)

// Package mongotest is a generated GoMock package.
package mongotest

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMongoSnapService is a mock of MongoSnapService interface.
type MockMongoSnapService struct {
	ctrl     *gomock.Controller
	recorder *MockMongoSnapServiceMockRecorder
}

// MockMongoSnapServiceMockRecorder is the mock recorder for MockMongoSnapService.
type MockMongoSnapServiceMockRecorder struct {
	mock *MockMongoSnapService
}

// NewMockMongoSnapService creates a new mock instance.
func NewMockMongoSnapService(ctrl *gomock.Controller) *MockMongoSnapService {
	mock := &MockMongoSnapService{ctrl: ctrl}
	mock.recorder = &MockMongoSnapServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMongoSnapService) EXPECT() *MockMongoSnapServiceMockRecorder {
	return m.recorder
}

// ConfigOverride mocks base method.
func (m *MockMongoSnapService) ConfigOverride() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfigOverride")
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfigOverride indicates an expected call of ConfigOverride.
func (mr *MockMongoSnapServiceMockRecorder) ConfigOverride() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfigOverride", reflect.TypeOf((*MockMongoSnapService)(nil).ConfigOverride))
}

// Exists mocks base method.
func (m *MockMongoSnapService) Exists() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockMongoSnapServiceMockRecorder) Exists() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockMongoSnapService)(nil).Exists))
}

// Install mocks base method.
func (m *MockMongoSnapService) Install() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Install")
	ret0, _ := ret[0].(error)
	return ret0
}

// Install indicates an expected call of Install.
func (mr *MockMongoSnapServiceMockRecorder) Install() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Install", reflect.TypeOf((*MockMongoSnapService)(nil).Install))
}

// Installed mocks base method.
func (m *MockMongoSnapService) Installed() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Installed")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Installed indicates an expected call of Installed.
func (mr *MockMongoSnapServiceMockRecorder) Installed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Installed", reflect.TypeOf((*MockMongoSnapService)(nil).Installed))
}

// Name mocks base method.
func (m *MockMongoSnapService) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockMongoSnapServiceMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockMongoSnapService)(nil).Name))
}

// Remove mocks base method.
func (m *MockMongoSnapService) Remove() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove")
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockMongoSnapServiceMockRecorder) Remove() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockMongoSnapService)(nil).Remove))
}

// Running mocks base method.
func (m *MockMongoSnapService) Running() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Running")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Running indicates an expected call of Running.
func (mr *MockMongoSnapServiceMockRecorder) Running() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Running", reflect.TypeOf((*MockMongoSnapService)(nil).Running))
}

// Start mocks base method.
func (m *MockMongoSnapService) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockMongoSnapServiceMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockMongoSnapService)(nil).Start))
}

// Stop mocks base method.
func (m *MockMongoSnapService) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockMongoSnapServiceMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockMongoSnapService)(nil).Stop))
}
