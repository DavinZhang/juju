// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package equinix

import (
	jujuos "github.com/DavinZhang/juju/core/os"
	"github.com/juju/errors"

	"github.com/DavinZhang/juju/cloudconfig/cloudinit"
	"github.com/DavinZhang/juju/cloudconfig/providerinit/renderers"
)

type EquinixRenderer struct{}

func (EquinixRenderer) Render(cfg cloudinit.CloudConfig, os jujuos.OSType) ([]byte, error) {
	switch os {
	case jujuos.Ubuntu, jujuos.CentOS:
		return renderers.RenderYAML(cfg)
	default:
		return nil, errors.Errorf("Cannot encode userdata for OS: %s", os.String())
	}
}
