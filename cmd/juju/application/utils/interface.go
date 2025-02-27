// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package utils

import (
	"github.com/DavinZhang/juju/api/common/charms"
	"github.com/DavinZhang/juju/resource"
)

//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/charmresource_mock.go github.com/DavinZhang/juju/cmd/juju/application/utils CharmClient
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/resourcefacade_mock.go github.com/DavinZhang/juju/cmd/juju/application/utils ResourceLister

// CharmClient defines a subset of the charms facade, as required
// by the upgrade-charm command and to GetMetaResources.
type CharmClient interface {
	CharmInfo(string) (*charms.CharmInfo, error)
}

// ResourceLister defines a subset of the resources facade, as required
// by the upgrade-charm command and to deploy bundles.
type ResourceLister interface {
	ListResources([]string) ([]resource.ApplicationResources, error)
}
