// Copyright 2015 Canonical Ltd.
// Copyright 2015 Cloudbase Solutions SRL
// Licensed under the AGPLv3, see LICENCE file for details.

package ec2_test

import (
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v2"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/cloudconfig/cloudinit/cloudinittest"
	"github.com/DavinZhang/juju/cloudconfig/providerinit/renderers"
	"github.com/DavinZhang/juju/core/os"
	"github.com/DavinZhang/juju/provider/ec2"
	"github.com/DavinZhang/juju/testing"
)

type UserdataSuite struct {
	testing.BaseSuite
}

var _ = gc.Suite(&UserdataSuite{})

func (s *UserdataSuite) TestAmazonUnix(c *gc.C) {
	renderer := ec2.AmazonRenderer{}
	cloudcfg := &cloudinittest.CloudConfig{YAML: []byte("yaml")}

	result, err := renderer.Render(cloudcfg, os.Ubuntu)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, jc.DeepEquals, utils.Gzip(cloudcfg.YAML))

	result, err = renderer.Render(cloudcfg, os.CentOS)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, jc.DeepEquals, utils.Gzip(cloudcfg.YAML))
}

func (s *UserdataSuite) TestAmazonWindows(c *gc.C) {
	renderer := ec2.AmazonRenderer{}
	cloudcfg := &cloudinittest.CloudConfig{YAML: []byte("yaml")}

	result, err := renderer.Render(cloudcfg, os.Windows)
	c.Assert(err, jc.ErrorIsNil)
	expected := []byte(`<powershell>` +
		string(renderers.WinEmbedInScript(cloudcfg.YAML)) +
		`</powershell>`)
	c.Assert(result, jc.DeepEquals, expected)
}

func (s *UserdataSuite) TestAmazonUnknownOS(c *gc.C) {
	renderer := ec2.AmazonRenderer{}
	cloudcfg := &cloudinittest.CloudConfig{}
	result, err := renderer.Render(cloudcfg, os.GenericLinux)
	c.Assert(result, gc.IsNil)
	c.Assert(err, gc.ErrorMatches, "Cannot encode userdata for OS: GenericLinux")
}
