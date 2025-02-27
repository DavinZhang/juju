// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/apiserver/facades/client/subnets (interfaces: Backing)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	networkingcommon "github.com/DavinZhang/juju/apiserver/common/networkingcommon"
	network "github.com/DavinZhang/juju/core/network"
	cloudspec "github.com/DavinZhang/juju/environs/cloudspec"
	config "github.com/DavinZhang/juju/environs/config"
	names "github.com/juju/names/v4"
	reflect "reflect"
)

// MockBacking is a mock of Backing interface
type MockBacking struct {
	ctrl     *gomock.Controller
	recorder *MockBackingMockRecorder
}

// MockBackingMockRecorder is the mock recorder for MockBacking
type MockBackingMockRecorder struct {
	mock *MockBacking
}

// NewMockBacking creates a new mock instance
func NewMockBacking(ctrl *gomock.Controller) *MockBacking {
	mock := &MockBacking{ctrl: ctrl}
	mock.recorder = &MockBackingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBacking) EXPECT() *MockBackingMockRecorder {
	return m.recorder
}

// AddSubnet mocks base method
func (m *MockBacking) AddSubnet(arg0 networkingcommon.BackingSubnetInfo) (networkingcommon.BackingSubnet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSubnet", arg0)
	ret0, _ := ret[0].(networkingcommon.BackingSubnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSubnet indicates an expected call of AddSubnet
func (mr *MockBackingMockRecorder) AddSubnet(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSubnet", reflect.TypeOf((*MockBacking)(nil).AddSubnet), arg0)
}

// AllSpaces mocks base method
func (m *MockBacking) AllSpaces() ([]networkingcommon.BackingSpace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllSpaces")
	ret0, _ := ret[0].([]networkingcommon.BackingSpace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllSpaces indicates an expected call of AllSpaces
func (mr *MockBackingMockRecorder) AllSpaces() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllSpaces", reflect.TypeOf((*MockBacking)(nil).AllSpaces))
}

// AllSubnets mocks base method
func (m *MockBacking) AllSubnets() ([]networkingcommon.BackingSubnet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllSubnets")
	ret0, _ := ret[0].([]networkingcommon.BackingSubnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllSubnets indicates an expected call of AllSubnets
func (mr *MockBackingMockRecorder) AllSubnets() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllSubnets", reflect.TypeOf((*MockBacking)(nil).AllSubnets))
}

// AvailabilityZones mocks base method
func (m *MockBacking) AvailabilityZones() (network.AvailabilityZones, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AvailabilityZones")
	ret0, _ := ret[0].(network.AvailabilityZones)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AvailabilityZones indicates an expected call of AvailabilityZones
func (mr *MockBackingMockRecorder) AvailabilityZones() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AvailabilityZones", reflect.TypeOf((*MockBacking)(nil).AvailabilityZones))
}

// CloudSpec mocks base method
func (m *MockBacking) CloudSpec() (cloudspec.CloudSpec, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudSpec")
	ret0, _ := ret[0].(cloudspec.CloudSpec)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloudSpec indicates an expected call of CloudSpec
func (mr *MockBackingMockRecorder) CloudSpec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudSpec", reflect.TypeOf((*MockBacking)(nil).CloudSpec))
}

// ModelConfig mocks base method
func (m *MockBacking) ModelConfig() (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelConfig")
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelConfig indicates an expected call of ModelConfig
func (mr *MockBackingMockRecorder) ModelConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelConfig", reflect.TypeOf((*MockBacking)(nil).ModelConfig))
}

// ModelTag mocks base method
func (m *MockBacking) ModelTag() names.ModelTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelTag")
	ret0, _ := ret[0].(names.ModelTag)
	return ret0
}

// ModelTag indicates an expected call of ModelTag
func (mr *MockBackingMockRecorder) ModelTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelTag", reflect.TypeOf((*MockBacking)(nil).ModelTag))
}

// SetAvailabilityZones mocks base method
func (m *MockBacking) SetAvailabilityZones(arg0 network.AvailabilityZones) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAvailabilityZones", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAvailabilityZones indicates an expected call of SetAvailabilityZones
func (mr *MockBackingMockRecorder) SetAvailabilityZones(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAvailabilityZones", reflect.TypeOf((*MockBacking)(nil).SetAvailabilityZones), arg0)
}

// SubnetByCIDR mocks base method
func (m *MockBacking) SubnetByCIDR(arg0 string) (networkingcommon.BackingSubnet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubnetByCIDR", arg0)
	ret0, _ := ret[0].(networkingcommon.BackingSubnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubnetByCIDR indicates an expected call of SubnetByCIDR
func (mr *MockBackingMockRecorder) SubnetByCIDR(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubnetByCIDR", reflect.TypeOf((*MockBacking)(nil).SubnetByCIDR), arg0)
}

// SubnetsByCIDR mocks base method
func (m *MockBacking) SubnetsByCIDR(arg0 string) ([]networkingcommon.BackingSubnet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubnetsByCIDR", arg0)
	ret0, _ := ret[0].([]networkingcommon.BackingSubnet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubnetsByCIDR indicates an expected call of SubnetsByCIDR
func (mr *MockBackingMockRecorder) SubnetsByCIDR(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubnetsByCIDR", reflect.TypeOf((*MockBacking)(nil).SubnetsByCIDR), arg0)
}
