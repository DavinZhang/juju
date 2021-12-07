// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/state/migrations (interfaces: MigrationExternalController,ExternalControllerSource,ExternalControllerModel)

// Package migrations is a generated GoMock package.
package migrations

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	description "github.com/juju/description/v3"
)

// MockMigrationExternalController is a mock of MigrationExternalController interface
type MockMigrationExternalController struct {
	ctrl     *gomock.Controller
	recorder *MockMigrationExternalControllerMockRecorder
}

// MockMigrationExternalControllerMockRecorder is the mock recorder for MockMigrationExternalController
type MockMigrationExternalControllerMockRecorder struct {
	mock *MockMigrationExternalController
}

// NewMockMigrationExternalController creates a new mock instance
func NewMockMigrationExternalController(ctrl *gomock.Controller) *MockMigrationExternalController {
	mock := &MockMigrationExternalController{ctrl: ctrl}
	mock.recorder = &MockMigrationExternalControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMigrationExternalController) EXPECT() *MockMigrationExternalControllerMockRecorder {
	return m.recorder
}

// Addrs mocks base method
func (m *MockMigrationExternalController) Addrs() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Addrs")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Addrs indicates an expected call of Addrs
func (mr *MockMigrationExternalControllerMockRecorder) Addrs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Addrs", reflect.TypeOf((*MockMigrationExternalController)(nil).Addrs))
}

// Alias mocks base method
func (m *MockMigrationExternalController) Alias() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Alias")
	ret0, _ := ret[0].(string)
	return ret0
}

// Alias indicates an expected call of Alias
func (mr *MockMigrationExternalControllerMockRecorder) Alias() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Alias", reflect.TypeOf((*MockMigrationExternalController)(nil).Alias))
}

// CACert mocks base method
func (m *MockMigrationExternalController) CACert() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CACert")
	ret0, _ := ret[0].(string)
	return ret0
}

// CACert indicates an expected call of CACert
func (mr *MockMigrationExternalControllerMockRecorder) CACert() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CACert", reflect.TypeOf((*MockMigrationExternalController)(nil).CACert))
}

// ID mocks base method
func (m *MockMigrationExternalController) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockMigrationExternalControllerMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockMigrationExternalController)(nil).ID))
}

// Models mocks base method
func (m *MockMigrationExternalController) Models() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Models")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Models indicates an expected call of Models
func (mr *MockMigrationExternalControllerMockRecorder) Models() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Models", reflect.TypeOf((*MockMigrationExternalController)(nil).Models))
}

// MockExternalControllerSource is a mock of ExternalControllerSource interface
type MockExternalControllerSource struct {
	ctrl     *gomock.Controller
	recorder *MockExternalControllerSourceMockRecorder
}

// MockExternalControllerSourceMockRecorder is the mock recorder for MockExternalControllerSource
type MockExternalControllerSourceMockRecorder struct {
	mock *MockExternalControllerSource
}

// NewMockExternalControllerSource creates a new mock instance
func NewMockExternalControllerSource(ctrl *gomock.Controller) *MockExternalControllerSource {
	mock := &MockExternalControllerSource{ctrl: ctrl}
	mock.recorder = &MockExternalControllerSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExternalControllerSource) EXPECT() *MockExternalControllerSourceMockRecorder {
	return m.recorder
}

// AllRemoteApplications mocks base method
func (m *MockExternalControllerSource) AllRemoteApplications() ([]MigrationRemoteApplication, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllRemoteApplications")
	ret0, _ := ret[0].([]MigrationRemoteApplication)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllRemoteApplications indicates an expected call of AllRemoteApplications
func (mr *MockExternalControllerSourceMockRecorder) AllRemoteApplications() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllRemoteApplications", reflect.TypeOf((*MockExternalControllerSource)(nil).AllRemoteApplications))
}

// ControllerForModel mocks base method
func (m *MockExternalControllerSource) ControllerForModel(arg0 string) (MigrationExternalController, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerForModel", arg0)
	ret0, _ := ret[0].(MigrationExternalController)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ControllerForModel indicates an expected call of ControllerForModel
func (mr *MockExternalControllerSourceMockRecorder) ControllerForModel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerForModel", reflect.TypeOf((*MockExternalControllerSource)(nil).ControllerForModel), arg0)
}

// MockExternalControllerModel is a mock of ExternalControllerModel interface
type MockExternalControllerModel struct {
	ctrl     *gomock.Controller
	recorder *MockExternalControllerModelMockRecorder
}

// MockExternalControllerModelMockRecorder is the mock recorder for MockExternalControllerModel
type MockExternalControllerModelMockRecorder struct {
	mock *MockExternalControllerModel
}

// NewMockExternalControllerModel creates a new mock instance
func NewMockExternalControllerModel(ctrl *gomock.Controller) *MockExternalControllerModel {
	mock := &MockExternalControllerModel{ctrl: ctrl}
	mock.recorder = &MockExternalControllerModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExternalControllerModel) EXPECT() *MockExternalControllerModelMockRecorder {
	return m.recorder
}

// AddExternalController mocks base method
func (m *MockExternalControllerModel) AddExternalController(arg0 description.ExternalControllerArgs) description.ExternalController {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddExternalController", arg0)
	ret0, _ := ret[0].(description.ExternalController)
	return ret0
}

// AddExternalController indicates an expected call of AddExternalController
func (mr *MockExternalControllerModelMockRecorder) AddExternalController(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddExternalController", reflect.TypeOf((*MockExternalControllerModel)(nil).AddExternalController), arg0)
}
