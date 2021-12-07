// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package migrationminion

import (
	"github.com/DavinZhang/juju/apiserver/facade"
	"github.com/DavinZhang/juju/state"
)

// NewFacade provides the signature required for facade registration.
func NewFacade(
	st *state.State,
	resources facade.Resources,
	authorizer facade.Authorizer,
) (*API, error) {
	return NewAPI(st, resources, authorizer)
}
