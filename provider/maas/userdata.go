// Copyright 2015 Canonical Ltd.
// Copyright 2015 Cloudbase Solutions SRL
// Licensed under the AGPLv3, see LICENCE file for details.

package maas

import (
	"github.com/juju/errors"
	"github.com/juju/utils/v2"

	"github.com/DavinZhang/juju/cloudconfig/cloudinit"
	"github.com/DavinZhang/juju/cloudconfig/providerinit/renderers"
	jujuos "github.com/DavinZhang/juju/core/os"
)

type MAASRenderer struct{}

func (MAASRenderer) Render(cfg cloudinit.CloudConfig, os jujuos.OSType) ([]byte, error) {
	switch os {
	case jujuos.Ubuntu, jujuos.CentOS:
		return renderers.RenderYAML(cfg, utils.Gzip, renderers.ToBase64)
	case jujuos.Windows:
		return renderers.RenderYAML(cfg, renderers.WinEmbedInScript, renderers.ToBase64)
	default:
		return nil, errors.Errorf("Cannot encode userdata for OS: %s", os.String())
	}
}
