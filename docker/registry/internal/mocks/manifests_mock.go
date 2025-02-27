// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/docker/registry/internal (interfaces: ArchitectureGetter)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	internal "github.com/DavinZhang/juju/docker/registry/internal"
)

// MockArchitectureGetter is a mock of ArchitectureGetter interface.
type MockArchitectureGetter struct {
	ctrl     *gomock.Controller
	recorder *MockArchitectureGetterMockRecorder
}

// MockArchitectureGetterMockRecorder is the mock recorder for MockArchitectureGetter.
type MockArchitectureGetterMockRecorder struct {
	mock *MockArchitectureGetter
}

// NewMockArchitectureGetter creates a new mock instance.
func NewMockArchitectureGetter(ctrl *gomock.Controller) *MockArchitectureGetter {
	mock := &MockArchitectureGetter{ctrl: ctrl}
	mock.recorder = &MockArchitectureGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArchitectureGetter) EXPECT() *MockArchitectureGetterMockRecorder {
	return m.recorder
}

// GetBlobs mocks base method.
func (m *MockArchitectureGetter) GetBlobs(arg0, arg1 string) (*internal.BlobsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlobs", arg0, arg1)
	ret0, _ := ret[0].(*internal.BlobsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlobs indicates an expected call of GetBlobs.
func (mr *MockArchitectureGetterMockRecorder) GetBlobs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlobs", reflect.TypeOf((*MockArchitectureGetter)(nil).GetBlobs), arg0, arg1)
}

// GetManifests mocks base method.
func (m *MockArchitectureGetter) GetManifests(arg0, arg1 string) (*internal.ManifestsResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManifests", arg0, arg1)
	ret0, _ := ret[0].(*internal.ManifestsResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManifests indicates an expected call of GetManifests.
func (mr *MockArchitectureGetterMockRecorder) GetManifests(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManifests", reflect.TypeOf((*MockArchitectureGetter)(nil).GetManifests), arg0, arg1)
}
