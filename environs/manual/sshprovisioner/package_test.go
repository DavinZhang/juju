// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.
package sshprovisioner_test

import (
	"runtime"
	stdtesting "testing"

	"github.com/DavinZhang/juju/testing"
)

func Test(t *stdtesting.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Manual provider is not supported on windows")
	}
	testing.MgoTestPackage(t)
}
