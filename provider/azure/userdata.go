// Copyright 2015 Canonical Ltd.
// Copyright 2015 Cloudbase Solutions SRL
// Licensed under the AGPLv3, see LICENCE file for details.

package azure

import (
	"github.com/juju/errors"
	"github.com/juju/utils/v2"

	"github.com/DavinZhang/juju/cloudconfig/cloudinit"
	"github.com/DavinZhang/juju/cloudconfig/providerinit/renderers"
	jujuos "github.com/DavinZhang/juju/core/os"
)

type AzureRenderer struct{}

func (AzureRenderer) Render(cfg cloudinit.CloudConfig, os jujuos.OSType) ([]byte, error) {
	switch os {
	case jujuos.Ubuntu:
		return renderers.RenderYAML(cfg, utils.Gzip, renderers.ToBase64)
	case jujuos.CentOS:
		return renderers.RenderScript(cfg, renderers.ToBase64)
	case jujuos.Windows:
		return renderers.RenderYAML(cfg, renderers.WinEmbedInScript, renderers.ToBase64)
	default:
		return nil, errors.Errorf("Cannot encode userdata for OS: %s", os)
	}
}
