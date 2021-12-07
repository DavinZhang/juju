// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package deployer

import (
	apideployer "github.com/DavinZhang/juju/api/deployer"
)

func MakeAPIShim(st *apideployer.State) API {
	return &apiShim{st}
}
