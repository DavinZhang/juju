// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package cloudsigma

import (
	"github.com/juju/errors"
	"github.com/juju/loggo"

	"github.com/DavinZhang/juju/cloudconfig/instancecfg"
	"github.com/DavinZhang/juju/cloudconfig/providerinit"
	"github.com/DavinZhang/juju/core/instance"
	"github.com/DavinZhang/juju/environs"
	"github.com/DavinZhang/juju/environs/context"
	"github.com/DavinZhang/juju/environs/imagemetadata"
	"github.com/DavinZhang/juju/environs/instances"
	"github.com/DavinZhang/juju/tools"
)

//
// Imlementation of InstanceBroker: methods for starting and stopping instances.
//

var findInstanceImage = func(
	matchingImages []*imagemetadata.ImageMetadata,
) (*imagemetadata.ImageMetadata, error) {
	if len(matchingImages) == 0 {
		return nil, errors.New("no matching image meta data")
	}
	return matchingImages[0], nil
}

// StartInstance asks for a new instance to be created, associated with
// the provided config in machineConfig. The given config describes the juju
// state for the new instance to connect to. The config MachineNonce, which must be
// unique within an environment, is used by juju to protect against the
// consequences of multiple instances being started with the same machine id.
func (env *environ) StartInstance(ctx context.ProviderCallContext, args environs.StartInstanceParams) (*environs.StartInstanceResult, error) {
	logger.Infof("sigmaEnviron.StartInstance...")

	if args.InstanceConfig == nil {
		return nil, errors.New("instance configuration is nil")
	}

	if len(args.Tools) == 0 {
		return nil, errors.New("agent binaries not found")
	}

	img, err := findInstanceImage(args.ImageMetadata)
	if err != nil {
		return nil, err
	}

	tools, err := args.Tools.Match(tools.Filter{Arch: img.Arch})
	if err != nil {
		return nil, errors.Errorf("chosen architecture %v not present in %v", img.Arch, args.Tools.Arches())
	}

	if err := args.InstanceConfig.SetTools(tools); err != nil {
		return nil, errors.Trace(err)
	}
	if err := instancecfg.FinishInstanceConfig(args.InstanceConfig, env.Config()); err != nil {
		return nil, err
	}
	userData, err := providerinit.ComposeUserData(args.InstanceConfig, nil, CloudSigmaRenderer{})
	if err != nil {
		return nil, errors.Annotate(err, "cannot make user data")
	}

	logger.Debugf("cloudsigma user data; %d bytes", len(userData))

	client := env.client
	cfg := env.Config()
	server, rootdrive, arch, err := client.newInstance(args, img, userData, cfg.AuthorizedKeys())
	if err != nil {
		return nil, errors.Errorf("failed start instance: %v", err)
	}

	inst := &sigmaInstance{server: server}

	// prepare hardware characteristics
	hwch, err := inst.hardware(arch, rootdrive.Size())
	if err != nil {
		return nil, err
	}

	logger.Debugf("hardware: %v", hwch)
	return &environs.StartInstanceResult{
		Instance: inst,
		Hardware: hwch,
	}, nil
}

// AllInstances returns all instances currently known to the broker.
func (env *environ) AllInstances(ctx context.ProviderCallContext) ([]instances.Instance, error) {
	return env.instancesForMethod(ctx, "AllInstances")
}

// AllRunningInstances returns all running, available instances currently known to the broker.
func (env *environ) AllRunningInstances(ctx context.ProviderCallContext) ([]instances.Instance, error) {
	return env.instancesForMethod(ctx, "AllRunningInstances")
}

func (env *environ) instancesForMethod(ctx context.ProviderCallContext, method string) ([]instances.Instance, error) {
	// Please note that this must *not* return instances that have not been
	// allocated as part of this environment -- if it does, juju will see they
	// are not tracked in state, assume they're stale/rogue, and shut them down.

	logger.Tracef("environ.%v...", method)

	servers, err := env.client.instances()
	if err != nil {
		logger.Tracef("environ.%v failed: %v", method, err)
		return nil, err
	}

	instances := make([]instances.Instance, 0, len(servers))
	for _, server := range servers {
		instance := sigmaInstance{server: server}
		instances = append(instances, instance)
	}

	if logger.LogLevel() <= loggo.TRACE {
		logger.Tracef("%v, len = %d:", method, len(instances))
		for _, instance := range instances {
			logger.Tracef("... id: %q, status: %q", instance.Id(), instance.Status(ctx))
		}
	}

	return instances, nil
}

// Instances returns a slice of instances corresponding to the
// given instance ids.  If no instances were found, but there
// was no other error, it will return ErrNoInstances.  If
// some but not all the instances were found, the returned slice
// will have some nil slots, and an ErrPartialInstances error
// will be returned.
func (env *environ) Instances(ctx context.ProviderCallContext, ids []instance.Id) ([]instances.Instance, error) {
	logger.Tracef("environ.Instances %#v", ids)
	// Please note that this must *not* return instances that have not been
	// allocated as part of this environment -- if it does, juju will see they
	// are not tracked in state, assume they're stale/rogue, and shut them down.
	// This advice applies even if an instance id passed in corresponds to a
	// real instance that's not part of the environment -- the Environ should
	// treat that no differently to a request for one that does not exist.

	m, err := env.client.instanceMap()
	if err != nil {
		return nil, errors.Annotate(err, "environ.Instances failed")
	}

	var found int
	r := make([]instances.Instance, len(ids))
	for i, id := range ids {
		if s, ok := m[string(id)]; ok {
			r[i] = sigmaInstance{server: s}
			found++
		}
	}

	if found == 0 {
		err = environs.ErrNoInstances
	} else if found != len(ids) {
		err = environs.ErrPartialInstances
	}

	return r, errors.Trace(err)
}

// StopInstances shuts down the given instances.
func (env *environ) StopInstances(ctx context.ProviderCallContext, instances ...instance.Id) error {
	logger.Debugf("stop instances %+v", instances)

	var err error

	for _, instance := range instances {
		if e := env.client.stopInstance(instance); e != nil {
			err = e
		}
	}

	return err
}
