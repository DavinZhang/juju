// Copyright 2011, 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package ec2_test

import (
	"testing"

	gc "gopkg.in/check.v1"
)

//go:generate go run github.com/golang/mock/mockgen -package ec2 -destination context_mock_test.go github.com/DavinZhang/juju/environs/context ProviderCallContext

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}
