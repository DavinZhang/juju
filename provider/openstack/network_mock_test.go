// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/provider/openstack (interfaces: SSLHostnameConfig,Networking,NetworkingBase,NetworkingNeutron,NetworkingAuthenticatingClient,NetworkingNova,NetworkingEnvironConfig)

// Package openstack is a generated GoMock package.
package openstack

import (
	reflect "reflect"

	neutron "github.com/go-goose/goose/v4/neutron"
	nova "github.com/go-goose/goose/v4/nova"
	gomock "github.com/golang/mock/gomock"
	set "github.com/juju/collections/set"
	instance "github.com/DavinZhang/juju/core/instance"
	network "github.com/DavinZhang/juju/core/network"
)

// MockSSLHostnameConfig is a mock of SSLHostnameConfig interface.
type MockSSLHostnameConfig struct {
	ctrl     *gomock.Controller
	recorder *MockSSLHostnameConfigMockRecorder
}

// MockSSLHostnameConfigMockRecorder is the mock recorder for MockSSLHostnameConfig.
type MockSSLHostnameConfigMockRecorder struct {
	mock *MockSSLHostnameConfig
}

// NewMockSSLHostnameConfig creates a new mock instance.
func NewMockSSLHostnameConfig(ctrl *gomock.Controller) *MockSSLHostnameConfig {
	mock := &MockSSLHostnameConfig{ctrl: ctrl}
	mock.recorder = &MockSSLHostnameConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSSLHostnameConfig) EXPECT() *MockSSLHostnameConfigMockRecorder {
	return m.recorder
}

// SSLHostnameVerification mocks base method.
func (m *MockSSLHostnameConfig) SSLHostnameVerification() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SSLHostnameVerification")
	ret0, _ := ret[0].(bool)
	return ret0
}

// SSLHostnameVerification indicates an expected call of SSLHostnameVerification.
func (mr *MockSSLHostnameConfigMockRecorder) SSLHostnameVerification() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SSLHostnameVerification", reflect.TypeOf((*MockSSLHostnameConfig)(nil).SSLHostnameVerification))
}

// MockNetworking is a mock of Networking interface.
type MockNetworking struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkingMockRecorder
}

// MockNetworkingMockRecorder is the mock recorder for MockNetworking.
type MockNetworkingMockRecorder struct {
	mock *MockNetworking
}

// NewMockNetworking creates a new mock instance.
func NewMockNetworking(ctrl *gomock.Controller) *MockNetworking {
	mock := &MockNetworking{ctrl: ctrl}
	mock.recorder = &MockNetworkingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNetworking) EXPECT() *MockNetworkingMockRecorder {
	return m.recorder
}

// AllocatePublicIP mocks base method.
func (m *MockNetworking) AllocatePublicIP(arg0 instance.Id) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllocatePublicIP", arg0)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllocatePublicIP indicates an expected call of AllocatePublicIP.
func (mr *MockNetworkingMockRecorder) AllocatePublicIP(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllocatePublicIP", reflect.TypeOf((*MockNetworking)(nil).AllocatePublicIP), arg0)
}

