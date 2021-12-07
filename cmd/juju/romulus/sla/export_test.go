// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package sla

import (
	"github.com/juju/romulus/api/sla"

	"github.com/DavinZhang/juju/api"
	"github.com/DavinZhang/juju/cmd/modelcmd"
)

var (
	NewAuthorizationClient = &newAuthorizationClient
	NewSLAClient           = &newSLAClient
	ModelId                = &modelId
	NewJujuClientStore     = &newJujuClientStore
)

// NewSLACommandForTest returns an slaCommand with apis provided by the given arguments
func NewSLACommandForTest(apiRoot api.Connection, slaC slaClient, authClient authorizationClient) modelcmd.ModelCommand {
	cmd := &slaCommand{
		newAPIRoot:   func() (api.Connection, error) { return apiRoot, nil },
		newSLAClient: func(api.Connection) slaClient { return slaC },
		newAuthorizationClient: func(options ...sla.ClientOption) (authorizationClient, error) {
			return authClient, nil
		},
	}
	return modelcmd.Wrap(cmd)
}
