// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package environs

import (
	"net"
	"strconv"
	"time"

	"github.com/juju/errors"
	"github.com/juju/names/v4"
	"github.com/juju/utils/v2"

	"github.com/DavinZhang/juju/api"
	"github.com/DavinZhang/juju/core/instance"
	"github.com/DavinZhang/juju/core/network"
	"github.com/DavinZhang/juju/environs/context"
	"github.com/DavinZhang/juju/environs/instances"
)

// AddressesRefreshAttempt is the attempt strategy used when
// refreshing instance addresses.
var AddressesRefreshAttempt = utils.AttemptStrategy{
	Total: 3 * time.Minute,
	Delay: 1 * time.Second,
}

// getAddresses queries and returns the Addresses for the given instances,
// ignoring nil instances or ones without addresses.
func getAddresses(ctx context.ProviderCallContext, instances []instances.Instance) []network.ProviderAddress {
	var allAddrs []network.ProviderAddress
	for _, inst := range instances {
		if inst == nil {
			continue
		}
		addrs, err := inst.Addresses(ctx)
		if err != nil {
			logger.Debugf(
				"failed to get addresses for %v: %v (ignoring)",
				inst.Id(), err,
			)
			continue
		}
		allAddrs = append(allAddrs, addrs...)
	}
	return allAddrs
}

// waitAnyInstanceAddresses waits for at least one of the instances
// to have addresses, and returns them.
func waitAnyInstanceAddresses(
	env Environ,
	ctx context.ProviderCallContext,
	instanceIds []instance.Id,
) ([]network.ProviderAddress, error) {
	var addrs []network.ProviderAddress
	for a := AddressesRefreshAttempt.Start(); len(addrs) == 0 && a.Next(); {
		instances, err := env.Instances(ctx, instanceIds)
		if err != nil && err != ErrPartialInstances {
			logger.Debugf("error getting state instances: %v", err)
			return nil, err
		}
		addrs = getAddresses(ctx, instances)
	}
	if len(addrs) == 0 {
		return nil, errors.NotFoundf("addresses for %v", instanceIds)
	}
	return addrs, nil
}

// APIInfo returns an api.Info for the environment. The result is populated
// with addresses and CA certificate, but no tag or password.
func APIInfo(
	ctx context.ProviderCallContext, controllerUUID, modelUUID, caCert string, apiPort int, env Environ,
) (*api.Info, error) {
	instanceIds, err := env.ControllerInstances(ctx, controllerUUID)
	if err != nil {
		return nil, err
	}
	logger.Debugf("ControllerInstances returned: %v", instanceIds)
	addrs, err := waitAnyInstanceAddresses(env, ctx, instanceIds)
	if err != nil {
		return nil, err
	}

	apiAddrs := make([]string, len(addrs))
	for i, addr := range addrs {
		apiAddrs[i] = net.JoinHostPort(addr.Host(), strconv.Itoa(apiPort))
	}

	apiInfo := &api.Info{Addrs: apiAddrs, CACert: caCert, ModelTag: names.NewModelTag(modelUUID)}
	return apiInfo, nil
}

// CheckProviderAPI returns an error if a simple API call
// to check a basic response from the specified environ fails.
func CheckProviderAPI(env InstanceBroker, ctx context.ProviderCallContext) error {
	// We will make a simple API call to the provider
	// to ensure the underlying substrate is ok.
	_, err := env.AllInstances(ctx)
	switch err {
	case nil, ErrPartialInstances, ErrNoInstances:
		return nil
	}
	return errors.Annotate(err, "cannot make API call to provider")
}
