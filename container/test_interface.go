// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package container

//go:generate go run github.com/golang/mock/mockgen -package testing -destination testing/interface_mock.go github.com/DavinZhang/juju/container TestLXDManager

// TestLXDManager for use in worker/provisioner tests
type TestLXDManager interface {
	Manager
	LXDProfileManager
	LXDProfileNameRetriever
}
