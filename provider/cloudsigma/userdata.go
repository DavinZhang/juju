// Copyright 2015 Canonical Ltd.
// Copyright 2015 Cloudbase Solutions SRL
// Licensed under the AGPLv3, see LICENCE file for details.

package cloudsigma

import (
	"github.com/juju/errors"

	"github.com/DavinZhang/juju/cloudconfig/cloudinit"
	"github.com/DavinZhang/juju/cloudconfig/providerinit/renderers"
	jujuos "github.com/DavinZhang/juju/core/os"
)

type CloudSigmaRenderer struct{}

func (CloudSigmaRenderer) Render(cfg cloudinit.CloudConfig, os jujuos.OSType) ([]byte, error) {
	switch os {
	case jujuos.Ubuntu, jujuos.CentOS:
		return renderers.RenderYAML(cfg, renderers.ToBase64)
	default:
		return nil, errors.Errorf("Cannot encode userdata for OS: %s", os.String())
	}
}
