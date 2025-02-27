// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package ec2_test

import (
	stdcontext "context"

	"github.com/aws/smithy-go"
	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/cloud"
	"github.com/DavinZhang/juju/environs"
	environscloudspec "github.com/DavinZhang/juju/environs/cloudspec"
	"github.com/DavinZhang/juju/environs/context"
	"github.com/DavinZhang/juju/provider/common"
	"github.com/DavinZhang/juju/provider/ec2"
	coretesting "github.com/DavinZhang/juju/testing"
)

type ProviderSuite struct {
	testing.IsolationSuite
	spec     environscloudspec.CloudSpec
	provider environs.EnvironProvider
}

var _ = gc.Suite(&ProviderSuite{})

func (s *ProviderSuite) SetUpTest(c *gc.C) {
	s.IsolationSuite.SetUpTest(c)

	credential := cloud.NewCredential(
		cloud.AccessKeyAuthType,
		map[string]string{
			"access-key": "foo",
			"secret-key": "bar",
		},
	)
	s.spec = environscloudspec.CloudSpec{
		Type:       "ec2",
		Name:       "aws",
		Region:     "us-east-1",
		Credential: &credential,
	}

	provider, err := environs.Provider("ec2")
	c.Assert(err, jc.ErrorIsNil)
	s.provider = provider
}

func (s *ProviderSuite) TestOpen(c *gc.C) {
	env, err := environs.Open(stdcontext.TODO(), s.provider, environs.OpenParams{
		Cloud:  s.spec,
		Config: coretesting.ModelConfig(c),
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(env, gc.NotNil)
}

func (s *ProviderSuite) TestOpenMissingCredential(c *gc.C) {
	s.spec.Credential = nil
	s.testOpenError(c, s.spec, `validating cloud spec: missing credential not valid`)
}

func (s *ProviderSuite) TestOpenUnsupportedCredential(c *gc.C) {
	credential := cloud.NewCredential(cloud.UserPassAuthType, map[string]string{})
	s.spec.Credential = &credential
	s.testOpenError(c, s.spec, `validating cloud spec: "userpass" auth-type not supported`)
}

func (s *ProviderSuite) testOpenError(c *gc.C, spec environscloudspec.CloudSpec, expect string) {
	_, err := environs.Open(stdcontext.TODO(), s.provider, environs.OpenParams{
		Cloud:  spec,
		Config: coretesting.ModelConfig(c),
	})
	c.Assert(err, gc.ErrorMatches, expect)
}

func (s *ProviderSuite) TestVerifyCredentialsErrs(c *gc.C) {
	err := ec2.VerifyCredentials(context.NewEmptyCloudCallContext())
	c.Assert(err, gc.Not(jc.ErrorIsNil))
	c.Assert(err, gc.Not(jc.Satisfies), common.IsCredentialNotValid)
}

func (s *ProviderSuite) TestMaybeConvertCredentialErrorIgnoresNil(c *gc.C) {
	err := ec2.MaybeConvertCredentialError(nil, context.NewEmptyCloudCallContext())
	c.Assert(err, jc.ErrorIsNil)
}

func (s *ProviderSuite) TestMaybeConvertCredentialErrorConvertsCredentialRelatedFailures(c *gc.C) {
	for _, code := range []string{
		"AuthFailure",
		"InvalidClientTokenId",
		"MissingAuthenticationToken",
		"Blocked",
		"CustomerKeyHasBeenRevoked",
		"PendingVerification",
		"SignatureDoesNotMatch",
	} {
		err := ec2.MaybeConvertCredentialError(
			&smithy.GenericAPIError{Code: code}, context.NewEmptyCloudCallContext())
		c.Assert(err, gc.NotNil)
		c.Assert(err, jc.Satisfies, common.IsCredentialNotValid)
	}
}

func (s *ProviderSuite) TestMaybeConvertCredentialErrorNotInvalidCredential(c *gc.C) {
	for _, code := range []string{
		"OptInRequired",
		"UnauthorizedOperation",
	} {
		err := ec2.MaybeConvertCredentialError(
			&smithy.GenericAPIError{Code: code}, context.NewEmptyCloudCallContext())
		c.Assert(err, gc.NotNil)
		c.Assert(err, gc.Not(jc.Satisfies), common.IsCredentialNotValid)
	}
}

func (s *ProviderSuite) TestMaybeConvertCredentialErrorHandlesOtherProviderErrors(c *gc.C) {
	// Any other ec2.Error is returned unwrapped.
	err := ec2.MaybeConvertCredentialError(&smithy.GenericAPIError{Code: "DryRunOperation"}, context.NewEmptyCloudCallContext())
	c.Assert(err, gc.Not(jc.ErrorIsNil))
	c.Assert(err, gc.Not(jc.Satisfies), common.IsCredentialNotValid)
}

func (s *ProviderSuite) TestConvertedCredentialError(c *gc.C) {
	// Trace() will keep error type
	inner := ec2.MaybeConvertCredentialError(
		&smithy.GenericAPIError{Code: "Blocked"}, context.NewEmptyCloudCallContext())
	traced := errors.Trace(inner)
	c.Assert(traced, gc.NotNil)
	c.Assert(traced, jc.Satisfies, common.IsCredentialNotValid)

	// Annotate() will keep error type
	annotated := errors.Annotate(inner, "annotation")
	c.Assert(annotated, gc.NotNil)
	c.Assert(annotated, jc.Satisfies, common.IsCredentialNotValid)

	// Running a CredentialNotValid through conversion call again is a no-op.
	again := ec2.MaybeConvertCredentialError(inner, context.NewEmptyCloudCallContext())
	c.Assert(again, gc.NotNil)
	c.Assert(again, jc.Satisfies, common.IsCredentialNotValid)
	c.Assert(again.Error(), jc.Contains, "\nYour Amazon account is currently blocked.: api error Blocked:")

	// Running an annotated CredentialNotValid through conversion call again is a no-op too.
	againAnotated := ec2.MaybeConvertCredentialError(annotated, context.NewEmptyCloudCallContext())
	c.Assert(againAnotated, gc.NotNil)
	c.Assert(againAnotated, jc.Satisfies, common.IsCredentialNotValid)
	c.Assert(againAnotated.Error(), jc.Contains, "\nYour Amazon account is currently blocked.: api error Blocked:")
}
