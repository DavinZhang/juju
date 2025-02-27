// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application_test

import (
	stdtesting "testing"

	"github.com/DavinZhang/juju/testing"
)

func TestAll(t *stdtesting.T) {
	testing.MgoTestPackage(t)
}

//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/storage_mock.go github.com/DavinZhang/juju/state/storage Storage
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/charm_mock.go github.com/DavinZhang/juju/apiserver/facades/client/application StateCharm
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/model_mock.go github.com/DavinZhang/juju/apiserver/facades/client/application StateModel
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/charmstore_mock.go github.com/DavinZhang/juju/apiserver/facades/client/application State
//go:generate go run github.com/golang/mock/mockgen -package application -destination updateseries_mocks_test.go github.com/DavinZhang/juju/apiserver/facades/client/application Application,Charm,CharmMeta,UpdateSeriesState,UpdateSeriesValidator,CharmhubClient
