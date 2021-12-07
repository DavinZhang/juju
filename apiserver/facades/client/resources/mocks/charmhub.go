// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/apiserver/facades/client/resources (interfaces: CharmHub)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	charmhub "github.com/DavinZhang/juju/charmhub"
	transport "github.com/DavinZhang/juju/charmhub/transport"
)

// MockCharmHub is a mock of CharmHub interface.
type MockCharmHub struct {
	ctrl     *gomock.Controller
	recorder *MockCharmHubMockRecorder
}

// MockCharmHubMockRecorder is the mock recorder for MockCharmHub.
type MockCharmHubMockRecorder struct {
	mock *MockCharmHub
}

// NewMockCharmHub creates a new mock instance.
func NewMockCharmHub(ctrl *gomock.Controller) *MockCharmHub {
	mock := &MockCharmHub{ctrl: ctrl}
	mock.recorder = &MockCharmHubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCharmHub) EXPECT() *MockCharmHubMockRecorder {
	return m.recorder
}

// ListResourceRevisions mocks base method.
func (m *MockCharmHub) ListResourceRevisions(arg0 context.Context, arg1, arg2 string) ([]transport.ResourceRevision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListResourceRevisions", arg0, arg1, arg2)
	ret0, _ := ret[0].([]transport.ResourceRevision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListResourceRevisions indicates an expected call of ListResourceRevisions.
func (mr *MockCharmHubMockRecorder) ListResourceRevisions(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListResourceRevisions", reflect.TypeOf((*MockCharmHub)(nil).ListResourceRevisions), arg0, arg1, arg2)
}

// Refresh mocks base method.
func (m *MockCharmHub) Refresh(arg0 context.Context, arg1 charmhub.RefreshConfig) ([]transport.RefreshResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", arg0, arg1)
	ret0, _ := ret[0].([]transport.RefreshResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh.
func (mr *MockCharmHubMockRecorder) Refresh(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockCharmHub)(nil).Refresh), arg0, arg1)
}
