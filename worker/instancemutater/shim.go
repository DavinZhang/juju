// Copyright 2019 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package instancemutater

import (
	"github.com/DavinZhang/juju/api/base"
	apiinstancemutater "github.com/DavinZhang/juju/api/instancemutater"
)

func NewClient(apiCaller base.APICaller) InstanceMutaterAPI {
	facade := apiinstancemutater.NewClient(apiCaller)
	return facade
}
