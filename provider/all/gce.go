// Copyright 2019 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

//go:build !minimal || provider_gce
// +build !minimal provider_gce

package all

import (
	// Register the provider.
	_ "github.com/DavinZhang/juju/provider/gce"
)
