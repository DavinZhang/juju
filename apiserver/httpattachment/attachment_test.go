// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package httpattachment_test

import (
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/testing"
)

type requestSuite struct {
	testing.BaseSuite
}

var _ = gc.Suite(&requestSuite{})

// TODO the functions in this package should be tested directly.
// https://bugs.launchpad.net/juju-core/+bug/1503990
