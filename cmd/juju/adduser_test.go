// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	gc "launchpad.net/gocheck"

	jujutesting "launchpad.net/juju-core/juju/testing"
	"launchpad.net/juju-core/testing"
)

type AdduserSuite struct {
	jujutesting.RepoSuite
}

var _ = gc.Suite(&AdduserSuite{})

func (s *AdduserSuite) Testadduser(c *gc.C) {

	_, err := testing.RunCommand(c, &AdduserCommand{}, []string{"foobar", "password"})
	c.Assert(err, gc.IsNil)

	_, err = testing.RunCommand(c, &AdduserCommand{}, []string{"foobar", "newpassword"})
	c.Assert(err, gc.ErrorMatches, "user already exists")
}
