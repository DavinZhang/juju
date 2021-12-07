// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/api/base (interfaces: APICaller)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	base "github.com/DavinZhang/juju/api/base"
	names "github.com/juju/names/v4"
	httprequest "gopkg.in/httprequest.v1"
	http "net/http"
	url "net/url"
	reflect "reflect"
)

// MockAPICaller is a mock of APICaller interface
type MockAPICaller struct {
	ctrl     *gomock.Controller
	recorder *MockAPICallerMockRecorder
}

// MockAPICallerMockRecorder is the mock recorder for MockAPICaller
type MockAPICallerMockRecorder struct {
	mock *MockAPICaller
}

// NewMockAPICaller creates a new mock instance
func NewMockAPICaller(ctrl *gomock.Controller) *MockAPICaller {
	mock := &MockAPICaller{ctrl: ctrl}
	mock.recorder = &MockAPICallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAPICaller) EXPECT() *MockAPICallerMockRecorder {
	return m.recorder
}

// APICall mocks base method
func (m *MockAPICaller) APICall(arg0 string, arg1 int, arg2, arg3 string, arg4, arg5 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "APICall", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// APICall indicates an expected call of APICall
func (mr *MockAPICallerMockRecorder) APICall(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APICall", reflect.TypeOf((*MockAPICaller)(nil).APICall), arg0, arg1, arg2, arg3, arg4, arg5)
}

// BakeryClient mocks base method
func (m *MockAPICaller) BakeryClient() base.MacaroonDischarger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BakeryClient")
	ret0, _ := ret[0].(base.MacaroonDischarger)
	return ret0
}

// BakeryClient indicates an expected call of BakeryClient
func (mr *MockAPICallerMockRecorder) BakeryClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BakeryClient", reflect.TypeOf((*MockAPICaller)(nil).BakeryClient))
}

// BestFacadeVersion mocks base method
func (m *MockAPICaller) BestFacadeVersion(arg0 string) int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BestFacadeVersion", arg0)
	ret0, _ := ret[0].(int)
	return ret0
}

// BestFacadeVersion indicates an expected call of BestFacadeVersion
func (mr *MockAPICallerMockRecorder) BestFacadeVersion(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BestFacadeVersion", reflect.TypeOf((*MockAPICaller)(nil).BestFacadeVersion), arg0)
}

// ConnectControllerStream mocks base method
func (m *MockAPICaller) ConnectControllerStream(arg0 string, arg1 url.Values, arg2 http.Header) (base.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectControllerStream", arg0, arg1, arg2)
	ret0, _ := ret[0].(base.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConnectControllerStream indicates an expected call of ConnectControllerStream
func (mr *MockAPICallerMockRecorder) ConnectControllerStream(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectControllerStream", reflect.TypeOf((*MockAPICaller)(nil).ConnectControllerStream), arg0, arg1, arg2)
}

// ConnectStream mocks base method
func (m *MockAPICaller) ConnectStream(arg0 string, arg1 url.Values) (base.Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectStream", arg0, arg1)
	ret0, _ := ret[0].(base.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConnectStream indicates an expected call of ConnectStream
func (mr *MockAPICallerMockRecorder) ConnectStream(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectStream", reflect.TypeOf((*MockAPICaller)(nil).ConnectStream), arg0, arg1)
}

// Context mocks base method
func (m *MockAPICaller) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockAPICallerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockAPICaller)(nil).Context))
}

// HTTPClient mocks base method
func (m *MockAPICaller) HTTPClient() (*httprequest.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HTTPClient")
	ret0, _ := ret[0].(*httprequest.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HTTPClient indicates an expected call of HTTPClient
func (mr *MockAPICallerMockRecorder) HTTPClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HTTPClient", reflect.TypeOf((*MockAPICaller)(nil).HTTPClient))
}

// ModelTag mocks base method
func (m *MockAPICaller) ModelTag() (names.ModelTag, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelTag")
	ret0, _ := ret[0].(names.ModelTag)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ModelTag indicates an expected call of ModelTag
func (mr *MockAPICallerMockRecorder) ModelTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelTag", reflect.TypeOf((*MockAPICaller)(nil).ModelTag))
}
