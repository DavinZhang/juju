// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/resource/repositories (interfaces: EntityRepository,ResourceGetter)

// Package mocks is a generated GoMock package.
package mocks

import (
	io "io"
	reflect "reflect"
	sync "sync"

	gomock "github.com/golang/mock/gomock"
	resource "github.com/juju/charm/v9/resource"
	charmstore "github.com/DavinZhang/juju/charmstore"
	resource0 "github.com/DavinZhang/juju/resource"
	repositories "github.com/DavinZhang/juju/resource/repositories"
	state "github.com/DavinZhang/juju/state"
)

// MockEntityRepository is a mock of EntityRepository interface.
type MockEntityRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEntityRepositoryMockRecorder
}

// MockEntityRepositoryMockRecorder is the mock recorder for MockEntityRepository.
type MockEntityRepositoryMockRecorder struct {
	mock *MockEntityRepository
}

// NewMockEntityRepository creates a new mock instance.
func NewMockEntityRepository(ctrl *gomock.Controller) *MockEntityRepository {
	mock := &MockEntityRepository{ctrl: ctrl}
	mock.recorder = &MockEntityRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEntityRepository) EXPECT() *MockEntityRepositoryMockRecorder {
	return m.recorder
}

// FetchLock mocks base method.
func (m *MockEntityRepository) FetchLock(arg0 string) sync.Locker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchLock", arg0)
	ret0, _ := ret[0].(sync.Locker)
	return ret0
}

// FetchLock indicates an expected call of FetchLock.
func (mr *MockEntityRepositoryMockRecorder) FetchLock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchLock", reflect.TypeOf((*MockEntityRepository)(nil).FetchLock), arg0)
}

// GetResource mocks base method.
func (m *MockEntityRepository) GetResource(arg0 string) (resource0.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResource", arg0)
	ret0, _ := ret[0].(resource0.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResource indicates an expected call of GetResource.
func (mr *MockEntityRepositoryMockRecorder) GetResource(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResource", reflect.TypeOf((*MockEntityRepository)(nil).GetResource), arg0)
}

// OpenResource mocks base method.
func (m *MockEntityRepository) OpenResource(arg0 string) (resource0.Resource, io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenResource", arg0)
	ret0, _ := ret[0].(resource0.Resource)
	ret1, _ := ret[1].(io.ReadCloser)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// OpenResource indicates an expected call of OpenResource.
func (mr *MockEntityRepositoryMockRecorder) OpenResource(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenResource", reflect.TypeOf((*MockEntityRepository)(nil).OpenResource), arg0)
}

// SetResource mocks base method.
func (m *MockEntityRepository) SetResource(arg0 resource.Resource, arg1 io.Reader, arg2 state.IncrementCharmModifiedVersionType) (resource0.Resource, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetResource", arg0, arg1, arg2)
	ret0, _ := ret[0].(resource0.Resource)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetResource indicates an expected call of SetResource.
func (mr *MockEntityRepositoryMockRecorder) SetResource(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetResource", reflect.TypeOf((*MockEntityRepository)(nil).SetResource), arg0, arg1, arg2)
}

// MockResourceGetter is a mock of ResourceGetter interface.
type MockResourceGetter struct {
	ctrl     *gomock.Controller
	recorder *MockResourceGetterMockRecorder
}

// MockResourceGetterMockRecorder is the mock recorder for MockResourceGetter.
type MockResourceGetterMockRecorder struct {
	mock *MockResourceGetter
}

// NewMockResourceGetter creates a new mock instance.
func NewMockResourceGetter(ctrl *gomock.Controller) *MockResourceGetter {
	mock := &MockResourceGetter{ctrl: ctrl}
	mock.recorder = &MockResourceGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResourceGetter) EXPECT() *MockResourceGetterMockRecorder {
	return m.recorder
}

// GetResource mocks base method.
func (m *MockResourceGetter) GetResource(arg0 repositories.ResourceRequest) (charmstore.ResourceData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResource", arg0)
	ret0, _ := ret[0].(charmstore.ResourceData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResource indicates an expected call of GetResource.
func (mr *MockResourceGetterMockRecorder) GetResource(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResource", reflect.TypeOf((*MockResourceGetter)(nil).GetResource), arg0)
}
