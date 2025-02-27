// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package firewaller

import (
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	"gopkg.in/macaroon.v2"

	"github.com/DavinZhang/juju/apiserver/common/firewall"
	"github.com/DavinZhang/juju/core/network"
	corefirewall "github.com/DavinZhang/juju/core/network/firewall"
	"github.com/DavinZhang/juju/state"
)

// State provides the subset of global state required by the
// remote firewaller facade.
type State interface {
	firewall.State

	ModelUUID() string
	GetMacaroon(entity names.Tag) (*macaroon.Macaroon, error)
	WatchOpenedPorts() state.StringsWatcher
	FindEntity(tag names.Tag) (state.Entity, error)
	FirewallRule(service corefirewall.WellKnownServiceType) (*state.FirewallRule, error)
	AllEndpointBindings() (map[string]map[string]string, error)
	SpaceInfos() (network.SpaceInfos, error)
}

// TODO(wallyworld) - for tests, remove when remaining firewaller tests become unit tests.
func StateShim(st *state.State, m *state.Model) stateShim {
	return stateShim{st: st, State: firewall.StateShim(st, m)}
}

type stateShim struct {
	firewall.State
	st *state.State
}

func (st stateShim) ModelUUID() string {
	return st.st.ModelUUID()
}

func (st stateShim) GetMacaroon(entity names.Tag) (*macaroon.Macaroon, error) {
	r := st.st.RemoteEntities()
	return r.GetMacaroon(entity)
}

func (st stateShim) FindEntity(tag names.Tag) (state.Entity, error) {
	return st.st.FindEntity(tag)
}

func (st stateShim) WatchOpenedPorts() state.StringsWatcher {
	return st.st.WatchOpenedPorts()
}

func (st stateShim) FirewallRule(service corefirewall.WellKnownServiceType) (*state.FirewallRule, error) {
	api := state.NewFirewallRules(st.st)
	return api.Rule(service)
}

func (st stateShim) AllEndpointBindings() (map[string]map[string]string, error) {
	model, err := st.st.Model()
	if err != nil {
		return nil, errors.Trace(err)
	}
	allEndpointBindings, err := model.AllEndpointBindings()
	if err != nil {
		return nil, errors.Trace(err)
	}

	res := make(map[string]map[string]string, len(allEndpointBindings))
	for appName, bindings := range allEndpointBindings {
		res[appName] = bindings.Map() // endpoint -> spaceID
	}
	return res, nil
}

func (st stateShim) SpaceInfos() (network.SpaceInfos, error) {
	return st.st.AllSpaceInfos()
}
