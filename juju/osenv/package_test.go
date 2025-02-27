// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package osenv_test

import (
	stdtesting "testing"

	gc "gopkg.in/check.v1"

	coretesting "github.com/DavinZhang/juju/testing"
)

func Test(t *stdtesting.T) {
	gc.TestingT(t)
}

type importSuite struct {
}

var _ = gc.Suite(&importSuite{})

func (*importSuite) TestDependencies(c *gc.C) {
	c.Assert(coretesting.FindJujuCoreImports(c, "github.com/DavinZhang/juju/juju/osenv"),
		gc.HasLen, 0)
}
