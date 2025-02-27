// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upgrader

import (
	"fmt"

	"github.com/juju/version/v2"

	"github.com/DavinZhang/juju/api/base"
	apiwatcher "github.com/DavinZhang/juju/api/watcher"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/core/watcher"
	"github.com/DavinZhang/juju/tools"
)

// State provides access to an upgrader worker's view of the state.
type State struct {
	facade base.FacadeCaller
}

// NewState returns a version of the state that provides functionality
// required by the upgrader worker.
func NewState(caller base.APICaller) *State {
	return &State{base.NewFacadeCaller(caller, "Upgrader")}
}

// SetVersion sets the tools version associated with the entity with
// the given tag, which must be the tag of the entity that the
// upgrader is running on behalf of.
func (st *State) SetVersion(tag string, v version.Binary) error {
	var results params.ErrorResults
	args := params.EntitiesVersion{
		AgentTools: []params.EntityVersion{{
			Tag:   tag,
			Tools: &params.Version{Version: v},
		}},
	}
	err := st.facade.FacadeCall("SetTools", args, &results)
	if err != nil {
		return err
	}
	return results.OneError()
}

func (st *State) DesiredVersion(tag string) (version.Number, error) {
	var results params.VersionResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: tag}},
	}
	err := st.facade.FacadeCall("DesiredVersion", args, &results)
	if err != nil {
		return version.Number{}, err
	}
	if len(results.Results) != 1 {
		return version.Number{}, fmt.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if err := result.Error; err != nil {
		return version.Number{}, err
	}
	if result.Version == nil {
		return version.Number{}, fmt.Errorf("received no error, but got a nil Version")
	}
	return *result.Version, nil
}

// Tools returns the agent tools that should run on the given entity,
// along with a flag whether to disable SSL hostname verification.
func (st *State) Tools(tag string) (tools.List, error) {
	var results params.ToolsResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: tag}},
	}
	err := st.facade.FacadeCall("Tools", args, &results)
	if err != nil {
		return nil, err
	}
	if len(results.Results) != 1 {
		return nil, fmt.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if err := result.Error; err != nil {
		return nil, err
	}
	return result.ToolsList, nil
}

func (st *State) WatchAPIVersion(agentTag string) (watcher.NotifyWatcher, error) {
	var results params.NotifyWatchResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: agentTag}},
	}
	err := st.facade.FacadeCall("WatchAPIVersion", args, &results)
	if err != nil {
		return nil, err
	}
	if len(results.Results) != 1 {
		return nil, fmt.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return nil, result.Error
	}
	w := apiwatcher.NewNotifyWatcher(st.facade.RawAPICaller(), result)
	return w, nil
}
