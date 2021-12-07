// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package maas

import (
	"github.com/DavinZhang/juju/cloudconfig/cloudinit"
	"github.com/DavinZhang/juju/environs"
)

var (
	ShortAttempt = &shortAttempt
)

func NewCloudinitConfig(env environs.Environ, hostname, series string) (cloudinit.CloudConfig, error) {
	return env.(*maasEnviron).newCloudinitConfig(hostname, series)
}
