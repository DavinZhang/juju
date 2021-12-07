// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package ecs

import (
	"github.com/DavinZhang/juju/caas"
	"github.com/DavinZhang/juju/caas/ecs/constants"
)

func init() {
	_ = caas.RegisterContainerProvider(constants.ECSProviderType, providerInstance)
}
