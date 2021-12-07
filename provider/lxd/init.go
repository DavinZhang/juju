// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxd

import (
	"github.com/DavinZhang/juju/environs"
	"github.com/DavinZhang/juju/provider/lxd/lxdnames"
)

func init() {
	environs.RegisterProvider(lxdnames.ProviderType, NewProvider())
}
