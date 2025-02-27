// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package space

import (
	"testing"

	gc "gopkg.in/check.v1"
)

//go:generate go run github.com/golang/mock/mockgen -package space -destination context_mock_test.go github.com/DavinZhang/juju/environs/context ProviderCallContext
//go:generate go run github.com/golang/mock/mockgen -package space -destination environs_mock_test.go github.com/DavinZhang/juju/environs BootstrapEnviron,NetworkingEnviron
//go:generate go run github.com/golang/mock/mockgen -package space -destination spaces_mock_test.go github.com/DavinZhang/juju/environs/space ReloadSpacesState,Space,Constraints

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}
