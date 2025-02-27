// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/DavinZhang/juju/apiserver/facades/client/machinemanager (interfaces: Machine,Application,Unit,Charm,CharmhubClient)

// Package machinemanager is a generated GoMock package.
package machinemanager

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	v9 "github.com/juju/charm/v9"
	charmhub "github.com/DavinZhang/juju/charmhub"
	transport "github.com/DavinZhang/juju/charmhub/transport"
	model "github.com/DavinZhang/juju/core/model"
	status "github.com/DavinZhang/juju/core/status"
	state "github.com/DavinZhang/juju/state"
	v4 "github.com/juju/names/v4"
	reflect "reflect"
	time "time"
)

// MockMachine is a mock of Machine interface
type MockMachine struct {
	ctrl     *gomock.Controller
	recorder *MockMachineMockRecorder
}

// MockMachineMockRecorder is the mock recorder for MockMachine
type MockMachineMockRecorder struct {
	mock *MockMachine
}

// NewMockMachine creates a new mock instance
func NewMockMachine(ctrl *gomock.Controller) *MockMachine {
	mock := &MockMachine{ctrl: ctrl}
	mock.recorder = &MockMachineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMachine) EXPECT() *MockMachineMockRecorder {
	return m.recorder
}

// ApplicationNames mocks base method
func (m *MockMachine) ApplicationNames() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplicationNames")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplicationNames indicates an expected call of ApplicationNames
func (mr *MockMachineMockRecorder) ApplicationNames() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplicationNames", reflect.TypeOf((*MockMachine)(nil).ApplicationNames))
}

// CompleteUpgradeSeries mocks base method
func (m *MockMachine) CompleteUpgradeSeries() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompleteUpgradeSeries")
	ret0, _ := ret[0].(error)
	return ret0
}

// CompleteUpgradeSeries indicates an expected call of CompleteUpgradeSeries
func (mr *MockMachineMockRecorder) CompleteUpgradeSeries() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteUpgradeSeries", reflect.TypeOf((*MockMachine)(nil).CompleteUpgradeSeries))
}

// CreateUpgradeSeriesLock mocks base method
func (m *MockMachine) CreateUpgradeSeriesLock(arg0 []string, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUpgradeSeriesLock", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUpgradeSeriesLock indicates an expected call of CreateUpgradeSeriesLock
func (mr *MockMachineMockRecorder) CreateUpgradeSeriesLock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUpgradeSeriesLock", reflect.TypeOf((*MockMachine)(nil).CreateUpgradeSeriesLock), arg0, arg1)
}

// Destroy mocks base method
func (m *MockMachine) Destroy() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy")
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy
func (mr *MockMachineMockRecorder) Destroy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockMachine)(nil).Destroy))
}

// ForceDestroy mocks base method
func (m *MockMachine) ForceDestroy(arg0 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForceDestroy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForceDestroy indicates an expected call of ForceDestroy
func (mr *MockMachineMockRecorder) ForceDestroy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForceDestroy", reflect.TypeOf((*MockMachine)(nil).ForceDestroy), arg0)
}

// GetUpgradeSeriesMessages mocks base method
func (m *MockMachine) GetUpgradeSeriesMessages() ([]string, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUpgradeSeriesMessages")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUpgradeSeriesMessages indicates an expected call of GetUpgradeSeriesMessages
func (mr *MockMachineMockRecorder) GetUpgradeSeriesMessages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpgradeSeriesMessages", reflect.TypeOf((*MockMachine)(nil).GetUpgradeSeriesMessages))
}

// Id mocks base method
func (m *MockMachine) Id() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(string)
	return ret0
}

// Id indicates an expected call of Id
func (mr *MockMachineMockRecorder) Id() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*MockMachine)(nil).Id))
}

// IsLockedForSeriesUpgrade mocks base method
func (m *MockMachine) IsLockedForSeriesUpgrade() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLockedForSeriesUpgrade")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsLockedForSeriesUpgrade indicates an expected call of IsLockedForSeriesUpgrade
func (mr *MockMachineMockRecorder) IsLockedForSeriesUpgrade() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLockedForSeriesUpgrade", reflect.TypeOf((*MockMachine)(nil).IsLockedForSeriesUpgrade))
}

// IsManager mocks base method
func (m *MockMachine) IsManager() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsManager")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsManager indicates an expected call of IsManager
func (mr *MockMachineMockRecorder) IsManager() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsManager", reflect.TypeOf((*MockMachine)(nil).IsManager))
}