// CreatePort mocks base method.
func (m *MockNetworking) CreatePort(arg0, arg1 string, arg2 network.Id) (*neutron.PortV2, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePort", arg0, arg1, arg2)
	ret0, _ := ret[0].(*neutron.PortV2)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePort indicates an expected call of CreatePort.
func (mr *MockNetworkingMockRecorder) CreatePort(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePort", reflect.TypeOf((*MockNetworking)(nil).CreatePort), arg0, arg1, arg2)
}

// DefaultNetworks mocks base method.
func (m *MockNetworking) DefaultNetworks() ([]nova.ServerNetworks, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultNetworks")
	ret0, _ := ret[0].([]nova.ServerNetworks)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DefaultNetworks indicates an expected call of DefaultNetworks.
func (mr *MockNetworkingMockRecorder) DefaultNetworks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultNetworks", reflect.TypeOf((*MockNetworking)(nil).DefaultNetworks))
}

// DeletePortByID mocks base method.
func (m *MockNetworking) DeletePortByID(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePortByID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePortByID indicates an expected call of DeletePortByID.
func (mr *MockNetworkingMockRecorder) DeletePortByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePortByID", reflect.TypeOf((*MockNetworking)(nil).DeletePortByID), arg0)
}

// FindNetworks mocks base method.
func (m *MockNetworking) FindNetworks(arg0 bool) (set.Strings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindNetworks", arg0)
	ret0, _ := ret[0].(set.Strings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindNetworks indicates an expected call of FindNetworks.
func (mr *MockNetworkingMockRecorder) FindNetworks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindNetworks", reflect.TypeOf((*MockNetworking)(nil).FindNetworks), arg0)
}

// NetworkInterfaces mocks base method.
func (m *MockNetworking) NetworkInterfaces(arg0 []instance.Id) ([]network.InterfaceInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NetworkInterfaces", arg0)
	ret0, _ := ret[0].([]network.InterfaceInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NetworkInterfaces indicates an expected call of NetworkInterfaces.
func (mr *MockNetworkingMockRecorder) NetworkInterfaces(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NetworkInterfaces", reflect.TypeOf((*MockNetworking)(nil).NetworkInterfaces), arg0)
}

// ResolveNetwork mocks base method.
func (m *MockNetworking) ResolveNetwork(arg0 string, arg1 bool) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveNetwork", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResolveNetwork indicates an expected call of ResolveNetwork.
func (mr *MockNetworkingMockRecorder) ResolveNetwork(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveNetwork", reflect.TypeOf((*MockNetworking)(nil).ResolveNetwork), arg0, arg1)
}

// Subnets mocks base method.
func (m *MockNetworking) Subnets(arg0 instance.Id, arg1 []network.Id) ([]network.SubnetInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subnets", arg0, arg1)
	ret0, _ := ret[0].([]network.SubnetInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subnets indicates an expected call of Subnets.
func (mr *MockNetworkingMockRecorder) Subnets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subnets", reflect.TypeOf((*MockNetworking)(nil).Subnets), arg0, arg1)
}

// MockNetworkingBase is a mock of NetworkingBase interface
type MockNetworkingBase struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkingBaseMockRecorder
}

// MockNetworkingBaseMockRecorder is the mock recorder for MockNetworkingBase
type MockNetworkingBaseMockRecorder struct {
	mock *MockNetworkingBase
}

// NewMockNetworkingBase creates a new mock instance
func NewMockNetworkingBase(ctrl *gomock.Controller) *MockNetworkingBase {
	mock := &MockNetworkingBase{ctrl: ctrl}
	mock.recorder = &MockNetworkingBaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNetworkingBase) EXPECT() *MockNetworkingBaseMockRecorder {
	return m.recorder
}

// client mocks base method
func (m *MockNetworkingBase) client() NetworkingAuthenticatingClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "client")
	ret0, _ := ret[0].(NetworkingAuthenticatingClient)
	return ret0
}

// client indicates an expected call of client
func (mr *MockNetworkingBaseMockRecorder) client() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "client", reflect.TypeOf((*MockNetworkingBase)(nil).client))
}

// ecfg mocks base method
func (m *MockNetworkingBase) ecfg() NetworkingEnvironConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ecfg")
	ret0, _ := ret[0].(NetworkingEnvironConfig)
	return ret0
}

// ecfg indicates an expected call of ecfg
func (mr *MockNetworkingBaseMockRecorder) ecfg() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ecfg", reflect.TypeOf((*MockNetworkingBase)(nil).ecfg))
}

// neutron mocks base method
func (m *MockNetworkingBase) neutron() NetworkingNeutron {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "neutron")
	ret0, _ := ret[0].(NetworkingNeutron)
	return ret0
}

// neutron indicates an expected call of neutron
func (mr *MockNetworkingBaseMockRecorder) neutron() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "neutron", reflect.TypeOf((*MockNetworkingBase)(nil).neutron))
}

// nova mocks base method
func (m *MockNetworkingBase) nova() NetworkingNova {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "nova")
	ret0, _ := ret[0].(NetworkingNova)
	return ret0
}

// nova indicates an expected call of nova
func (mr *MockNetworkingBaseMockRecorder) nova() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "nova", reflect.TypeOf((*MockNetworkingBase)(nil).nova))
}

