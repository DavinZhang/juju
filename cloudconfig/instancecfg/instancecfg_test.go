// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package instancecfg_test

import (
	jc "github.com/juju/testing/checkers"
	"github.com/juju/version/v2"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/cloudconfig/instancecfg"
	"github.com/DavinZhang/juju/core/model"
	"github.com/DavinZhang/juju/environs/config"
	"github.com/DavinZhang/juju/testing"
	coretools "github.com/DavinZhang/juju/tools"
)

type instancecfgSuite struct {
	testing.BaseSuite
}

var _ = gc.Suite(&instancecfgSuite{})

func (*instancecfgSuite) TestInstanceTagsController(c *gc.C) {
	cfg := testing.CustomModelConfig(c, testing.Attrs{})
	controllerJobs := []model.MachineJob{model.JobManageModel}
	nonControllerJobs := []model.MachineJob{model.JobHostUnits}
	testInstanceTags(c, cfg, controllerJobs, map[string]string{
		"juju-model-uuid":      testing.ModelTag.Id(),
		"juju-controller-uuid": testing.ControllerTag.Id(),
		"juju-is-controller":   "true",
	})
	testInstanceTags(c, cfg, nonControllerJobs, map[string]string{
		"juju-model-uuid":      testing.ModelTag.Id(),
		"juju-controller-uuid": testing.ControllerTag.Id(),
	})
}

func (*instancecfgSuite) TestInstanceTagsUserSpecified(c *gc.C) {
	cfg := testing.CustomModelConfig(c, testing.Attrs{
		"resource-tags": "a=b c=",
	})
	testInstanceTags(c, cfg, nil, map[string]string{
		"juju-model-uuid":      testing.ModelTag.Id(),
		"juju-controller-uuid": testing.ControllerTag.Id(),
		"a":                    "b",
		"c":                    "",
	})
}

func testInstanceTags(c *gc.C, cfg *config.Config, jobs []model.MachineJob, expectTags map[string]string) {
	tags := instancecfg.InstanceTags(testing.ModelTag.Id(), testing.ControllerTag.Id(), cfg, jobs)
	c.Assert(tags, jc.DeepEquals, expectTags)
}

func (*instancecfgSuite) TestAgentVersionZero(c *gc.C) {
	var icfg instancecfg.InstanceConfig
	c.Assert(icfg.AgentVersion(), gc.Equals, version.Binary{})
}

func (*instancecfgSuite) TestAgentVersion(c *gc.C) {
	var icfg instancecfg.InstanceConfig
	list := coretools.List{
		&coretools.Tools{Version: version.MustParseBinary("2.3.4-ubuntu-amd64")},
	}
	err := icfg.SetTools(list)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(icfg.AgentVersion(), gc.Equals, list[0].Version)
}

func (*instancecfgSuite) TestSetToolsSameVersions(c *gc.C) {
	var icfg instancecfg.InstanceConfig
	list := coretools.List{
		&coretools.Tools{Version: version.MustParseBinary("2.3.4-ubuntu-amd64")},
		&coretools.Tools{Version: version.MustParseBinary("2.3.4-ubuntu-amd64")},
	}
	err := icfg.SetTools(list)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(icfg.ToolsList(), jc.DeepEquals, list)
}

func (*instancecfgSuite) TestSetToolsDifferentVersions(c *gc.C) {
	var icfg instancecfg.InstanceConfig
	list := coretools.List{
		&coretools.Tools{Version: version.MustParseBinary("2.3.4-ubuntu-amd64")},
		&coretools.Tools{Version: version.MustParseBinary("2.3.5-ubuntu-amd64")},
	}
	err := icfg.SetTools(list)
	c.Assert(err, gc.ErrorMatches, `agent binary info mismatch.*2\.3\.4.*2\.3\.5.*`)
	c.Assert(icfg.ToolsList(), gc.HasLen, 0)
}

func (*instancecfgSuite) TestJujuTools(c *gc.C) {
	icfg := &instancecfg.InstanceConfig{
		DataDir: "/path/to/datadir/",
	}
	err := icfg.SetTools(coretools.List{
		&coretools.Tools{
			Version: version.MustParseBinary("2.3.4-ubuntu-amd64"),
			URL:     "/tools/2.3.4-ubuntu-amd64",
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(icfg.JujuTools(), gc.Equals, "/path/to/datadir/tools/2.3.4-ubuntu-amd64")
}

func (*instancecfgSuite) TestCharmDir(c *gc.C) {
	icfg := &instancecfg.InstanceConfig{
		DataDir: "/path/to/datadir/",
	}
	c.Assert(icfg.CharmDir(), gc.Equals, "/path/to/datadir/charms")
}
