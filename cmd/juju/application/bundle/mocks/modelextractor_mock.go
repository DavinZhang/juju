// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/cmd/juju/application/bundle (interfaces: ModelExtractor)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	params "github.com/DavinZhang/juju/apiserver/params"
	constraints "github.com/DavinZhang/juju/core/constraints"
	reflect "reflect"
)

// MockModelExtractor is a mock of ModelExtractor interface
type MockModelExtractor struct {
	ctrl     *gomock.Controller
	recorder *MockModelExtractorMockRecorder
}

// MockModelExtractorMockRecorder is the mock recorder for MockModelExtractor
type MockModelExtractorMockRecorder struct {
	mock *MockModelExtractor
}

// NewMockModelExtractor creates a new mock instance
func NewMockModelExtractor(ctrl *gomock.Controller) *MockModelExtractor {
	mock := &MockModelExtractor{ctrl: ctrl}
	mock.recorder = &MockModelExtractorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockModelExtractor) EXPECT() *MockModelExtractorMockRecorder {
	return m.recorder
}

// GetAnnotations mocks base method
func (m *MockModelExtractor) GetAnnotations(arg0 []string) ([]params.AnnotationsGetResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnnotations", arg0)
	ret0, _ := ret[0].([]params.AnnotationsGetResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnnotations indicates an expected call of GetAnnotations
func (mr *MockModelExtractorMockRecorder) GetAnnotations(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnnotations", reflect.TypeOf((*MockModelExtractor)(nil).GetAnnotations), arg0)
}

// GetConfig mocks base method
func (m *MockModelExtractor) GetConfig(arg0 string, arg1 ...string) ([]map[string]interface{}, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetConfig", varargs...)
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfig indicates an expected call of GetConfig
func (mr *MockModelExtractorMockRecorder) GetConfig(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockModelExtractor)(nil).GetConfig), varargs...)
}

// GetConstraints mocks base method
func (m *MockModelExtractor) GetConstraints(arg0 ...string) ([]constraints.Value, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetConstraints", varargs...)
	ret0, _ := ret[0].([]constraints.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConstraints indicates an expected call of GetConstraints
func (mr *MockModelExtractorMockRecorder) GetConstraints(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConstraints", reflect.TypeOf((*MockModelExtractor)(nil).GetConstraints), arg0...)
}

// Sequences mocks base method
func (m *MockModelExtractor) Sequences() (map[string]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sequences")
	ret0, _ := ret[0].(map[string]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sequences indicates an expected call of Sequences
func (mr *MockModelExtractorMockRecorder) Sequences() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sequences", reflect.TypeOf((*MockModelExtractor)(nil).Sequences))
}
