// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package credentialvalidator

import (
	"github.com/DavinZhang/juju/api/base"
	"github.com/DavinZhang/juju/api/credentialvalidator"
)

// NewFacade creates a *credentialvalidator.Facade and returns it as a Facade.
func NewFacade(apiCaller base.APICaller) (Facade, error) {
	facade := credentialvalidator.NewFacade(apiCaller)
	return facade, nil
}
