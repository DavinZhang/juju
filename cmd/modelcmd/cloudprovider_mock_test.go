// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/cmd/modelcmd (interfaces: TestCloudProvider)

// Package modelcmd is a generated GoMock package.
package modelcmd

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	jsonschema "github.com/juju/jsonschema"
	cloud "github.com/DavinZhang/juju/cloud"
	environs "github.com/DavinZhang/juju/environs"
	config "github.com/DavinZhang/juju/environs/config"
	context "github.com/DavinZhang/juju/environs/context"
)

// MockTestCloudProvider is a mock of TestCloudProvider interface.
type MockTestCloudProvider struct {
	ctrl     *gomock.Controller
	recorder *MockTestCloudProviderMockRecorder
}

// MockTestCloudProviderMockRecorder is the mock recorder for MockTestCloudProvider.
type MockTestCloudProviderMockRecorder struct {
	mock *MockTestCloudProvider
}

// NewMockTestCloudProvider creates a new mock instance.
func NewMockTestCloudProvider(ctrl *gomock.Controller) *MockTestCloudProvider {
	mock := &MockTestCloudProvider{ctrl: ctrl}
	mock.recorder = &MockTestCloudProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTestCloudProvider) EXPECT() *MockTestCloudProviderMockRecorder {
	return m.recorder
}

// CloudSchema mocks base method.
func (m *MockTestCloudProvider) CloudSchema() *jsonschema.Schema {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudSchema")
	ret0, _ := ret[0].(*jsonschema.Schema)
	return ret0
}

// CloudSchema indicates an expected call of CloudSchema.
func (mr *MockTestCloudProviderMockRecorder) CloudSchema() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudSchema", reflect.TypeOf((*MockTestCloudProvider)(nil).CloudSchema))
}

// CredentialSchemas mocks base method.
func (m *MockTestCloudProvider) CredentialSchemas() map[cloud.AuthType]cloud.CredentialSchema {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CredentialSchemas")
	ret0, _ := ret[0].(map[cloud.AuthType]cloud.CredentialSchema)
	return ret0
}

// CredentialSchemas indicates an expected call of CredentialSchemas.
func (mr *MockTestCloudProviderMockRecorder) CredentialSchemas() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CredentialSchemas", reflect.TypeOf((*MockTestCloudProvider)(nil).CredentialSchemas))
}

// DetectCredentials mocks base method.
func (m *MockTestCloudProvider) DetectCredentials(arg0 string) (*cloud.CloudCredential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetectCredentials", arg0)
	ret0, _ := ret[0].(*cloud.CloudCredential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DetectCredentials indicates an expected call of DetectCredentials.
func (mr *MockTestCloudProviderMockRecorder) DetectCredentials(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetectCredentials", reflect.TypeOf((*MockTestCloudProvider)(nil).DetectCredentials), arg0)
}

// FinalizeCredential mocks base method.
func (m *MockTestCloudProvider) FinalizeCredential(arg0 environs.FinalizeCredentialContext, arg1 environs.FinalizeCredentialParams) (*cloud.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinalizeCredential", arg0, arg1)
	ret0, _ := ret[0].(*cloud.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FinalizeCredential indicates an expected call of FinalizeCredential.
func (mr *MockTestCloudProviderMockRecorder) FinalizeCredential(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinalizeCredential", reflect.TypeOf((*MockTestCloudProvider)(nil).FinalizeCredential), arg0, arg1)
}

// Ping mocks base method.
func (m *MockTestCloudProvider) Ping(arg0 context.ProviderCallContext, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockTestCloudProviderMockRecorder) Ping(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockTestCloudProvider)(nil).Ping), arg0, arg1)
}

// PrepareConfig mocks base method.
func (m *MockTestCloudProvider) PrepareConfig(arg0 environs.PrepareConfigParams) (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareConfig", arg0)
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareConfig indicates an expected call of PrepareConfig.
func (mr *MockTestCloudProviderMockRecorder) PrepareConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareConfig", reflect.TypeOf((*MockTestCloudProvider)(nil).PrepareConfig), arg0)
}

// RegisterCredentials mocks base method.
func (m *MockTestCloudProvider) RegisterCredentials(arg0 cloud.Cloud) (map[string]*cloud.CloudCredential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterCredentials", arg0)
	ret0, _ := ret[0].(map[string]*cloud.CloudCredential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterCredentials indicates an expected call of RegisterCredentials.
func (mr *MockTestCloudProviderMockRecorder) RegisterCredentials(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterCredentials", reflect.TypeOf((*MockTestCloudProvider)(nil).RegisterCredentials), arg0)
}

// Validate mocks base method.
func (m *MockTestCloudProvider) Validate(arg0, arg1 *config.Config) (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", arg0, arg1)
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Validate indicates an expected call of Validate.
func (mr *MockTestCloudProviderMockRecorder) Validate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockTestCloudProvider)(nil).Validate), arg0, arg1)
}

// Version mocks base method.
func (m *MockTestCloudProvider) Version() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Version")
	ret0, _ := ret[0].(int)
	return ret0
}

// Version indicates an expected call of Version.
func (mr *MockTestCloudProviderMockRecorder) Version() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Version", reflect.TypeOf((*MockTestCloudProvider)(nil).Version))
}
