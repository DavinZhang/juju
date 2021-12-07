// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/worker/uniter/runner/context (interfaces: HookUnit)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	charm "github.com/juju/charm/v9"
	uniter "github.com/DavinZhang/juju/api/uniter"
	params "github.com/DavinZhang/juju/apiserver/params"
	status "github.com/DavinZhang/juju/core/status"
	names "github.com/juju/names/v4"
)

// MockHookUnit is a mock of HookUnit interface.
type MockHookUnit struct {
	ctrl     *gomock.Controller
	recorder *MockHookUnitMockRecorder
}

// MockHookUnitMockRecorder is the mock recorder for MockHookUnit.
type MockHookUnitMockRecorder struct {
	mock *MockHookUnit
}

// NewMockHookUnit creates a new mock instance.
func NewMockHookUnit(ctrl *gomock.Controller) *MockHookUnit {
	mock := &MockHookUnit{ctrl: ctrl}
	mock.recorder = &MockHookUnitMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHookUnit) EXPECT() *MockHookUnitMockRecorder {
	return m.recorder
}

// Application mocks base method.
func (m *MockHookUnit) Application() (*uniter.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Application")
	ret0, _ := ret[0].(*uniter.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Application indicates an expected call of Application.
func (mr *MockHookUnitMockRecorder) Application() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Application", reflect.TypeOf((*MockHookUnit)(nil).Application))
}

// ApplicationName mocks base method.
func (m *MockHookUnit) ApplicationName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplicationName")
	ret0, _ := ret[0].(string)
	return ret0
}

// ApplicationName indicates an expected call of ApplicationName.
func (mr *MockHookUnitMockRecorder) ApplicationName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplicationName", reflect.TypeOf((*MockHookUnit)(nil).ApplicationName))
}

// CommitHookChanges mocks base method.
func (m *MockHookUnit) CommitHookChanges(arg0 params.CommitHookChangesArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitHookChanges", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitHookChanges indicates an expected call of CommitHookChanges.
func (mr *MockHookUnitMockRecorder) CommitHookChanges(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitHookChanges", reflect.TypeOf((*MockHookUnit)(nil).CommitHookChanges), arg0)
}

// ConfigSettings mocks base method.
func (m *MockHookUnit) ConfigSettings() (charm.Settings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfigSettings")
	ret0, _ := ret[0].(charm.Settings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConfigSettings indicates an expected call of ConfigSettings.
func (mr *MockHookUnitMockRecorder) ConfigSettings() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfigSettings", reflect.TypeOf((*MockHookUnit)(nil).ConfigSettings))
}

// LogActionMessage mocks base method.
func (m *MockHookUnit) LogActionMessage(arg0 names.ActionTag, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogActionMessage", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogActionMessage indicates an expected call of LogActionMessage.
func (mr *MockHookUnitMockRecorder) LogActionMessage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogActionMessage", reflect.TypeOf((*MockHookUnit)(nil).LogActionMessage), arg0, arg1)
}

// Name mocks base method.
func (m *MockHookUnit) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockHookUnitMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockHookUnit)(nil).Name))
}

// NetworkInfo mocks base method.
func (m *MockHookUnit) NetworkInfo(arg0 []string, arg1 *int) (map[string]params.NetworkInfoResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NetworkInfo", arg0, arg1)
	ret0, _ := ret[0].(map[string]params.NetworkInfoResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NetworkInfo indicates an expected call of NetworkInfo.
func (mr *MockHookUnitMockRecorder) NetworkInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NetworkInfo", reflect.TypeOf((*MockHookUnit)(nil).NetworkInfo), arg0, arg1)
}

// RequestReboot mocks base method.
func (m *MockHookUnit) RequestReboot() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestReboot")
	ret0, _ := ret[0].(error)
	return ret0
}

// RequestReboot indicates an expected call of RequestReboot.
func (mr *MockHookUnitMockRecorder) RequestReboot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestReboot", reflect.TypeOf((*MockHookUnit)(nil).RequestReboot))
}

// SetAgentStatus mocks base method.
func (m *MockHookUnit) SetAgentStatus(arg0 status.Status, arg1 string, arg2 map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAgentStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAgentStatus indicates an expected call of SetAgentStatus.
func (mr *MockHookUnitMockRecorder) SetAgentStatus(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAgentStatus", reflect.TypeOf((*MockHookUnit)(nil).SetAgentStatus), arg0, arg1, arg2)
}

// SetUnitStatus mocks base method.
func (m *MockHookUnit) SetUnitStatus(arg0 status.Status, arg1 string, arg2 map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUnitStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUnitStatus indicates an expected call of SetUnitStatus.
func (mr *MockHookUnitMockRecorder) SetUnitStatus(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUnitStatus", reflect.TypeOf((*MockHookUnit)(nil).SetUnitStatus), arg0, arg1, arg2)
}

// State mocks base method.
func (m *MockHookUnit) State() (params.UnitStateResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(params.UnitStateResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// State indicates an expected call of State.
func (mr *MockHookUnitMockRecorder) State() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockHookUnit)(nil).State))
}

// Tag mocks base method.
func (m *MockHookUnit) Tag() names.UnitTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tag")
	ret0, _ := ret[0].(names.UnitTag)
	return ret0
}

// Tag indicates an expected call of Tag.
func (mr *MockHookUnitMockRecorder) Tag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tag", reflect.TypeOf((*MockHookUnit)(nil).Tag))
}

// UnitStatus mocks base method.
func (m *MockHookUnit) UnitStatus() (params.StatusResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnitStatus")
	ret0, _ := ret[0].(params.StatusResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnitStatus indicates an expected call of UnitStatus.
func (mr *MockHookUnitMockRecorder) UnitStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnitStatus", reflect.TypeOf((*MockHookUnit)(nil).UnitStatus))
}

// UpdateNetworkInfo mocks base method.
func (m *MockHookUnit) UpdateNetworkInfo() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNetworkInfo")
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNetworkInfo indicates an expected call of UpdateNetworkInfo.
func (mr *MockHookUnitMockRecorder) UpdateNetworkInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNetworkInfo", reflect.TypeOf((*MockHookUnit)(nil).UpdateNetworkInfo))
}
