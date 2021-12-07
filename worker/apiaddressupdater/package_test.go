// Copyright 2014-2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiaddressupdater_test

import (
	stdtesting "testing"

	coretesting "github.com/DavinZhang/juju/testing"
)

func TestPackage(t *stdtesting.T) {
	coretesting.MgoTestPackage(t)
}
