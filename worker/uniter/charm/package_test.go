// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charm_test

import (
	stdtesting "testing"

	coretesting "github.com/DavinZhang/juju/testing"
)

//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/mocks.go github.com/DavinZhang/juju/worker/uniter/charm BundleReader,BundleInfo,Bundle

func TestPackage(t *stdtesting.T) {
	// TODO(fwereade) 2014-03-21 not-worth-a-bug-number
	// rewrite BundlesDir tests to use the mocks below and not require an API
	// server and associated gubbins.
	coretesting.MgoTestPackage(t)
}
