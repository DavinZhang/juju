// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package watcher_test

import (
	stdtesting "testing"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	coretesting "github.com/DavinZhang/juju/testing"
)

func TestPackage(t *stdtesting.T) {
	gc.TestingT(t)
}

type ImportTest struct{}

var _ = gc.Suite(&ImportTest{})

func (*ImportTest) TestImports(c *gc.C) {
	found := coretesting.FindJujuCoreImports(c, "github.com/DavinZhang/juju/core/watcher")

	// This package brings in nothing else from outside DavinZhang/juju/core
	c.Assert(found, jc.SameContents, []string{
		"core/life",
		"core/migration",
		"core/network",
		"core/secrets",
		"core/status",
		//  TODO: these have been brought in from migration and this is BAD.
		"resource",
	})

}
