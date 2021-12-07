// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charms_test

import (
	stdtesting "testing"

	"github.com/DavinZhang/juju/testing"
)

func TestAll(t *stdtesting.T) {
	testing.MgoTestPackage(t)
}

//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/state_mock.go github.com/DavinZhang/juju/apiserver/facades/client/charms/interfaces BackendState,BackendModel,Application,Machine,Unit,Downloader
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/repository.go github.com/DavinZhang/juju/core/charm Repository,RepositoryFactory,CharmArchive