// MockNetworkingNeutron is a mock of NetworkingNeutron interface
type MockNetworkingNeutron struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkingNeutronMockRecorder
}

// MockNetworkingNeutronMockRecorder is the mock recorder for MockNetworkingNeutron
type MockNetworkingNeutronMockRecorder struct {
	mock *MockNetworkingNeutron
}

// NewMockNetworkingNeutron creates a new mock instance
func NewMockNetworkingNeutron(ctrl *gomock.Controller) *MockNetworkingNeutron {
	mock := &MockNetworkingNeutron{ctrl: ctrl}
	mock.recorder = &MockNetworkingNeutronMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNetworkingNeutron) EXPECT() *MockNetworkingNeutronMockRecorder {
	return m.recorder
}

// AllocateFloatingIPV2 mocks base method
func (m *MockNetworkingNeutron) AllocateFloatingIPV2(arg0 string) (*neutron.FloatingIPV2, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllocateFloatingIPV2", arg0)
	ret0, _ := ret[0].(*neutron.FloatingIPV2)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllocateFloatingIPV2 indicates an expected call of AllocateFloatingIPV2
func (mr *MockNetworkingNeutronMockRecorder) AllocateFloatingIPV2(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllocateFloatingIPV2", reflect.TypeOf((*MockNetworkingNeutron)(nil).AllocateFloatingIPV2), arg0)
}

// CreatePortV2 mocks base method
func (m *MockNetworkingNeutron) CreatePortV2(arg0 neutron.PortV2) (*neutron.PortV2, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePortV2", arg0)
	ret0, _ := ret[0].(*neutron.PortV2)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePortV2 indicates an expected call of CreatePortV2
func (mr *MockNetworkingNeutronMockRecorder) CreatePortV2(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePortV2", reflect.TypeOf((*MockNetworkingNeutron)(nil).CreatePortV2), arg0)
}

// DeletePortV2 mocks base method
func (m *MockNetworkingNeutron) DeletePortV2(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePortV2", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePortV2 indicates an expected call of DeletePortV2
func (mr *MockNetworkingNeutronMockRecorder) DeletePortV2(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePortV2", reflect.TypeOf((*MockNetworkingNeutron)(nil).DeletePortV2), arg0)
}

// GetNetworkV2 mocks base method
func (m *MockNetworkingNeutron) GetNetworkV2(arg0 string) (*neutron.NetworkV2, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetworkV2", arg0)
	ret0, _ := ret[0].(*neutron.NetworkV2)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNetworkV2 indicates an expected call of GetNetworkV2
func (mr *MockNetworkingNeutronMockRecorder) GetNetworkV2(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkV2", reflect.TypeOf((*MockNetworkingNeutron)(nil).GetNetworkV2), arg0)
}

// ListFloatingIPsV2 mocks base method
func (m *MockNetworkingNeutron) ListFloatingIPsV2(arg0 ...*neutron.Filter) ([]neutron.FloatingIPV2, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListFloatingIPsV2", varargs...)
	ret0, _ := ret[0].([]neutron.FloatingIPV2)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFloatingIPsV2 indicates an expected call of ListFloatingIPsV2
func (mr *MockNetworkingNeutronMockRecorder) ListFloatingIPsV2(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFloatingIPsV2", reflect.TypeOf((*MockNetworkingNeutron)(nil).ListFloatingIPsV2), arg0...)
}

// ListNetworksV2 mocks base method
func (m *MockNetworkingNeutron) ListNetworksV2(arg0 ...*neutron.Filter) ([]neutron.NetworkV2, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListNetworksV2", varargs...)
	ret0, _ := ret[0].([]neutron.NetworkV2)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNetworksV2 indicates an expected call of ListNetworksV2
func (mr *MockNetworkingNeutronMockRecorder) ListNetworksV2(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNetworksV2", reflect.TypeOf((*MockNetworkingNeutron)(nil).ListNetworksV2), arg0...)
}

// ListSubnetsV2 mocks base method
func (m *MockNetworkingNeutron) ListSubnetsV2() ([]neutron.SubnetV2, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSubnetsV2")
	ret0, _ := ret[0].([]neutron.SubnetV2)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSubnetsV2 indicates an expected call of ListSubnetsV2
func (mr *MockNetworkingNeutronMockRecorder) ListSubnetsV2() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSubnetsV2", reflect.TypeOf((*MockNetworkingNeutron)(nil).ListSubnetsV2))
}

// MockNetworkingAuthenticatingClient is a mock of NetworkingAuthenticatingClient interface
type MockNetworkingAuthenticatingClient struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkingAuthenticatingClientMockRecorder
}

// MockNetworkingAuthenticatingClientMockRecorder is the mock recorder for MockNetworkingAuthenticatingClient
type MockNetworkingAuthenticatingClientMockRecorder struct {
	mock *MockNetworkingAuthenticatingClient
}

// NewMockNetworkingAuthenticatingClient creates a new mock instance
func NewMockNetworkingAuthenticatingClient(ctrl *gomock.Controller) *MockNetworkingAuthenticatingClient {
	mock := &MockNetworkingAuthenticatingClient{ctrl: ctrl}
	mock.recorder = &MockNetworkingAuthenticatingClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNetworkingAuthenticatingClient) EXPECT() *MockNetworkingAuthenticatingClientMockRecorder {
	return m.recorder
}

// TenantId mocks base method
func (m *MockNetworkingAuthenticatingClient) TenantId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TenantId")
	ret0, _ := ret[0].(string)
	return ret0
}

// TenantId indicates an expected call of TenantId
func (mr *MockNetworkingAuthenticatingClientMockRecorder) TenantId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TenantId", reflect.TypeOf((*MockNetworkingAuthenticatingClient)(nil).TenantId))
}

// MockNetworkingNova is a mock of NetworkingNova interface
type MockNetworkingNova struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkingNovaMockRecorder
}

// MockNetworkingNovaMockRecorder is the mock recorder for MockNetworkingNova
type MockNetworkingNovaMockRecorder struct {
	mock *MockNetworkingNova
}

// NewMockNetworkingNova creates a new mock instance
func NewMockNetworkingNova(ctrl *gomock.Controller) *MockNetworkingNova {
	mock := &MockNetworkingNova{ctrl: ctrl}
	mock.recorder = &MockNetworkingNovaMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNetworkingNova) EXPECT() *MockNetworkingNovaMockRecorder {
	return m.recorder
}

// GetServer mocks base method
func (m *MockNetworkingNova) GetServer(arg0 string) (*nova.ServerDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServer", arg0)
	ret0, _ := ret[0].(*nova.ServerDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServer indicates an expected call of GetServer
func (mr *MockNetworkingNovaMockRecorder) GetServer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServer", reflect.TypeOf((*MockNetworkingNova)(nil).GetServer), arg0)
}

// MockNetworkingEnvironConfig is a mock of NetworkingEnvironConfig interface
type MockNetworkingEnvironConfig struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkingEnvironConfigMockRecorder
}

// MockNetworkingEnvironConfigMockRecorder is the mock recorder for MockNetworkingEnvironConfig
type MockNetworkingEnvironConfigMockRecorder struct {
	mock *MockNetworkingEnvironConfig
}

// NewMockNetworkingEnvironConfig creates a new mock instance
func NewMockNetworkingEnvironConfig(ctrl *gomock.Controller) *MockNetworkingEnvironConfig {
	mock := &MockNetworkingEnvironConfig{ctrl: ctrl}
	mock.recorder = &MockNetworkingEnvironConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNetworkingEnvironConfig) EXPECT() *MockNetworkingEnvironConfigMockRecorder {
	return m.recorder
}

// externalNetwork mocks base method
func (m *MockNetworkingEnvironConfig) externalNetwork() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "externalNetwork")
	ret0, _ := ret[0].(string)
	return ret0
}

// externalNetwork indicates an expected call of externalNetwork
func (mr *MockNetworkingEnvironConfigMockRecorder) externalNetwork() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "externalNetwork", reflect.TypeOf((*MockNetworkingEnvironConfig)(nil).externalNetwork))
}

// network mocks base method
func (m *MockNetworkingEnvironConfig) network() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "network")
	ret0, _ := ret[0].(string)
	return ret0
}

// network indicates an expected call of network
func (mr *MockNetworkingEnvironConfigMockRecorder) network() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "network", reflect.TypeOf((*MockNetworkingEnvironConfig)(nil).network))
}
