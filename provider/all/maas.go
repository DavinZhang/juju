// Copyright 2019 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

//go:build !minimal || provider_maas
// +build !minimal provider_maas

package all

import (
	// Register the provider.
	_ "github.com/DavinZhang/juju/provider/maas"
)
