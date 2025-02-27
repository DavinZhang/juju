// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package manual_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/cloud"
	"github.com/DavinZhang/juju/environs"
	envtesting "github.com/DavinZhang/juju/environs/testing"
	"github.com/DavinZhang/juju/testing"
)

type credentialsSuite struct {
	testing.BaseSuite
	provider environs.EnvironProvider
}

var _ = gc.Suite(&credentialsSuite{})

func (s *credentialsSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)

	var err error
	s.provider, err = environs.Provider("manual")
	c.Assert(err, jc.ErrorIsNil)
}

func (s *credentialsSuite) TestCredentialSchemas(c *gc.C) {
	envtesting.AssertProviderAuthTypes(c, s.provider, "empty")
}

func (s *credentialsSuite) TestDetectCredentials(c *gc.C) {
	credentials, err := s.provider.DetectCredentials("")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(credentials, jc.DeepEquals, cloud.NewEmptyCloudCredential())
}
