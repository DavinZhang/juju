// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package presence_test

import (
	"testing"

	gc "gopkg.in/check.v1"

	coretesting "github.com/DavinZhang/juju/testing"
)

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}

type ImportTest struct{}

var _ = gc.Suite(&ImportTest{})

func (*ImportTest) TestImports(c *gc.C) {
	found := coretesting.FindJujuCoreImports(c, "github.com/DavinZhang/juju/core/presence")

	// This package brings in nothing else from DavinZhang/juju
	c.Assert(found, gc.HasLen, 0)
}
