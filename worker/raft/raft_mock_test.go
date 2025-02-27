// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/worker/raft (interfaces: Raft,ApplierMetrics)

// Package raft is a generated GoMock package.
package raft

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	raft "github.com/hashicorp/raft"
)

// MockRaft is a mock of Raft interface.
type MockRaft struct {
	ctrl     *gomock.Controller
	recorder *MockRaftMockRecorder
}

// MockRaftMockRecorder is the mock recorder for MockRaft.
type MockRaftMockRecorder struct {
	mock *MockRaft
}

// NewMockRaft creates a new mock instance.
func NewMockRaft(ctrl *gomock.Controller) *MockRaft {
	mock := &MockRaft{ctrl: ctrl}
	mock.recorder = &MockRaftMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRaft) EXPECT() *MockRaftMockRecorder {
	return m.recorder
}

// Apply mocks base method.
func (m *MockRaft) Apply(arg0 []byte, arg1 time.Duration) raft.ApplyFuture {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Apply", arg0, arg1)
	ret0, _ := ret[0].(raft.ApplyFuture)
	return ret0
}

// Apply indicates an expected call of Apply.
func (mr *MockRaftMockRecorder) Apply(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Apply", reflect.TypeOf((*MockRaft)(nil).Apply), arg0, arg1)
}

// GetConfiguration mocks base method.
func (m *MockRaft) GetConfiguration() raft.ConfigurationFuture {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfiguration")
	ret0, _ := ret[0].(raft.ConfigurationFuture)
	return ret0
}

// GetConfiguration indicates an expected call of GetConfiguration.
func (mr *MockRaftMockRecorder) GetConfiguration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfiguration", reflect.TypeOf((*MockRaft)(nil).GetConfiguration))
}

// Leader mocks base method.
func (m *MockRaft) Leader() raft.ServerAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Leader")
	ret0, _ := ret[0].(raft.ServerAddress)
	return ret0
}

// Leader indicates an expected call of Leader.
func (mr *MockRaftMockRecorder) Leader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Leader", reflect.TypeOf((*MockRaft)(nil).Leader))
}

// State mocks base method.
func (m *MockRaft) State() raft.RaftState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(raft.RaftState)
	return ret0
}

// State indicates an expected call of State.
func (mr *MockRaftMockRecorder) State() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockRaft)(nil).State))
}

// MockApplierMetrics is a mock of ApplierMetrics interface.
type MockApplierMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockApplierMetricsMockRecorder
}

// MockApplierMetricsMockRecorder is the mock recorder for MockApplierMetrics.
type MockApplierMetricsMockRecorder struct {
	mock *MockApplierMetrics
}

// NewMockApplierMetrics creates a new mock instance.
func NewMockApplierMetrics(ctrl *gomock.Controller) *MockApplierMetrics {
	mock := &MockApplierMetrics{ctrl: ctrl}
	mock.recorder = &MockApplierMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplierMetrics) EXPECT() *MockApplierMetricsMockRecorder {
	return m.recorder
}

// Record mocks base method.
func (m *MockApplierMetrics) Record(arg0 time.Time, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Record", arg0, arg1)
}

// Record indicates an expected call of Record.
func (mr *MockApplierMetricsMockRecorder) Record(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Record", reflect.TypeOf((*MockApplierMetrics)(nil).Record), arg0, arg1)
}

// RecordLeaderError mocks base method.
func (m *MockApplierMetrics) RecordLeaderError(arg0 time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RecordLeaderError", arg0)
}

// RecordLeaderError indicates an expected call of RecordLeaderError.
func (mr *MockApplierMetricsMockRecorder) RecordLeaderError(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecordLeaderError", reflect.TypeOf((*MockApplierMetrics)(nil).RecordLeaderError), arg0)
}
