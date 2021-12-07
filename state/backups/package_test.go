// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package backups_test

import (
	stdtesting "testing"

	"github.com/DavinZhang/juju/core/os"
	"github.com/DavinZhang/juju/testing"
)

func Test(t *stdtesting.T) {
	if os.HostOS() == os.Ubuntu {
		testing.MgoTestPackage(t)
	}
}
