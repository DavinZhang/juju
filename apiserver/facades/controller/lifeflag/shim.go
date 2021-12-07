// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lifeflag

import (
	"github.com/DavinZhang/juju/apiserver/facade"
	"github.com/DavinZhang/juju/state"
)

// NewExternalFacade is for API registration.
func NewExternalFacade(st *state.State, resources facade.Resources, authorizer facade.Authorizer) (*Facade, error) {
	return NewFacade(st, resources, authorizer)
}
