// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/cmd/containeragent/utils (interfaces: FileReaderWriter)

// Package mocks is a generated GoMock package.
package mocks

import (
	io "io"
	os "os"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFileReaderWriter is a mock of FileReaderWriter interface
type MockFileReaderWriter struct {
	ctrl     *gomock.Controller
	recorder *MockFileReaderWriterMockRecorder
}

// MockFileReaderWriterMockRecorder is the mock recorder for MockFileReaderWriter
type MockFileReaderWriterMockRecorder struct {
	mock *MockFileReaderWriter
}

// NewMockFileReaderWriter creates a new mock instance
func NewMockFileReaderWriter(ctrl *gomock.Controller) *MockFileReaderWriter {
	mock := &MockFileReaderWriter{ctrl: ctrl}
	mock.recorder = &MockFileReaderWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileReaderWriter) EXPECT() *MockFileReaderWriterMockRecorder {
	return m.recorder
}

// MkdirAll mocks base method
func (m *MockFileReaderWriter) MkdirAll(arg0 string, arg1 os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MkdirAll", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// MkdirAll indicates an expected call of MkdirAll
func (mr *MockFileReaderWriterMockRecorder) MkdirAll(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MkdirAll", reflect.TypeOf((*MockFileReaderWriter)(nil).MkdirAll), arg0, arg1)
}

// ReadFile mocks base method
func (m *MockFileReaderWriter) ReadFile(arg0 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFile", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadFile indicates an expected call of ReadFile
func (mr *MockFileReaderWriterMockRecorder) ReadFile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFile", reflect.TypeOf((*MockFileReaderWriter)(nil).ReadFile), arg0)
}

// Reader mocks base method
func (m *MockFileReaderWriter) Reader(arg0 string) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reader", arg0)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Reader indicates an expected call of Reader
func (mr *MockFileReaderWriterMockRecorder) Reader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reader", reflect.TypeOf((*MockFileReaderWriter)(nil).Reader), arg0)
}

// RemoveAll mocks base method
func (m *MockFileReaderWriter) RemoveAll(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAll", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAll indicates an expected call of RemoveAll
func (mr *MockFileReaderWriterMockRecorder) RemoveAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAll", reflect.TypeOf((*MockFileReaderWriter)(nil).RemoveAll), arg0)
}

// Symlink mocks base method
func (m *MockFileReaderWriter) Symlink(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Symlink", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Symlink indicates an expected call of Symlink
func (mr *MockFileReaderWriterMockRecorder) Symlink(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Symlink", reflect.TypeOf((*MockFileReaderWriter)(nil).Symlink), arg0, arg1)
}

// WriteFile mocks base method
func (m *MockFileReaderWriter) WriteFile(arg0 string, arg1 []byte, arg2 os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteFile indicates an expected call of WriteFile
func (mr *MockFileReaderWriterMockRecorder) WriteFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFile", reflect.TypeOf((*MockFileReaderWriter)(nil).WriteFile), arg0, arg1, arg2)
}

// Writer mocks base method
func (m *MockFileReaderWriter) Writer(arg0 string, arg1 os.FileMode) (io.WriteCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Writer", arg0, arg1)
	ret0, _ := ret[0].(io.WriteCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Writer indicates an expected call of Writer
func (mr *MockFileReaderWriterMockRecorder) Writer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Writer", reflect.TypeOf((*MockFileReaderWriter)(nil).Writer), arg0, arg1)
}
