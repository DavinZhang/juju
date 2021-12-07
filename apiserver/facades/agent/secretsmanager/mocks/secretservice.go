// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/secrets (interfaces: SecretsService)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	secrets "github.com/DavinZhang/juju/core/secrets"
	secrets0 "github.com/DavinZhang/juju/secrets"
)

// MockSecretsService is a mock of SecretsService interface.
type MockSecretsService struct {
	ctrl     *gomock.Controller
	recorder *MockSecretsServiceMockRecorder
}

// MockSecretsServiceMockRecorder is the mock recorder for MockSecretsService.
type MockSecretsServiceMockRecorder struct {
	mock *MockSecretsService
}

// NewMockSecretsService creates a new mock instance.
func NewMockSecretsService(ctrl *gomock.Controller) *MockSecretsService {
	mock := &MockSecretsService{ctrl: ctrl}
	mock.recorder = &MockSecretsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSecretsService) EXPECT() *MockSecretsServiceMockRecorder {
	return m.recorder
}

// CreateSecret mocks base method.
func (m *MockSecretsService) CreateSecret(arg0 context.Context, arg1 *secrets.URL, arg2 secrets0.CreateParams) (*secrets.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSecret", arg0, arg1, arg2)
	ret0, _ := ret[0].(*secrets.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSecret indicates an expected call of CreateSecret.
func (mr *MockSecretsServiceMockRecorder) CreateSecret(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSecret", reflect.TypeOf((*MockSecretsService)(nil).CreateSecret), arg0, arg1, arg2)
}

// DeleteSecret mocks base method.
func (m *MockSecretsService) DeleteSecret(arg0 context.Context, arg1 *secrets.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecret", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSecret indicates an expected call of DeleteSecret.
func (mr *MockSecretsServiceMockRecorder) DeleteSecret(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecret", reflect.TypeOf((*MockSecretsService)(nil).DeleteSecret), arg0, arg1)
}

// GetSecret mocks base method.
func (m *MockSecretsService) GetSecret(arg0 context.Context, arg1 *secrets.URL) (*secrets.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", arg0, arg1)
	ret0, _ := ret[0].(*secrets.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockSecretsServiceMockRecorder) GetSecret(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockSecretsService)(nil).GetSecret), arg0, arg1)
}

// GetSecretValue mocks base method.
func (m *MockSecretsService) GetSecretValue(arg0 context.Context, arg1 *secrets.URL) (secrets.SecretValue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecretValue", arg0, arg1)
	ret0, _ := ret[0].(secrets.SecretValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecretValue indicates an expected call of GetSecretValue.
func (mr *MockSecretsServiceMockRecorder) GetSecretValue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecretValue", reflect.TypeOf((*MockSecretsService)(nil).GetSecretValue), arg0, arg1)
}

// ListSecrets mocks base method.
func (m *MockSecretsService) ListSecrets(arg0 context.Context, arg1 secrets0.Filter) ([]*secrets.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSecrets", arg0, arg1)
	ret0, _ := ret[0].([]*secrets.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSecrets indicates an expected call of ListSecrets.
func (mr *MockSecretsServiceMockRecorder) ListSecrets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSecrets", reflect.TypeOf((*MockSecretsService)(nil).ListSecrets), arg0, arg1)
}

// UpdateSecret mocks base method.
func (m *MockSecretsService) UpdateSecret(arg0 context.Context, arg1 *secrets.URL, arg2 secrets0.UpdateParams) (*secrets.SecretMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSecret", arg0, arg1, arg2)
	ret0, _ := ret[0].(*secrets.SecretMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSecret indicates an expected call of UpdateSecret.
func (mr *MockSecretsServiceMockRecorder) UpdateSecret(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSecret", reflect.TypeOf((*MockSecretsService)(nil).UpdateSecret), arg0, arg1, arg2)
}
