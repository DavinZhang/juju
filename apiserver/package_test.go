// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver_test

import (
	"testing"

	coretesting "github.com/DavinZhang/juju/testing"
)

func TestPackage(t *testing.T) {
	coretesting.MgoTestPackage(t)
}
