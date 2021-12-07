// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package azure

import (
	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/utils/v2/ssh"

	"github.com/DavinZhang/juju/environs"
	"github.com/DavinZhang/juju/provider/azure/internal/azureauth"
	"github.com/DavinZhang/juju/provider/azure/internal/azurecli"
)

const (
	// ProviderType defines the Azure provider.
	ProviderType = "azure"
)

// NewProvider instantiates and returns the Azure EnvironProvider using the
// given configuration.
func NewProvider(config ProviderConfig) (environs.CloudEnvironProvider, error) {
	environProvider, err := NewEnvironProvider(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return environProvider, nil
}

func init() {
	environProvider, err := NewProvider(ProviderConfig{
		RetryClock:                 &clock.WallClock,
		RandomWindowsAdminPassword: randomAdminPassword,
		GenerateSSHKey:             ssh.GenerateKey,
		ServicePrincipalCreator:    &azureauth.ServicePrincipalCreator{},
		AzureCLI:                   azurecli.AzureCLI{},
	})
	if err != nil {
		panic(err)
	}

	environs.RegisterProvider(ProviderType, environProvider)

	// TODO(axw) register an image metadata data source that queries
	// the Azure image registry, and introduce a way to disable the
	// common simplestreams source.
}
