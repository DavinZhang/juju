// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lifeflag

import (
	"github.com/juju/names/v4"

	"github.com/DavinZhang/juju/apiserver/common"
	apiservererrors "github.com/DavinZhang/juju/apiserver/errors"
	"github.com/DavinZhang/juju/apiserver/facade"
	"github.com/DavinZhang/juju/state"
)

type Backend interface {
	ModelUUID() string
	state.EntityFinder
}

func NewFacade(backend Backend, resources facade.Resources, authorizer facade.Authorizer) (*Facade, error) {
	if !authorizer.AuthController() {
		return nil, apiservererrors.ErrPerm
	}
	expect := names.NewModelTag(backend.ModelUUID())
	getCanAccess := func() (common.AuthFunc, error) {
		return func(tag names.Tag) bool {
			return tag == expect
		}, nil
	}
	life := common.NewLifeGetter(backend, getCanAccess)
	watch := common.NewAgentEntityWatcher(backend, resources, getCanAccess)
	return &Facade{
		LifeGetter:         life,
		AgentEntityWatcher: watch,
	}, nil
}

type Facade struct {
	*common.LifeGetter
	*common.AgentEntityWatcher
}
