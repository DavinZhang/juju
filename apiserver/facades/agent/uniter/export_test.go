// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package uniter

import (
	"github.com/DavinZhang/juju/apiserver/common"
	"github.com/DavinZhang/juju/apiserver/facade"
	"github.com/DavinZhang/juju/apiserver/facades/agent/meterstatus"
	"github.com/DavinZhang/juju/caas"
	"github.com/DavinZhang/juju/state"
)

var (
	GetZone                = &getZone
	WatchStorageAttachment = watchStorageAttachment

	_ meterstatus.MeterStatus = (*UniterAPI)(nil)
)

type (
	Backend                    backend
	StorageStateInterface      storageInterface
	StorageVolumeInterface     = storageVolumeInterface
	StorageFilesystemInterface = storageFilesystemInterface
)

func NewStorageAPI(
	backend backend,
	storage storageAccess,
	resources facade.Resources,
	accessUnit common.GetAuthFunc,
) (*StorageAPI, error) {
	return newStorageAPI(backend, storage, resources, accessUnit)
}

func SetNewContainerBrokerFunc(api *UniterAPI, newBroker caas.NewContainerBrokerFunc) {
	api.containerBrokerFunc = newBroker
}

type patcher interface {
	PatchValue(interface{}, interface{})
}

func PatchGetStorageStateError(patcher patcher, err error) {
	patcher.PatchValue(&getStorageState, func(st *state.State) (storageAccess, error) { return nil, err })
}

func (n *NetworkInfoIAAS) MachineNetworkInfos() (map[string][]NetInfoAddress, error) {
	err := n.populateMachineAddresses()
	return n.machineAddresses, err
}
