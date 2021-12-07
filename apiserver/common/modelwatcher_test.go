// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common_test

import (
	"context"
	"fmt"

	"github.com/juju/cmd/v3/cmdtesting"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/apiserver/common"
	"github.com/DavinZhang/juju/apiserver/params"
	apiservertesting "github.com/DavinZhang/juju/apiserver/testing"
	"github.com/DavinZhang/juju/cmd/modelcmd"
	"github.com/DavinZhang/juju/environs/bootstrap"
	"github.com/DavinZhang/juju/environs/config"
	"github.com/DavinZhang/juju/jujuclient"
	"github.com/DavinZhang/juju/provider/dummy"
	"github.com/DavinZhang/juju/state"
	"github.com/DavinZhang/juju/testing"
)

type modelWatcherSuite struct {
	testing.BaseSuite

	testingEnvConfig *config.Config
}

var _ = gc.Suite(&modelWatcherSuite{})

type fakeModelAccessor struct {
	modelConfig      *config.Config
	modelConfigError error
}

func (*fakeModelAccessor) WatchForModelConfigChanges() state.NotifyWatcher {
	return apiservertesting.NewFakeNotifyWatcher()
}

func (f *fakeModelAccessor) ModelConfig() (*config.Config, error) {
	if f.modelConfigError != nil {
		return nil, f.modelConfigError
	}
	return f.modelConfig, nil
}

func (s *modelWatcherSuite) TearDownTest(c *gc.C) {
	dummy.Reset(c)
	s.BaseSuite.TearDownTest(c)
}

func (s *modelWatcherSuite) TestWatchSuccess(c *gc.C) {
	resources := common.NewResources()
	s.AddCleanup(func(_ *gc.C) { resources.StopAll() })
	e := common.NewModelWatcher(
		&fakeModelAccessor{},
		resources,
		nil,
	)
	result, err := e.WatchForModelConfigChanges()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, gc.DeepEquals, params.NotifyWatchResult{"1", nil})
	c.Assert(resources.Count(), gc.Equals, 1)
}

func (*modelWatcherSuite) TestModelConfigSuccess(c *gc.C) {
	authorizer := apiservertesting.FakeAuthorizer{
		Tag:        names.NewMachineTag("0"),
		Controller: true,
	}
	testingModelConfig := testingEnvConfig(c)
	e := common.NewModelWatcher(
		&fakeModelAccessor{modelConfig: testingModelConfig},
		nil,
		authorizer,
	)
	result, err := e.ModelConfig()
	c.Assert(err, jc.ErrorIsNil)
	// Make sure we can read the secret attribute (i.e. it's not masked).
	c.Check(result.Config["secret"], gc.Equals, "pork")
	c.Check(map[string]interface{}(result.Config), jc.DeepEquals, testingModelConfig.AllAttrs())
}

func (*modelWatcherSuite) TestModelConfigFetchError(c *gc.C) {
	authorizer := apiservertesting.FakeAuthorizer{
		Tag:        names.NewMachineTag("0"),
		Controller: true,
	}
	e := common.NewModelWatcher(
		&fakeModelAccessor{
			modelConfigError: fmt.Errorf("pow"),
		},
		nil,
		authorizer,
	)
	_, err := e.ModelConfig()
	c.Assert(err, gc.ErrorMatches, "pow")
}

func testingEnvConfig(c *gc.C) *config.Config {
	env, err := bootstrap.PrepareController(
		false,
		modelcmd.BootstrapContext(context.Background(), cmdtesting.Context(c)),
		jujuclient.NewMemStore(),
		bootstrap.PrepareParams{
			ControllerConfig: testing.FakeControllerConfig(),
			ControllerName:   "dummycontroller",
			ModelConfig:      dummy.SampleConfig(),
			Cloud:            dummy.SampleCloudSpec(),
			AdminSecret:      "admin-secret",
		},
	)
	c.Assert(err, jc.ErrorIsNil)
	return env.Config()
}