// Principals mocks base method
func (m *MockMachine) Principals() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Principals")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Principals indicates an expected call of Principals
func (mr *MockMachineMockRecorder) Principals() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Principals", reflect.TypeOf((*MockMachine)(nil).Principals))
}

// RemoveUpgradeSeriesLock mocks base method
func (m *MockMachine) RemoveUpgradeSeriesLock() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUpgradeSeriesLock")
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUpgradeSeriesLock indicates an expected call of RemoveUpgradeSeriesLock
func (mr *MockMachineMockRecorder) RemoveUpgradeSeriesLock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUpgradeSeriesLock", reflect.TypeOf((*MockMachine)(nil).RemoveUpgradeSeriesLock))
}

// Series mocks base method
func (m *MockMachine) Series() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Series")
	ret0, _ := ret[0].(string)
	return ret0
}

// Series indicates an expected call of Series
func (mr *MockMachineMockRecorder) Series() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Series", reflect.TypeOf((*MockMachine)(nil).Series))
}

// SetKeepInstance mocks base method
func (m *MockMachine) SetKeepInstance(arg0 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetKeepInstance", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetKeepInstance indicates an expected call of SetKeepInstance
func (mr *MockMachineMockRecorder) SetKeepInstance(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetKeepInstance", reflect.TypeOf((*MockMachine)(nil).SetKeepInstance), arg0)
}

// SetUpgradeSeriesStatus mocks base method
func (m *MockMachine) SetUpgradeSeriesStatus(arg0 model.UpgradeSeriesStatus, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUpgradeSeriesStatus", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUpgradeSeriesStatus indicates an expected call of SetUpgradeSeriesStatus
func (mr *MockMachineMockRecorder) SetUpgradeSeriesStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUpgradeSeriesStatus", reflect.TypeOf((*MockMachine)(nil).SetUpgradeSeriesStatus), arg0, arg1)
}

// Tag mocks base method
func (m *MockMachine) Tag() v4.Tag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tag")
	ret0, _ := ret[0].(v4.Tag)
	return ret0
}

// Tag indicates an expected call of Tag
func (mr *MockMachineMockRecorder) Tag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tag", reflect.TypeOf((*MockMachine)(nil).Tag))
}

// Units mocks base method
func (m *MockMachine) Units() ([]Unit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Units")
	ret0, _ := ret[0].([]Unit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Units indicates an expected call of Units
func (mr *MockMachineMockRecorder) Units() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Units", reflect.TypeOf((*MockMachine)(nil).Units))
}

// UpgradeSeriesStatus mocks base method
func (m *MockMachine) UpgradeSeriesStatus() (model.UpgradeSeriesStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpgradeSeriesStatus")
	ret0, _ := ret[0].(model.UpgradeSeriesStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpgradeSeriesStatus indicates an expected call of UpgradeSeriesStatus
func (mr *MockMachineMockRecorder) UpgradeSeriesStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpgradeSeriesStatus", reflect.TypeOf((*MockMachine)(nil).UpgradeSeriesStatus))
}

// WatchUpgradeSeriesNotifications mocks base method
func (m *MockMachine) WatchUpgradeSeriesNotifications() (state.NotifyWatcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchUpgradeSeriesNotifications")
	ret0, _ := ret[0].(state.NotifyWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchUpgradeSeriesNotifications indicates an expected call of WatchUpgradeSeriesNotifications
func (mr *MockMachineMockRecorder) WatchUpgradeSeriesNotifications() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchUpgradeSeriesNotifications", reflect.TypeOf((*MockMachine)(nil).WatchUpgradeSeriesNotifications))
}

// MockApplication is a mock of Application interface
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// Charm mocks base method
func (m *MockApplication) Charm() (Charm, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Charm")
	ret0, _ := ret[0].(Charm)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Charm indicates an expected call of Charm
func (mr *MockApplicationMockRecorder) Charm() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Charm", reflect.TypeOf((*MockApplication)(nil).Charm))
}

// CharmOrigin mocks base method
func (m *MockApplication) CharmOrigin() *state.CharmOrigin {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CharmOrigin")
	ret0, _ := ret[0].(*state.CharmOrigin)
	return ret0
}

// CharmOrigin indicates an expected call of CharmOrigin
func (mr *MockApplicationMockRecorder) CharmOrigin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CharmOrigin", reflect.TypeOf((*MockApplication)(nil).CharmOrigin))
}

// Name mocks base method
func (m *MockApplication) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockApplicationMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockApplication)(nil).Name))
}

