// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxd

import (
	"github.com/juju/errors"
	"github.com/lxc/lxd/shared/api"

	"github.com/DavinZhang/juju/container/lxd"
	"github.com/DavinZhang/juju/core/instance"
	"github.com/DavinZhang/juju/core/network"
	"github.com/DavinZhang/juju/core/status"
	"github.com/DavinZhang/juju/environs/context"
	"github.com/DavinZhang/juju/environs/instances"
)

type environInstance struct {
	container *lxd.Container
	env       *environ
}

var _ instances.Instance = (*environInstance)(nil)

func newInstance(container *lxd.Container, env *environ) *environInstance {
	return &environInstance{
		container: container,
		env:       env,
	}
}

// Id implements instances.Instance.
func (i *environInstance) Id() instance.Id {
	return instance.Id(i.container.Name)
}

// Status implements instances.Instance.
func (i *environInstance) Status(ctx context.ProviderCallContext) instance.Status {
	var jujuStatus status.Status
	code := i.container.StatusCode
	switch code {
	case api.Starting, api.Started:
		jujuStatus = status.Allocating
	case api.Running:
		jujuStatus = status.Running
	case api.Freezing, api.Frozen, api.Thawed, api.Stopping, api.Stopped:
		jujuStatus = status.Empty
	default:
		jujuStatus = status.Empty
	}
	return instance.Status{
		Status:  jujuStatus,
		Message: code.String(),
	}

}

// Addresses implements instances.Instance.
func (i *environInstance) Addresses(_ context.ProviderCallContext) (network.ProviderAddresses, error) {
	addrs, err := i.env.server().ContainerAddresses(i.container.Name)
	return addrs, errors.Trace(err)
}
