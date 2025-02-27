// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/provider/oci (interfaces: IdentityClient)

// Package testing is a generated GoMock package.
package testing

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	identity "github.com/oracle/oci-go-sdk/v47/identity"
)

// MockIdentityClient is a mock of IdentityClient interface.
type MockIdentityClient struct {
	ctrl     *gomock.Controller
	recorder *MockIdentityClientMockRecorder
}

// MockIdentityClientMockRecorder is the mock recorder for MockIdentityClient.
type MockIdentityClientMockRecorder struct {
	mock *MockIdentityClient
}

// NewMockIdentityClient creates a new mock instance.
func NewMockIdentityClient(ctrl *gomock.Controller) *MockIdentityClient {
	mock := &MockIdentityClient{ctrl: ctrl}
	mock.recorder = &MockIdentityClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIdentityClient) EXPECT() *MockIdentityClientMockRecorder {
	return m.recorder
}

// ListAvailabilityDomains mocks base method.
func (m *MockIdentityClient) ListAvailabilityDomains(arg0 context.Context, arg1 identity.ListAvailabilityDomainsRequest) (identity.ListAvailabilityDomainsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAvailabilityDomains", arg0, arg1)
	ret0, _ := ret[0].(identity.ListAvailabilityDomainsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAvailabilityDomains indicates an expected call of ListAvailabilityDomains.
func (mr *MockIdentityClientMockRecorder) ListAvailabilityDomains(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAvailabilityDomains", reflect.TypeOf((*MockIdentityClient)(nil).ListAvailabilityDomains), arg0, arg1)
}
