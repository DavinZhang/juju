// Copyright 2019 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package broker_test

import (
	stdtesting "testing"

	"github.com/DavinZhang/juju/testing"
)

func TestPackage(t *stdtesting.T) {
	testing.MgoTestPackage(t)
}
