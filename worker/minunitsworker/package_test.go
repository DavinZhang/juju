// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package minunitsworker_test

import (
	stdtesting "testing"

	coretesting "github.com/DavinZhang/juju/testing"
)

func TestPackage(t *stdtesting.T) {
	coretesting.MgoTestPackage(t)
}
