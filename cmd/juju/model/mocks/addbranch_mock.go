// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/cmd/juju/model (interfaces: AddBranchCommandAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAddBranchCommandAPI is a mock of AddBranchCommandAPI interface
type MockAddBranchCommandAPI struct {
	ctrl     *gomock.Controller
	recorder *MockAddBranchCommandAPIMockRecorder
}

// MockAddBranchCommandAPIMockRecorder is the mock recorder for MockAddBranchCommandAPI
type MockAddBranchCommandAPIMockRecorder struct {
	mock *MockAddBranchCommandAPI
}

// NewMockAddBranchCommandAPI creates a new mock instance
func NewMockAddBranchCommandAPI(ctrl *gomock.Controller) *MockAddBranchCommandAPI {
	mock := &MockAddBranchCommandAPI{ctrl: ctrl}
	mock.recorder = &MockAddBranchCommandAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAddBranchCommandAPI) EXPECT() *MockAddBranchCommandAPIMockRecorder {
	return m.recorder
}

// AddBranch mocks base method
func (m *MockAddBranchCommandAPI) AddBranch(arg0 string) error {
	ret := m.ctrl.Call(m, "AddBranch", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBranch indicates an expected call of AddBranch
func (mr *MockAddBranchCommandAPIMockRecorder) AddBranch(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBranch", reflect.TypeOf((*MockAddBranchCommandAPI)(nil).AddBranch), arg0)
}

// Close mocks base method
func (m *MockAddBranchCommandAPI) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockAddBranchCommandAPIMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAddBranchCommandAPI)(nil).Close))
}
