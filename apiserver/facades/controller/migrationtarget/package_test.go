// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package migrationtarget_test

import (
	stdtesting "testing"

	"github.com/DavinZhang/juju/component/all"
	"github.com/DavinZhang/juju/testing"
)

func TestAll(t *stdtesting.T) {
	testing.MgoTestPackage(t)
}

func init() {
	// Required for resources.
	if err := all.RegisterForServer(); err != nil {
		panic(err)
	}
}
