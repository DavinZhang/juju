// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/state (interfaces: Resources)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	resource "github.com/juju/charm/v9/resource"
	resource0 "github.com/DavinZhang/juju/resource"
	state "github.com/DavinZhang/juju/state"
	txn "github.com/juju/mgo/v2/txn"
	io "io"
	reflect "reflect"
	time "time"
)

// MockResources is a mock of Resources interface
type MockResources struct {
	ctrl     *gomock.Controller
	recorder *MockResourcesMockRecorder
}

// MockResourcesMockRecorder is the mock recorder for MockResources
type MockResourcesMockRecorder struct {
	mock *MockResources
}

// NewMockResources creates a new mock instance
func NewMockResources(ctrl *gomock.Controller) *MockResources {
	mock := &MockResources{ctrl: ctrl}
	mock.recorder = &MockResourcesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockResources) EXPECT() *MockResourcesMockRecorder {
	return m.recorder
}

// AddPendingResource mocks base method
func (m *MockResources) AddPendingResource(arg0, arg1 string, arg2 resource.Resource) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPendingResource", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddPendingResource indicates an expected call of AddPendingResource
func (mr *MockResourcesMockRecorder) AddPendingResource(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPendingResource", reflect.TypeOf((*MockResources)(nil).AddPendingResource), arg0, arg1, arg2)
}

// GetPendingResource mocks base method
func (m *MockResources) GetPendingResource(arg0, arg1, arg2 string) (resource0.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingResource", arg0, arg1, arg2)
	ret0, _ := ret[0].(resource0.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPendingResource indicates an expected call of GetPendingResource
func (mr *MockResourcesMockRecorder) GetPendingResource(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingResource", reflect.TypeOf((*MockResources)(nil).GetPendingResource), arg0, arg1, arg2)
}

// GetResource mocks base method
func (m *MockResources) GetResource(arg0, arg1 string) (resource0.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResource", arg0, arg1)
	ret0, _ := ret[0].(resource0.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResource indicates an expected call of GetResource
func (mr *MockResourcesMockRecorder) GetResource(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResource", reflect.TypeOf((*MockResources)(nil).GetResource), arg0, arg1)
}

// ListPendingResources mocks base method
func (m *MockResources) ListPendingResources(arg0 string) ([]resource0.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPendingResources", arg0)
	ret0, _ := ret[0].([]resource0.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPendingResources indicates an expected call of ListPendingResources
func (mr *MockResourcesMockRecorder) ListPendingResources(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPendingResources", reflect.TypeOf((*MockResources)(nil).ListPendingResources), arg0)
}

// ListResources mocks base method
func (m *MockResources) ListResources(arg0 string) (resource0.ApplicationResources, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListResources", arg0)
	ret0, _ := ret[0].(resource0.ApplicationResources)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListResources indicates an expected call of ListResources
func (mr *MockResourcesMockRecorder) ListResources(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListResources", reflect.TypeOf((*MockResources)(nil).ListResources), arg0)
}

// NewResolvePendingResourcesOps mocks base method
func (m *MockResources) NewResolvePendingResourcesOps(arg0 string, arg1 map[string]string) ([]txn.Op, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewResolvePendingResourcesOps", arg0, arg1)
	ret0, _ := ret[0].([]txn.Op)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewResolvePendingResourcesOps indicates an expected call of NewResolvePendingResourcesOps
func (mr *MockResourcesMockRecorder) NewResolvePendingResourcesOps(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewResolvePendingResourcesOps", reflect.TypeOf((*MockResources)(nil).NewResolvePendingResourcesOps), arg0, arg1)
}

// OpenResource mocks base method
func (m *MockResources) OpenResource(arg0, arg1 string) (resource0.Resource, io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenResource", arg0, arg1)
	ret0, _ := ret[0].(resource0.Resource)
	ret1, _ := ret[1].(io.ReadCloser)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// OpenResource indicates an expected call of OpenResource
func (mr *MockResourcesMockRecorder) OpenResource(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenResource", reflect.TypeOf((*MockResources)(nil).OpenResource), arg0, arg1)
}

// OpenResourceForUniter mocks base method
func (m *MockResources) OpenResourceForUniter(arg0 resource0.Unit, arg1 string) (resource0.Resource, io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenResourceForUniter", arg0, arg1)
	ret0, _ := ret[0].(resource0.Resource)
	ret1, _ := ret[1].(io.ReadCloser)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// OpenResourceForUniter indicates an expected call of OpenResourceForUniter
func (mr *MockResourcesMockRecorder) OpenResourceForUniter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenResourceForUniter", reflect.TypeOf((*MockResources)(nil).OpenResourceForUniter), arg0, arg1)
}

// RemovePendingAppResources mocks base method
func (m *MockResources) RemovePendingAppResources(arg0 string, arg1 map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemovePendingAppResources", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePendingAppResources indicates an expected call of RemovePendingAppResources
func (mr *MockResourcesMockRecorder) RemovePendingAppResources(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePendingAppResources", reflect.TypeOf((*MockResources)(nil).RemovePendingAppResources), arg0, arg1)
}

// SetCharmStoreResources mocks base method
func (m *MockResources) SetCharmStoreResources(arg0 string, arg1 []resource.Resource, arg2 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCharmStoreResources", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCharmStoreResources indicates an expected call of SetCharmStoreResources
func (mr *MockResourcesMockRecorder) SetCharmStoreResources(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCharmStoreResources", reflect.TypeOf((*MockResources)(nil).SetCharmStoreResources), arg0, arg1, arg2)
}

// SetResource mocks base method
func (m *MockResources) SetResource(arg0, arg1 string, arg2 resource.Resource, arg3 io.Reader, arg4 state.IncrementCharmModifiedVersionType) (resource0.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetResource", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(resource0.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetResource indicates an expected call of SetResource
func (mr *MockResourcesMockRecorder) SetResource(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetResource", reflect.TypeOf((*MockResources)(nil).SetResource), arg0, arg1, arg2, arg3, arg4)
}

// SetUnitResource mocks base method
func (m *MockResources) SetUnitResource(arg0, arg1 string, arg2 resource.Resource) (resource0.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUnitResource", arg0, arg1, arg2)
	ret0, _ := ret[0].(resource0.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetUnitResource indicates an expected call of SetUnitResource
func (mr *MockResourcesMockRecorder) SetUnitResource(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUnitResource", reflect.TypeOf((*MockResources)(nil).SetUnitResource), arg0, arg1, arg2)
}

// UpdatePendingResource mocks base method
func (m *MockResources) UpdatePendingResource(arg0, arg1, arg2 string, arg3 resource.Resource, arg4 io.Reader) (resource0.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePendingResource", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(resource0.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePendingResource indicates an expected call of UpdatePendingResource
func (mr *MockResourcesMockRecorder) UpdatePendingResource(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePendingResource", reflect.TypeOf((*MockResources)(nil).UpdatePendingResource), arg0, arg1, arg2, arg3, arg4)
}
