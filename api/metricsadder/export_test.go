// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package metricsadder

import (
	"github.com/DavinZhang/juju/api/base/testing"
)

// PatchFacadeCall patches the State's facade such that
// FacadeCall method calls are diverted to the provided
// function.
func PatchFacadeCall(p testing.Patcher, client *Client, f func(request string, params, response interface{}) error) {
	testing.PatchFacadeCall(p, &client.facade, f)
}
