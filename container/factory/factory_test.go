// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package factory_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/container"
	"github.com/DavinZhang/juju/container/factory"
	"github.com/DavinZhang/juju/core/instance"
	"github.com/DavinZhang/juju/testing"
)

type factorySuite struct {
	testing.BaseSuite
}

var _ = gc.Suite(&factorySuite{})

func (*factorySuite) TestNewContainerManager(c *gc.C) {
	for _, test := range []struct {
		containerType instance.ContainerType
		valid         bool
	}{{
		containerType: instance.LXD,
		valid:         true,
	}, {
		containerType: instance.LXD,
		valid:         true,
	}, {
		containerType: instance.KVM,
		valid:         true,
	}, {
		containerType: instance.NONE,
		valid:         false,
	}, {
		containerType: instance.ContainerType("other"),
		valid:         false,
	}} {
		conf := container.ManagerConfig{container.ConfigModelUUID: testing.ModelTag.Id()}
		manager, err := factory.NewContainerManager(test.containerType, conf)
		if test.valid {
			c.Assert(err, jc.ErrorIsNil)
			c.Assert(manager, gc.NotNil)
		} else {
			c.Assert(err, gc.ErrorMatches, `unknown container type: ".*"`)
			c.Assert(manager, gc.IsNil)
		}
	}
}
