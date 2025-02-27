// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common_test

import (
	"io"

	"github.com/DavinZhang/juju/core/instance"
	"github.com/DavinZhang/juju/core/network"
	"github.com/DavinZhang/juju/environs"
	"github.com/DavinZhang/juju/environs/config"
	"github.com/DavinZhang/juju/environs/context"
	"github.com/DavinZhang/juju/environs/instances"
	"github.com/DavinZhang/juju/environs/simplestreams"
	"github.com/DavinZhang/juju/environs/storage"
	jujustorage "github.com/DavinZhang/juju/storage"
)

type allInstancesFunc func(context.ProviderCallContext) ([]instances.Instance, error)
type instancesFunc func(context.ProviderCallContext, []instance.Id) ([]instances.Instance, error)
type startInstanceFunc func(context.ProviderCallContext, environs.StartInstanceParams) (instances.Instance, *instance.HardwareCharacteristics, network.InterfaceInfos, error)
type stopInstancesFunc func(context.ProviderCallContext, []instance.Id) error
type getToolsSourcesFunc func() ([]simplestreams.DataSource, error)
type configFunc func() *config.Config
type setConfigFunc func(*config.Config) error

type mockEnviron struct {
	storage          storage.Storage
	allInstances     allInstancesFunc
	instances        instancesFunc
	startInstance    startInstanceFunc
	stopInstances    stopInstancesFunc
	getToolsSources  getToolsSourcesFunc
	config           configFunc
	setConfig        setConfigFunc
	storageProviders jujustorage.StaticProviderRegistry
	environs.Environ // stub out other methods with panics
}

func (env *mockEnviron) Storage() storage.Storage {
	return env.storage
}

func (env *mockEnviron) AllInstances(ctx context.ProviderCallContext) ([]instances.Instance, error) {
	return env.allInstances(ctx)
}

func (env *mockEnviron) AllRunningInstances(ctx context.ProviderCallContext) ([]instances.Instance, error) {
	return env.allInstances(ctx)
}

func (env *mockEnviron) Instances(ctx context.ProviderCallContext, ids []instance.Id) ([]instances.Instance, error) {
	return env.instances(ctx, ids)
}

func (env *mockEnviron) StartInstance(ctx context.ProviderCallContext, args environs.StartInstanceParams) (*environs.StartInstanceResult, error) {
	inst, hw, networkInfo, err := env.startInstance(ctx, args)
	if err != nil {
		return nil, err
	}
	return &environs.StartInstanceResult{
		Instance:    inst,
		Hardware:    hw,
		NetworkInfo: networkInfo,
	}, nil
}

func (env *mockEnviron) StopInstances(ctx context.ProviderCallContext, ids ...instance.Id) error {
	return env.stopInstances(ctx, ids)
}

func (env *mockEnviron) Config() *config.Config {
	return env.config()
}

func (env *mockEnviron) SetConfig(cfg *config.Config) error {
	if env.setConfig != nil {
		return env.setConfig(cfg)
	}
	return nil
}

func (env *mockEnviron) GetToolsSources() ([]simplestreams.DataSource, error) {
	if env.getToolsSources != nil {
		return env.getToolsSources()
	}
	datasource := storage.NewStorageSimpleStreamsDataSource("test cloud storage", env.Storage(), storage.BaseToolsPath, simplestreams.SPECIFIC_CLOUD_DATA, false)
	return []simplestreams.DataSource{datasource}, nil
}

func (env *mockEnviron) StorageProviderTypes() ([]jujustorage.ProviderType, error) {
	return env.storageProviders.StorageProviderTypes()
}

func (env *mockEnviron) StorageProvider(t jujustorage.ProviderType) (jujustorage.Provider, error) {
	return env.storageProviders.StorageProvider(t)
}

type availabilityZonesFunc func(context.ProviderCallContext) (network.AvailabilityZones, error)
type instanceAvailabilityZoneNamesFunc func(context.ProviderCallContext, []instance.Id) (map[instance.Id]string, error)
type deriveAvailabilityZonesFunc func(context.ProviderCallContext, environs.StartInstanceParams) ([]string, error)

type mockZonedEnviron struct {
	mockEnviron
	availabilityZones             availabilityZonesFunc
	instanceAvailabilityZoneNames instanceAvailabilityZoneNamesFunc
	deriveAvailabilityZones       deriveAvailabilityZonesFunc
}

func (env *mockZonedEnviron) AvailabilityZones(ctx context.ProviderCallContext) (network.AvailabilityZones, error) {
	return env.availabilityZones(ctx)
}

func (env *mockZonedEnviron) InstanceAvailabilityZoneNames(ctx context.ProviderCallContext, ids []instance.Id) (map[instance.Id]string, error) {
	return env.instanceAvailabilityZoneNames(ctx, ids)
}

func (env *mockZonedEnviron) DeriveAvailabilityZones(ctx context.ProviderCallContext, args environs.StartInstanceParams) ([]string, error) {
	return env.deriveAvailabilityZones(ctx, args)
}

type mockInstance struct {
	id                 string
	addresses          network.ProviderAddresses
	addressesErr       error
	dnsName            string
	dnsNameErr         error
	status             instance.Status
	instances.Instance // stub out other methods with panics
}

func (inst *mockInstance) Id() instance.Id {
	return instance.Id(inst.id)
}

func (inst *mockInstance) Status(context.ProviderCallContext) instance.Status {
	return inst.status
}

func (inst *mockInstance) Addresses(context.ProviderCallContext) (network.ProviderAddresses, error) {
	return inst.addresses, inst.addressesErr
}

type mockStorage struct {
	storage.Storage
	putErr       error
	removeAllErr error
}

func (stor *mockStorage) Put(name string, reader io.Reader, size int64) error {
	if stor.putErr != nil {
		return stor.putErr
	}
	return stor.Storage.Put(name, reader, size)
}

func (stor *mockStorage) RemoveAll() error {
	if stor.removeAllErr != nil {
		return stor.removeAllErr
	}
	return stor.Storage.RemoveAll()
}

type mockAvailabilityZone struct {
	name      string
	available bool
}

func (z *mockAvailabilityZone) Name() string {
	return z.name
}

func (z *mockAvailabilityZone) Available() bool {
	return z.available
}
