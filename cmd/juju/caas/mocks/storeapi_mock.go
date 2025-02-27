// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/cmd/juju/caas (interfaces: CredentialStoreAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	cloud "github.com/DavinZhang/juju/cloud"
)

// MockCredentialStoreAPI is a mock of CredentialStoreAPI interface
type MockCredentialStoreAPI struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialStoreAPIMockRecorder
}

// MockCredentialStoreAPIMockRecorder is the mock recorder for MockCredentialStoreAPI
type MockCredentialStoreAPIMockRecorder struct {
	mock *MockCredentialStoreAPI
}

// NewMockCredentialStoreAPI creates a new mock instance
func NewMockCredentialStoreAPI(ctrl *gomock.Controller) *MockCredentialStoreAPI {
	mock := &MockCredentialStoreAPI{ctrl: ctrl}
	mock.recorder = &MockCredentialStoreAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCredentialStoreAPI) EXPECT() *MockCredentialStoreAPIMockRecorder {
	return m.recorder
}

// UpdateCredential mocks base method
func (m *MockCredentialStoreAPI) UpdateCredential(arg0 string, arg1 cloud.CloudCredential) error {
	ret := m.ctrl.Call(m, "UpdateCredential", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCredential indicates an expected call of UpdateCredential
func (mr *MockCredentialStoreAPIMockRecorder) UpdateCredential(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCredential", reflect.TypeOf((*MockCredentialStoreAPI)(nil).UpdateCredential), arg0, arg1)
}
