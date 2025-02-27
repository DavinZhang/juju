// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasagent_test

import (
	"github.com/juju/names/v4"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/apiserver/common"
	"github.com/DavinZhang/juju/apiserver/facade/facadetest"
	"github.com/DavinZhang/juju/apiserver/facades/agent/caasagent"
	apiservertesting "github.com/DavinZhang/juju/apiserver/testing"
	coretesting "github.com/DavinZhang/juju/testing"
)

var _ = gc.Suite(&caasagentSuite{})

type caasagentSuite struct {
	coretesting.BaseSuite

	resources  *common.Resources
	authorizer *apiservertesting.FakeAuthorizer
	facade     *caasagent.FacadeV2
}

func (s *caasagentSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)

	s.resources = common.NewResources()
	s.AddCleanup(func(_ *gc.C) { s.resources.StopAll() })

	s.authorizer = &apiservertesting.FakeAuthorizer{
		Tag: names.NewMachineTag("0"),
	}
}

func (s *caasagentSuite) TestPermission(c *gc.C) {
	s.authorizer = &apiservertesting.FakeAuthorizer{
		Tag: names.NewApplicationTag("someapp"),
	}
	_, err := caasagent.NewStateFacadeV2(facadetest.Context{Auth_: s.authorizer})
	c.Assert(err, gc.ErrorMatches, "permission denied")
}
