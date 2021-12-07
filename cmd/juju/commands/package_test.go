// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package commands_test

import (
	stdtesting "testing"

	"github.com/DavinZhang/juju/component/all"
	"github.com/DavinZhang/juju/testing"
)

//go:generate go run github.com/golang/mock/mockgen -package commands -destination mockenvirons_test.go github.com/DavinZhang/juju/environs Environ,PrecheckJujuUpgradeStep
//go:generate go run github.com/golang/mock/mockgen -package commands -destination mockupgradeenvirons_test.go github.com/DavinZhang/juju/cmd/juju/commands UpgradePrecheckEnviron

func init() {
	if err := all.RegisterForClient(); err != nil {
		panic(err)
	}
}

func TestPackage(t *stdtesting.T) {
	testing.MgoTestPackage(t)
}
