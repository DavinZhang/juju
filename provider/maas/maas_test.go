// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package maas

import (
	"path/filepath"

	"github.com/juju/os/v2/series"
	"github.com/juju/utils/v2"
	"github.com/juju/utils/v2/arch"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/environs/context"
	sstesting "github.com/DavinZhang/juju/environs/simplestreams/testing"
	envtesting "github.com/DavinZhang/juju/environs/testing"
	envtools "github.com/DavinZhang/juju/environs/tools"
	"github.com/DavinZhang/juju/juju/keys"
	coretesting "github.com/DavinZhang/juju/testing"
	"github.com/DavinZhang/juju/version"
)

type baseProviderSuite struct {
	coretesting.FakeJujuXDGDataHomeSuite
	envtesting.ToolsFixture
	controllerUUID string

	callCtx           *context.CloudCallContext
	invalidCredential bool
}

func (s *baseProviderSuite) setupFakeTools(c *gc.C) {
	s.PatchValue(&keys.JujuPublicKey, sstesting.SignedMetadataPublicKey)
	storageDir := c.MkDir()
	toolsDir := filepath.Join(storageDir, "tools")
	s.PatchValue(&envtools.DefaultBaseURL, utils.MakeFileURL(toolsDir))
	s.UploadFakeToolsToDirectory(c, storageDir, "released", "released")
}

func (s *baseProviderSuite) SetUpSuite(c *gc.C) {
	s.FakeJujuXDGDataHomeSuite.SetUpSuite(c)
	restoreTimeouts := envtesting.PatchAttemptStrategies(&shortAttempt)
	restoreFinishBootstrap := envtesting.DisableFinishBootstrap()
	s.AddCleanup(func(*gc.C) {
		restoreFinishBootstrap()
		restoreTimeouts()
	})
}

func (s *baseProviderSuite) SetUpTest(c *gc.C) {
	s.FakeJujuXDGDataHomeSuite.SetUpTest(c)
	s.ToolsFixture.SetUpTest(c)
	s.PatchValue(&version.Current, coretesting.FakeVersionNumber)
	s.PatchValue(&arch.HostArch, func() string { return arch.AMD64 })
	s.PatchValue(&series.HostSeries, func() (string, error) { return version.DefaultSupportedLTS(), nil })
	s.callCtx = &context.CloudCallContext{
		InvalidateCredentialFunc: func(string) error {
			s.invalidCredential = true
			return nil
		},
	}
}

func (s *baseProviderSuite) TearDownTest(c *gc.C) {
	s.invalidCredential = false
	s.ToolsFixture.TearDownTest(c)
	s.FakeJujuXDGDataHomeSuite.TearDownTest(c)
}

func (s *baseProviderSuite) TearDownSuite(c *gc.C) {
	s.FakeJujuXDGDataHomeSuite.TearDownSuite(c)
}
