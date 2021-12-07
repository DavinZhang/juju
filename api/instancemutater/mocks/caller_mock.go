// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/api/base (interfaces: APICaller,FacadeCaller)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	http "net/http"
	url "net/url"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	base "github.com/DavinZhang/juju/api/base"
	names_v3 "github.com/juju/names/v4"
	httprequest_v1 "gopkg.in/httprequest.v1"
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
	ret := m.ctrl.Call(m, "APICall", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// APICall indicates an expected call of APICall
func (mr *MockAPICallerMockRecorder) APICall(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "APICall", reflect.TypeOf((*MockAPICaller)(nil).APICall), arg0, arg1, arg2, arg3, arg4, arg5)
}

// BakeryClient mocks base method
func (m *MockAPICaller) BakeryClient() base.MacaroonDischarger {
	ret := m.ctrl.Call(m, "BakeryClient")
	ret0, _ := ret[0].(base.MacaroonDischarger)
	return ret0
}

// BakeryClient indicates an expected call of BakeryClient
func (mr *MockAPICallerMockRecorder) BakeryClient() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BakeryClient", reflect.TypeOf((*MockAPICaller)(nil).BakeryClient))
}

// BestFacadeVersion mocks base method
func (m *MockAPICaller) BestFacadeVersion(arg0 string) int {
	ret := m.ctrl.Call(m, "BestFacadeVersion", arg0)
	ret0, _ := ret[0].(int)
	return ret0
}

// BestFacadeVersion indicates an expected call of BestFacadeVersion
func (mr *MockAPICallerMockRecorder) BestFacadeVersion(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BestFacadeVersion", reflect.TypeOf((*MockAPICaller)(nil).BestFacadeVersion), arg0)
}

// ConnectControllerStream mocks base method
func (m *MockAPICaller) ConnectControllerStream(arg0 string, arg1 url.Values, arg2 http.Header) (base.Stream, error) {
	ret := m.ctrl.Call(m, "ConnectControllerStream", arg0, arg1, arg2)
	ret0, _ := ret[0].(base.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConnectControllerStream indicates an expected call of ConnectControllerStream
func (mr *MockAPICallerMockRecorder) ConnectControllerStream(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectControllerStream", reflect.TypeOf((*MockAPICaller)(nil).ConnectControllerStream), arg0, arg1, arg2)
}

// ConnectStream mocks base method
func (m *MockAPICaller) ConnectStream(arg0 string, arg1 url.Values) (base.Stream, error) {
	ret := m.ctrl.Call(m, "ConnectStream", arg0, arg1)
	ret0, _ := ret[0].(base.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConnectStream indicates an expected call of ConnectStream
func (mr *MockAPICallerMockRecorder) ConnectStream(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectStream", reflect.TypeOf((*MockAPICaller)(nil).ConnectStream), arg0, arg1)
}

// Context mocks base method
func (m *MockAPICaller) Context() context.Context {
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockAPICallerMockRecorder) Context() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockAPICaller)(nil).Context))
}

// HTTPClient mocks base method
func (m *MockAPICaller) HTTPClient() (*httprequest_v1.Client, error) {
	ret := m.ctrl.Call(m, "HTTPClient")
	ret0, _ := ret[0].(*httprequest_v1.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HTTPClient indicates an expected call of HTTPClient
func (mr *MockAPICallerMockRecorder) HTTPClient() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HTTPClient", reflect.TypeOf((*MockAPICaller)(nil).HTTPClient))
}

// ModelTag mocks base method
func (m *MockAPICaller) ModelTag() (names_v3.ModelTag, bool) {
	ret := m.ctrl.Call(m, "ModelTag")
	ret0, _ := ret[0].(names_v3.ModelTag)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// ModelTag indicates an expected call of ModelTag
func (mr *MockAPICallerMockRecorder) ModelTag() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelTag", reflect.TypeOf((*MockAPICaller)(nil).ModelTag))
}

// MockFacadeCaller is a mock of FacadeCaller interface
type MockFacadeCaller struct {
	ctrl     *gomock.Controller
	recorder *MockFacadeCallerMockRecorder
}

// MockFacadeCallerMockRecorder is the mock recorder for MockFacadeCaller
type MockFacadeCallerMockRecorder struct {
	mock *MockFacadeCaller
}

// NewMockFacadeCaller creates a new mock instance
func NewMockFacadeCaller(ctrl *gomock.Controller) *MockFacadeCaller {
	mock := &MockFacadeCaller{ctrl: ctrl}
	mock.recorder = &MockFacadeCallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFacadeCaller) EXPECT() *MockFacadeCallerMockRecorder {
	return m.recorder
}

// BestAPIVersion mocks base method
func (m *MockFacadeCaller) BestAPIVersion() int {
	ret := m.ctrl.Call(m, "BestAPIVersion")
	ret0, _ := ret[0].(int)
	return ret0
}

// BestAPIVersion indicates an expected call of BestAPIVersion
func (mr *MockFacadeCallerMockRecorder) BestAPIVersion() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BestAPIVersion", reflect.TypeOf((*MockFacadeCaller)(nil).BestAPIVersion))
}

// FacadeCall mocks base method
func (m *MockFacadeCaller) FacadeCall(arg0 string, arg1, arg2 interface{}) error {
	ret := m.ctrl.Call(m, "FacadeCall", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// FacadeCall indicates an expected call of FacadeCall
func (mr *MockFacadeCallerMockRecorder) FacadeCall(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FacadeCall", reflect.TypeOf((*MockFacadeCaller)(nil).FacadeCall), arg0, arg1, arg2)
}

// Name mocks base method
func (m *MockFacadeCaller) Name() string {
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockFacadeCallerMockRecorder) Name() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockFacadeCaller)(nil).Name))
}

// RawAPICaller mocks base method
func (m *MockFacadeCaller) RawAPICaller() base.APICaller {
	ret := m.ctrl.Call(m, "RawAPICaller")
	ret0, _ := ret[0].(base.APICaller)
	return ret0
}

// RawAPICaller indicates an expected call of RawAPICaller
func (mr *MockFacadeCallerMockRecorder) RawAPICaller() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RawAPICaller", reflect.TypeOf((*MockFacadeCaller)(nil).RawAPICaller))
}
