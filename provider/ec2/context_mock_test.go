// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/environs/context (interfaces: ProviderCallContext)

// Package ec2 is a generated GoMock package.
package ec2

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockProviderCallContext is a mock of ProviderCallContext interface.
type MockProviderCallContext struct {
	ctrl     *gomock.Controller
	recorder *MockProviderCallContextMockRecorder
}

// MockProviderCallContextMockRecorder is the mock recorder for MockProviderCallContext.
type MockProviderCallContextMockRecorder struct {
	mock *MockProviderCallContext
}

// NewMockProviderCallContext creates a new mock instance.
func NewMockProviderCallContext(ctrl *gomock.Controller) *MockProviderCallContext {
	mock := &MockProviderCallContext{ctrl: ctrl}
	mock.recorder = &MockProviderCallContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProviderCallContext) EXPECT() *MockProviderCallContextMockRecorder {
	return m.recorder
}

// Deadline mocks base method.
func (m *MockProviderCallContext) Deadline() (time.Time, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deadline")
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Deadline indicates an expected call of Deadline.
func (mr *MockProviderCallContextMockRecorder) Deadline() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deadline", reflect.TypeOf((*MockProviderCallContext)(nil).Deadline))
}

// Done mocks base method.
func (m *MockProviderCallContext) Done() <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Done")
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Done indicates an expected call of Done.
func (mr *MockProviderCallContextMockRecorder) Done() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Done", reflect.TypeOf((*MockProviderCallContext)(nil).Done))
}

// Err mocks base method.
func (m *MockProviderCallContext) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err.
func (mr *MockProviderCallContextMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockProviderCallContext)(nil).Err))
}

// InvalidateCredential mocks base method.
func (m *MockProviderCallContext) InvalidateCredential(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvalidateCredential", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InvalidateCredential indicates an expected call of InvalidateCredential.
func (mr *MockProviderCallContextMockRecorder) InvalidateCredential(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvalidateCredential", reflect.TypeOf((*MockProviderCallContext)(nil).InvalidateCredential), arg0)
}

// Value mocks base method.
func (m *MockProviderCallContext) Value(arg0 interface{}) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Value", arg0)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Value indicates an expected call of Value.
func (mr *MockProviderCallContextMockRecorder) Value(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Value", reflect.TypeOf((*MockProviderCallContext)(nil).Value), arg0)
}
