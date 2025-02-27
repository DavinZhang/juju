// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// This package exists solely to avoid circular imports.

package factory

import (
	"github.com/juju/errors"

	"github.com/DavinZhang/juju/container"
	"github.com/DavinZhang/juju/container/kvm"
	"github.com/DavinZhang/juju/container/lxd"
	"github.com/DavinZhang/juju/core/instance"
)

// NewContainerManager creates the appropriate container.Manager for the
// specified container type.
var NewContainerManager = func(forType instance.ContainerType, conf container.ManagerConfig) (container.Manager, error) {
	switch forType {
	case instance.LXD:
		svr, err := lxd.MaybeNewLocalServer()
		if err != nil {
			return nil, errors.Annotate(err, "creating LXD container manager")
		}
		return lxd.NewContainerManager(conf, svr)
	case instance.KVM:
		return kvm.NewContainerManager(conf)
	}
	return nil, errors.Errorf("unknown container type: %q", forType)
}
