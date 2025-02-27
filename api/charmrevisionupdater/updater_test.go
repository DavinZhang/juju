// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charmrevisionupdater_test

import (
	gc "gopkg.in/check.v1"

	basetesting "github.com/DavinZhang/juju/api/base/testing"
	"github.com/DavinZhang/juju/api/charmrevisionupdater"
	"github.com/DavinZhang/juju/apiserver/params"
	coretesting "github.com/DavinZhang/juju/testing"
)

type versionUpdaterSuite struct {
	coretesting.BaseSuite
}

var _ = gc.Suite(&versionUpdaterSuite{})

func (s *versionUpdaterSuite) TestUpdateRevisions(c *gc.C) {
	apiCaller := basetesting.APICallerFunc(func(objType string, version int, id, request string, arg, result interface{}) error {
		c.Check(objType, gc.Equals, "CharmRevisionUpdater")
		c.Check(version, gc.Equals, 0)
		c.Check(id, gc.Equals, "")
		c.Check(request, gc.Equals, "UpdateLatestRevisions")
		c.Check(arg, gc.IsNil)
		c.Assert(result, gc.FitsTypeOf, &params.ErrorResult{})
		*(result.(*params.ErrorResult)) = params.ErrorResult{
			Error: &params.Error{Message: "boom"},
		}
		return nil
	})

	client := charmrevisionupdater.NewClient(apiCaller)
	err := client.UpdateLatestRevisions()
	c.Assert(err, gc.ErrorMatches, "boom")
}
