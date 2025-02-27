// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upgrader

import (
	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/names/v4"
	"github.com/juju/version/v2"

	"github.com/DavinZhang/juju/apiserver/common"
	apiservererrors "github.com/DavinZhang/juju/apiserver/errors"
	"github.com/DavinZhang/juju/apiserver/facade"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/environs/config"
	"github.com/DavinZhang/juju/state"
	"github.com/DavinZhang/juju/state/stateenvirons"
	"github.com/DavinZhang/juju/state/watcher"
	jujuversion "github.com/DavinZhang/juju/version"
)

var logger = loggo.GetLogger("juju.apiserver.upgrader")

// The upgrader facade is a bit unique vs the other API Facades, as it
// has two implementations that actually expose the same API and which
// one gets returned depends on who is calling.  Both of them conform
// to the exact Upgrader API, so the actual calls that are available
// do not depend on who is currently connected.

// NewUpgraderFacade provides the signature required for facade registration.
func NewUpgraderFacade(ctx facade.Context) (Upgrader, error) {
	auth := ctx.Auth()
	st := ctx.State()
	ctrlSt := ctx.StatePool().SystemState()
	resources := ctx.Resources()
	// The type of upgrader we return depends on who is asking.
	// Machines get an UpgraderAPI, units get a UnitUpgraderAPI.
	// This is tested in the api/upgrader package since there
	// are currently no direct srvRoot tests.
	// TODO(dfc) this is redundant
	tag, err := names.ParseTag(auth.GetAuthTag().String())
	if err != nil {
		return nil, apiservererrors.ErrPerm
	}
	model, err := st.Model()
	if err != nil {
		return nil, errors.Trace(err)
	}
	switch tag.(type) {
	case names.MachineTag, names.ControllerAgentTag, names.ApplicationTag, names.ModelTag:
		return NewUpgraderAPI(ctrlSt, st, resources, auth)
	case names.UnitTag:
		if model.Type() == state.ModelTypeCAAS {
			// For sidecar applications.
			return NewUpgraderAPI(ctrlSt, st, resources, auth)
		}
		return NewUnitUpgraderAPI(st, resources, auth)
	}
	// Not a machine or unit.
	return nil, apiservererrors.ErrPerm
}

type Upgrader interface {
	WatchAPIVersion(args params.Entities) (params.NotifyWatchResults, error)
	DesiredVersion(args params.Entities) (params.VersionResults, error)
	Tools(args params.Entities) (params.ToolsResults, error)
	SetTools(args params.EntitiesVersion) (params.ErrorResults, error)
}

// UpgraderAPI provides access to the Upgrader API facade.
type UpgraderAPI struct {
	*common.ToolsGetter
	*common.ToolsSetter

	st         *state.State
	m          *state.Model
	resources  facade.Resources
	authorizer facade.Authorizer
}

// NewUpgraderAPI creates a new server-side UpgraderAPI facade.
func NewUpgraderAPI(
	ctrlSt *state.State,
	st *state.State,
	resources facade.Resources,
	authorizer facade.Authorizer,
) (*UpgraderAPI, error) {
	if !authorizer.AuthMachineAgent() && !authorizer.AuthApplicationAgent() && !authorizer.AuthModelAgent() && !authorizer.AuthUnitAgent() {
		return nil, apiservererrors.ErrPerm
	}
	getCanReadWrite := func() (common.AuthFunc, error) {
		return authorizer.AuthOwner, nil
	}
	model, err := st.Model()
	if err != nil {
		return nil, err
	}
	urlGetter := common.NewToolsURLGetter(model.UUID(), ctrlSt)
	configGetter := stateenvirons.EnvironConfigGetter{Model: model}
	newEnviron := common.EnvironFuncForModel(model, configGetter)
	return &UpgraderAPI{
		ToolsGetter: common.NewToolsGetter(st, configGetter, st, urlGetter, getCanReadWrite, newEnviron),
		ToolsSetter: common.NewToolsSetter(st, getCanReadWrite),
		st:          st,
		m:           model,
		resources:   resources,
		authorizer:  authorizer,
	}, nil
}

// WatchAPIVersion starts a watcher to track if there is a new version
// of the API that we want to upgrade to
func (u *UpgraderAPI) WatchAPIVersion(args params.Entities) (params.NotifyWatchResults, error) {
	result := params.NotifyWatchResults{
		Results: make([]params.NotifyWatchResult, len(args.Entities)),
	}
	for i, agent := range args.Entities {
		tag, err := names.ParseTag(agent.Tag)
		if err != nil {
			return params.NotifyWatchResults{}, errors.Trace(err)
		}
		err = apiservererrors.ErrPerm
		if u.authorizer.AuthOwner(tag) {
			watch := u.m.WatchForModelConfigChanges()
			// Consume the initial event. Technically, API
			// calls to Watch 'transmit' the initial event
			// in the Watch response. But NotifyWatchers
			// have no state to transmit.
			if _, ok := <-watch.Changes(); ok {
				result.Results[i].NotifyWatcherId = u.resources.Register(watch)
				err = nil
			} else {
				err = watcher.EnsureErr(watch)
			}
		}
		result.Results[i].Error = apiservererrors.ServerError(err)
	}
	return result, nil
}

func (u *UpgraderAPI) getGlobalAgentVersion() (version.Number, *config.Config, error) {
	// Get the Agent Version requested in the Model Config
	cfg, err := u.m.ModelConfig()
	if err != nil {
		return version.Number{}, nil, err
	}
	agentVersion, ok := cfg.AgentVersion()
	if !ok {
		return version.Number{}, nil, errors.New("agent version not set in model config")
	}
	return agentVersion, cfg, nil
}

type hasIsManager interface {
	IsManager() bool
}

func (u *UpgraderAPI) entityIsManager(tag names.Tag) bool {
	entity, err := u.st.FindEntity(tag)
	if err != nil {
		return false
	}
	if m, ok := entity.(hasIsManager); !ok {
		return false
	} else {
		return m.IsManager()
	}
}

// DesiredVersion reports the Agent Version that we want that agent to be running
func (u *UpgraderAPI) DesiredVersion(args params.Entities) (params.VersionResults, error) {
	results := make([]params.VersionResult, len(args.Entities))
	if len(args.Entities) == 0 {
		return params.VersionResults{}, nil
	}
	agentVersion, _, err := u.getGlobalAgentVersion()
	if err != nil {
		return params.VersionResults{}, apiservererrors.ServerError(err)
	}
	// Is the desired version greater than the current API server version?
	isNewerVersion := agentVersion.Compare(jujuversion.Current) > 0
	for i, entity := range args.Entities {
		tag, err := names.ParseTag(entity.Tag)
		if err != nil {
			results[i].Error = apiservererrors.ServerError(err)
			continue
		}
		err = apiservererrors.ErrPerm
		if u.authorizer.AuthOwner(tag) {
			// Only return the globally desired agent version if the
			// asking entity is a machine agent with JobManageModel or
			// if this API server is running the globally desired agent
			// version. Otherwise report this API server's current
			// agent version.
			//
			// This ensures that state machine agents will upgrade
			// first - once they have restarted and are running the
			// new version other agents will start to see the new
			// agent version.
			if !isNewerVersion || u.entityIsManager(tag) {
				results[i].Version = &agentVersion
			} else {
				logger.Debugf("desired version is %s, but current version is %s and agent is not a manager node", agentVersion, jujuversion.Current)
				results[i].Version = &jujuversion.Current
			}
			err = nil
		}
		results[i].Error = apiservererrors.ServerError(err)
	}
	return params.VersionResults{Results: results}, nil
}
