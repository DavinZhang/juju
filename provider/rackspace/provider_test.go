// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package rackspace_test

import (
	stdcontext "context"

	"github.com/juju/errors"
	"github.com/juju/jsonschema"
	"github.com/juju/testing"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/cloud"
	"github.com/DavinZhang/juju/environs"
	environscloudspec "github.com/DavinZhang/juju/environs/cloudspec"
	"github.com/DavinZhang/juju/environs/config"
	"github.com/DavinZhang/juju/environs/context"
	"github.com/DavinZhang/juju/provider/rackspace"
	coretesting "github.com/DavinZhang/juju/testing"
)

type providerSuite struct {
	provider      environs.EnvironProvider
	innerProvider *fakeProvider
}

var _ = gc.Suite(&providerSuite{})

func (s *providerSuite) SetUpTest(c *gc.C) {
	s.innerProvider = new(fakeProvider)
	s.provider = rackspace.NewProvider(s.innerProvider)
}

func (s *providerSuite) TestValidate(c *gc.C) {
	cfg, err := config.New(config.UseDefaults, map[string]interface{}{
		"name":            "some-name",
		"type":            "some-type",
		"uuid":            coretesting.ModelTag.Id(),
		"controller-uuid": coretesting.ControllerTag.Id(),
		"authorized-keys": "key",
	})
	c.Check(err, gc.IsNil)
	_, err = s.provider.Validate(cfg, nil)
	c.Check(err, gc.IsNil)
	s.innerProvider.CheckCallNames(c, "Validate")
}

func (s *providerSuite) TestPrepareConfig(c *gc.C) {
	args := environs.PrepareConfigParams{
		Cloud: environscloudspec.CloudSpec{
			Region: "dfw",
		},
	}
	s.provider.PrepareConfig(args)

	expect := args
	expect.Cloud.Region = "DFW"
	s.innerProvider.CheckCalls(c, []testing.StubCall{
		{"PrepareConfig", []interface{}{expect}},
	})
}

type fakeProvider struct {
	testing.Stub
}

func (p *fakeProvider) Version() int {
	p.MethodCall(p, "Version")
	return 0
}

func (p *fakeProvider) Open(_ stdcontext.Context, args environs.OpenParams) (environs.Environ, error) {
	p.MethodCall(p, "Open", args)
	return nil, nil
}

func (p *fakeProvider) PrepareForCreateEnvironment(controllerUUID string, cfg *config.Config) (*config.Config, error) {
	p.MethodCall(p, "PrepareForCreateEnvironment", controllerUUID, cfg)
	return nil, nil
}

func (p *fakeProvider) PrepareConfig(args environs.PrepareConfigParams) (*config.Config, error) {
	p.MethodCall(p, "PrepareConfig", args)
	return nil, nil
}

func (p *fakeProvider) PrepareForBootstrap(ctx environs.BootstrapContext, cfg *config.Config) (environs.Environ, error) {
	p.MethodCall(p, "PrepareForBootstrap", ctx, cfg)
	return nil, nil
}

func (p *fakeProvider) Validate(cfg, old *config.Config) (valid *config.Config, err error) {
	p.MethodCall(p, "Validate", cfg, old)
	return cfg, nil
}

func (p *fakeProvider) CloudSchema() *jsonschema.Schema {
	p.MethodCall(p, "CloudSchema")
	return nil
}

// Ping tests the connection to the cloud, to verify the endpoint is valid.
func (p *fakeProvider) Ping(callCtx context.ProviderCallContext, endpoint string) error {
	return errors.NotImplementedf("Ping")
}

func (p *fakeProvider) CredentialSchemas() map[cloud.AuthType]cloud.CredentialSchema {
	p.MethodCall(p, "CredentialSchemas")
	return nil
}

func (p *fakeProvider) DetectCredentials(cloudName string) (*cloud.CloudCredential, error) {
	p.MethodCall(p, "DetectCredentials")
	return nil, errors.NotFoundf("credentials")
}

func (p *fakeProvider) FinalizeCredential(ctx environs.FinalizeCredentialContext, args environs.FinalizeCredentialParams) (*cloud.Credential, error) {
	p.MethodCall(p, "FinalizeCredential", ctx, args)
	return &args.Credential, nil
}
