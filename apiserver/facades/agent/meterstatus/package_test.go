// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package meterstatus_test

import (
	stdtesting "testing"

	coretesting "github.com/DavinZhang/juju/testing"
)

func TestAll(t *stdtesting.T) {
	coretesting.MgoTestPackage(t)
}
