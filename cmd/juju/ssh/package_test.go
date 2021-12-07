// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package ssh_test

import (
	stdtesting "testing"

	"github.com/DavinZhang/juju/testing"
)

//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/ssh_container_mock.go github.com/DavinZhang/juju/cmd/juju/ssh CloudCredentialAPI,ApplicationAPI,ModelAPI,CharmsAPI
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/context_mock.go github.com/DavinZhang/juju/cmd/juju/commands Context
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/k8s_exec_mock.go github.com/DavinZhang/juju/caas/kubernetes/provider/exec Executor
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination mocks/statusapi_mock.go github.com/DavinZhang/juju/cmd/juju/ssh StatusAPI

func TestPackage(t *stdtesting.T) {
	testing.MgoTestPackage(t)
}
