// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package crosscontroller

import (
	"github.com/juju/errors"

	"github.com/DavinZhang/juju/apiserver/common"
	apiservererrors "github.com/DavinZhang/juju/apiserver/errors"
	"github.com/DavinZhang/juju/apiserver/facade"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/state"
	"github.com/DavinZhang/juju/state/watcher"
)

type localControllerInfoFunc func() ([]string, string, error)
type publicDNSAddressFunc func() (string, error)
type watchLocalControllerInfoFunc func() state.NotifyWatcher

// CrossControllerAPI provides access to the CrossModelRelations API facade.
type CrossControllerAPI struct {
	resources                facade.Resources
	localControllerInfo      localControllerInfoFunc
	publicDNSAddress         publicDNSAddressFunc
	watchLocalControllerInfo watchLocalControllerInfoFunc
}

// NewStateCrossControllerAPI creates a new server-side CrossModelRelations API facade
// backed by global state.
func NewStateCrossControllerAPI(ctx facade.Context) (*CrossControllerAPI, error) {
	st := ctx.State()
	return NewCrossControllerAPI(
		ctx.Resources(),
		func() ([]string, string, error) {
			return common.StateControllerInfo(st)
		},
		func() (string, error) {
			config, err := st.ControllerConfig()
			if err != nil {
				return "", errors.Trace(err)
			}
			return config.PublicDNSAddress(), nil
		},
		st.WatchAPIHostPortsForClients,
	)
}

// NewCrossControllerAPI returns a new server-side CrossControllerAPI facade.
func NewCrossControllerAPI(
	resources facade.Resources,
	localControllerInfo localControllerInfoFunc,
	publicDNSAddress publicDNSAddressFunc,
	watchLocalControllerInfo watchLocalControllerInfoFunc,
) (*CrossControllerAPI, error) {
	return &CrossControllerAPI{
		resources:                resources,
		localControllerInfo:      localControllerInfo,
		publicDNSAddress:         publicDNSAddress,
		watchLocalControllerInfo: watchLocalControllerInfo,
	}, nil
}

// WatchControllerInfo creates a watcher that notifies when the API info
// for the controller changes.
func (api *CrossControllerAPI) WatchControllerInfo() (params.NotifyWatchResults, error) {
	results := params.NotifyWatchResults{
		Results: make([]params.NotifyWatchResult, 1),
	}
	w := api.watchLocalControllerInfo()
	if _, ok := <-w.Changes(); !ok {
		results.Results[0].Error = apiservererrors.ServerError(watcher.EnsureErr(w))
		return results, nil
	}
	results.Results[0].NotifyWatcherId = api.resources.Register(w)
	return results, nil
}

// ControllerInfo returns the API info for the controller.
func (api *CrossControllerAPI) ControllerInfo() (params.ControllerAPIInfoResults, error) {
	results := params.ControllerAPIInfoResults{
		Results: make([]params.ControllerAPIInfoResult, 1),
	}
	addrs, caCert, err := api.localControllerInfo()
	if err != nil {
		results.Results[0].Error = apiservererrors.ServerError(err)
		return results, nil
	}
	publicDNSAddress, err := api.publicDNSAddress()
	if err != nil {
		results.Results[0].Error = apiservererrors.ServerError(err)
		return results, nil
	}
	results.Results[0].Addresses = append([]string{publicDNSAddress}, addrs...)
	results.Results[0].CACert = caCert
	return results, nil
}
