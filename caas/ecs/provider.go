// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package ecs

import (
	"net/url"

	jujuclock "github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/jsonschema"
	"github.com/juju/loggo"

	"github.com/DavinZhang/juju/caas"
	"github.com/DavinZhang/juju/caas/ecs/constants"
	"github.com/DavinZhang/juju/cloud"
	"github.com/DavinZhang/juju/environs"
	"github.com/DavinZhang/juju/environs/cloudspec"
	"github.com/DavinZhang/juju/environs/config"
	"github.com/DavinZhang/juju/environs/context"
)

var logger = loggo.GetLogger("juju.ecs.provider")

type environProvider struct {
	providerCredentials
}

var (
	_                environs.EnvironProvider = (*environProvider)(nil)
	providerInstance                          = environProvider{}
)

// Version is part of the EnvironProvider interface.
func (environProvider) Version() int {
	return 0
}

// Open is specified in the EnvironProvider interface.
func (p environProvider) Open(args environs.OpenParams) (caas.Broker, error) {
	if err := args.Cloud.Validate(); err != nil {
		return nil, errors.Trace(err)
	}
	awsCfg, err := cloudSpecToAWSConfig(args.Cloud)
	if err != nil {
		return nil, errors.Trace(err)
	}
	clusterName := args.Cloud.Credential.Attributes()[credAttrClusterName]
	return newEnviron(
		args.ControllerUUID,
		clusterName,
		jujuclock.WallClock,
		args.Config, awsCfg,
		newECSClient,
	)
}

// CloudSchema returns the schema used to validate input for add-cloud.  Since
// this provider does not support custom clouds, this always returns nil.
func (p environProvider) CloudSchema() *jsonschema.Schema {
	return nil
}

// Ping tests the connection to the cloud, to verify the endpoint is valid.
func (p environProvider) Ping(ctx context.ProviderCallContext, endpoint string) error {
	return errors.NotImplementedf("Ping")
}

// PrepareConfig is specified in the EnvironProvider interface.
func (p environProvider) PrepareConfig(args environs.PrepareConfigParams) (*config.Config, error) {
	if err := p.validateCloudSpec(args.Cloud); err != nil {
		return nil, errors.Annotate(err, "validating cloud spec")
	}
	// Set the default storage sources.
	attrs := make(map[string]interface{})
	if _, ok := args.Config.StorageDefaultBlockSource(); !ok {
		attrs[config.StorageDefaultBlockSourceKey] = constants.StorageProviderType
	}
	if _, ok := args.Config.StorageDefaultFilesystemSource(); !ok {
		attrs[config.StorageDefaultFilesystemSourceKey] = constants.StorageProviderType
	}
	return args.Config.Apply(attrs)
}

// DetectRegions is specified in the environs.CloudRegionDetector interface.
func (p environProvider) DetectRegions() ([]cloud.Region, error) {
	return nil, errors.NotFoundf("regions")
}

func (p environProvider) validateCloudSpec(spec cloudspec.CloudSpec) error {
	if err := spec.Validate(); err != nil {
		return errors.Trace(err)
	}
	if err := validateCloudCredential(spec.Credential); err != nil {
		return errors.Trace(err)
	}
	if _, err := url.Parse(spec.Endpoint); err != nil {
		return errors.NotValidf("endpoint %q", spec.Endpoint)
	}
	if spec.Credential == nil {
		return errors.NotValidf("missing credential")
	}
	if authType := spec.Credential.AuthType(); authType != cloud.AccessKeyAuthType {
		return errors.NotSupportedf("%q auth-type", authType)
	}
	return nil
}
