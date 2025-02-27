// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/core/watcher (interfaces: StringsWatcher)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	watcher "github.com/DavinZhang/juju/core/watcher"
)

// MockStringsWatcher is a mock of StringsWatcher interface
type MockStringsWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockStringsWatcherMockRecorder
}

// MockStringsWatcherMockRecorder is the mock recorder for MockStringsWatcher
type MockStringsWatcherMockRecorder struct {
	mock *MockStringsWatcher
}

// NewMockStringsWatcher creates a new mock instance
func NewMockStringsWatcher(ctrl *gomock.Controller) *MockStringsWatcher {
	mock := &MockStringsWatcher{ctrl: ctrl}
	mock.recorder = &MockStringsWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStringsWatcher) EXPECT() *MockStringsWatcherMockRecorder {
	return m.recorder
}

// Changes mocks base method
func (m *MockStringsWatcher) Changes() watcher.StringsChannel {
	ret := m.ctrl.Call(m, "Changes")
	ret0, _ := ret[0].(watcher.StringsChannel)
	return ret0
}

// Changes indicates an expected call of Changes
func (mr *MockStringsWatcherMockRecorder) Changes() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Changes", reflect.TypeOf((*MockStringsWatcher)(nil).Changes))
}

// Kill mocks base method
func (m *MockStringsWatcher) Kill() {
	m.ctrl.Call(m, "Kill")
}

// Kill indicates an expected call of Kill
func (mr *MockStringsWatcherMockRecorder) Kill() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kill", reflect.TypeOf((*MockStringsWatcher)(nil).Kill))
}

// Wait mocks base method
func (m *MockStringsWatcher) Wait() error {
	ret := m.ctrl.Call(m, "Wait")
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait
func (mr *MockStringsWatcherMockRecorder) Wait() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockStringsWatcher)(nil).Wait))
}
