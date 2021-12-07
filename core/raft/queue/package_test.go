// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package queue

import (
	"testing"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	coretesting "github.com/DavinZhang/juju/testing"
)

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}

type ImportTest struct{}

var _ = gc.Suite(&ImportTest{})

func (*ImportTest) TestImports(c *gc.C) {
	found := coretesting.FindJujuCoreImports(c, "github.com/DavinZhang/juju/core/raft/queue")

	// This package brings in nothing else from outside DavinZhang/juju/core
	c.Assert(found, jc.SameContents, []string{})

}
