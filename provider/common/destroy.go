// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common

import (
	"strings"

	"github.com/juju/errors"

	"github.com/DavinZhang/juju/core/instance"
	"github.com/DavinZhang/juju/environs"
	"github.com/DavinZhang/juju/environs/context"
	"github.com/DavinZhang/juju/storage"
)

// Destroy is a common implementation of the Destroy method defined on
// environs.Environ; we strongly recommend that this implementation be
// used when writing a new provider.
func Destroy(env environs.Environ, ctx context.ProviderCallContext) error {
	logger.Infof("destroying model %q", env.Config().Name())
	if err := destroyInstances(env, ctx); err != nil {
		return errors.Annotate(err, "destroying instances")
	}
	if err := destroyStorage(env, ctx); err != nil {
		return errors.Annotate(err, "destroying storage")
	}
	return nil
}

func destroyInstances(env environs.Environ, ctx context.ProviderCallContext) error {
	logger.Infof("destroying instances")
	instances, err := env.AllInstances(ctx)
	switch err {
	case nil:
		ids := make([]instance.Id, len(instances))
		for i, inst := range instances {
			ids[i] = inst.Id()
		}
		if err := env.StopInstances(ctx, ids...); err != nil {
			return err
		}
		fallthrough
	case environs.ErrNoInstances:
		return nil
	default:
		return err
	}
}

// TODO(axw) we should just make it Environ.Destroy's responsibility
// to destroy persistent storage. Trying to include it in the storage
// source abstraction doesn't work well with dynamic, non-persistent
// storage like tmpfs, rootfs, etc.
func destroyStorage(env environs.Environ, ctx context.ProviderCallContext) error {
	logger.Infof("destroying storage")
	storageProviderTypes, err := env.StorageProviderTypes()
	if err != nil {
		return errors.Trace(err)
	}
	for _, storageProviderType := range storageProviderTypes {
		storageProvider, err := env.StorageProvider(storageProviderType)
		if err != nil {
			return errors.Trace(err)
		}
		if !storageProvider.Dynamic() {
			continue
		}
		if storageProvider.Scope() != storage.ScopeEnviron {
			continue
		}
		storageConfig, err := storage.NewConfig(
			string(storageProviderType),
			storageProviderType,
			map[string]interface{}{},
		)
		if err != nil {
			return errors.Trace(err)
		}
		if storageProvider.Supports(storage.StorageKindBlock) {
			volumeSource, err := storageProvider.VolumeSource(storageConfig)
			if err != nil {
				return errors.Annotate(err, "getting volume source")
			}
			if err := destroyVolumes(volumeSource, ctx); err != nil {
				return errors.Trace(err)
			}
		}
	}
	return nil
}

func destroyVolumes(volumeSource storage.VolumeSource, ctx context.ProviderCallContext) error {
	volumeIds, err := volumeSource.ListVolumes(ctx)
	if err != nil {
		return errors.Annotate(err, "listing volumes")
	}

	var errStrings []string
	errs, err := volumeSource.DestroyVolumes(ctx, volumeIds)
	if err != nil {
		return errors.Annotate(err, "destroying volumes")
	}
	for _, err := range errs {
		if err != nil {
			errStrings = append(errStrings, err.Error())
		}
	}
	if len(errStrings) > 0 {
		return errors.Errorf("destroying volumes: %s", strings.Join(errStrings, ", "))
	}
	return nil
}