// MockUnit is a mock of Unit interface
type MockUnit struct {
	ctrl     *gomock.Controller
	recorder *MockUnitMockRecorder
}

// MockUnitMockRecorder is the mock recorder for MockUnit
type MockUnitMockRecorder struct {
	mock *MockUnit
}

// NewMockUnit creates a new mock instance
func NewMockUnit(ctrl *gomock.Controller) *MockUnit {
	mock := &MockUnit{ctrl: ctrl}
	mock.recorder = &MockUnitMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUnit) EXPECT() *MockUnitMockRecorder {
	return m.recorder
}

// AgentStatus mocks base method
func (m *MockUnit) AgentStatus() (status.StatusInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AgentStatus")
	ret0, _ := ret[0].(status.StatusInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AgentStatus indicates an expected call of AgentStatus
func (mr *MockUnitMockRecorder) AgentStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AgentStatus", reflect.TypeOf((*MockUnit)(nil).AgentStatus))
}

// Name mocks base method
func (m *MockUnit) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockUnitMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockUnit)(nil).Name))
}

// Status mocks base method
func (m *MockUnit) Status() (status.StatusInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(status.StatusInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Status indicates an expected call of Status
func (mr *MockUnitMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockUnit)(nil).Status))
}

// UnitTag mocks base method
func (m *MockUnit) UnitTag() v4.UnitTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnitTag")
	ret0, _ := ret[0].(v4.UnitTag)
	return ret0
}

// UnitTag indicates an expected call of UnitTag
func (mr *MockUnitMockRecorder) UnitTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnitTag", reflect.TypeOf((*MockUnit)(nil).UnitTag))
}

// MockCharm is a mock of Charm interface
type MockCharm struct {
	ctrl     *gomock.Controller
	recorder *MockCharmMockRecorder
}

// MockCharmMockRecorder is the mock recorder for MockCharm
type MockCharmMockRecorder struct {
	mock *MockCharm
}

// NewMockCharm creates a new mock instance
func NewMockCharm(ctrl *gomock.Controller) *MockCharm {
	mock := &MockCharm{ctrl: ctrl}
	mock.recorder = &MockCharmMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCharm) EXPECT() *MockCharmMockRecorder {
	return m.recorder
}

// Manifest mocks base method
func (m *MockCharm) Manifest() *v9.Manifest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Manifest")
	ret0, _ := ret[0].(*v9.Manifest)
	return ret0
}

// Manifest indicates an expected call of Manifest
func (mr *MockCharmMockRecorder) Manifest() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Manifest", reflect.TypeOf((*MockCharm)(nil).Manifest))
}

// Meta mocks base method
func (m *MockCharm) Meta() *v9.Meta {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Meta")
	ret0, _ := ret[0].(*v9.Meta)
	return ret0
}

// Meta indicates an expected call of Meta
func (mr *MockCharmMockRecorder) Meta() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Meta", reflect.TypeOf((*MockCharm)(nil).Meta))
}

// String mocks base method
func (m *MockCharm) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String
func (mr *MockCharmMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockCharm)(nil).String))
}

// URL mocks base method
func (m *MockCharm) URL() *v9.URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "URL")
	ret0, _ := ret[0].(*v9.URL)
	return ret0
}

// URL indicates an expected call of URL
func (mr *MockCharmMockRecorder) URL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "URL", reflect.TypeOf((*MockCharm)(nil).URL))
}

// MockCharmhubClient is a mock of CharmhubClient interface
type MockCharmhubClient struct {
	ctrl     *gomock.Controller
	recorder *MockCharmhubClientMockRecorder
}

// MockCharmhubClientMockRecorder is the mock recorder for MockCharmhubClient
type MockCharmhubClientMockRecorder struct {
	mock *MockCharmhubClient
}

// NewMockCharmhubClient creates a new mock instance
func NewMockCharmhubClient(ctrl *gomock.Controller) *MockCharmhubClient {
	mock := &MockCharmhubClient{ctrl: ctrl}
	mock.recorder = &MockCharmhubClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCharmhubClient) EXPECT() *MockCharmhubClientMockRecorder {
	return m.recorder
}

// Refresh mocks base method
func (m *MockCharmhubClient) Refresh(arg0 context.Context, arg1 charmhub.RefreshConfig) ([]transport.RefreshResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", arg0, arg1)
	ret0, _ := ret[0].([]transport.RefreshResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh
func (mr *MockCharmhubClientMockRecorder) Refresh(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockCharmhubClient)(nil).Refresh), arg0, arg1)
}
