// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package logger_test

import (
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/api"
	"github.com/DavinZhang/juju/api/logger"
	"github.com/DavinZhang/juju/core/watcher/watchertest"
	jujutesting "github.com/DavinZhang/juju/juju/testing"
	"github.com/DavinZhang/juju/state"
)

type loggerSuite struct {
	jujutesting.JujuConnSuite

	// These are raw State objects. Use them for setup and assertions, but
	// should never be touched by the API calls themselves
	rawMachine     *state.Machine
	rawCharm       *state.Charm
	rawApplication *state.Application
	rawUnit        *state.Unit

	logger *logger.State
}

var _ = gc.Suite(&loggerSuite{})

func (s *loggerSuite) SetUpTest(c *gc.C) {
	s.JujuConnSuite.SetUpTest(c)
	var stateAPI api.Connection
	stateAPI, s.rawMachine = s.OpenAPIAsNewMachine(c)
	// Create the logger facade.
	s.logger = logger.NewState(stateAPI)
	c.Assert(s.logger, gc.NotNil)
}

func (s *loggerSuite) TestLoggingConfigWrongMachine(c *gc.C) {
	config, err := s.logger.LoggingConfig(names.NewMachineTag("42"))
	c.Assert(err, gc.ErrorMatches, "permission denied")
	c.Assert(config, gc.Equals, "")
}

func (s *loggerSuite) TestLoggingConfig(c *gc.C) {
	config, err := s.logger.LoggingConfig(s.rawMachine.Tag())
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(config, gc.Not(gc.Equals), "")
}

func (s *loggerSuite) setLoggingConfig(c *gc.C, loggingConfig string) {
	err := s.Model.UpdateModelConfig(map[string]interface{}{"logging-config": loggingConfig}, nil)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *loggerSuite) TestWatchLoggingConfig(c *gc.C) {
	watcher, err := s.logger.WatchLoggingConfig(s.rawMachine.Tag())
	c.Assert(err, jc.ErrorIsNil)
	wc := watchertest.NewNotifyWatcherC(c, watcher, s.BackingState.StartSync)
	defer wc.AssertStops()

	// Initial event
	wc.AssertOneChange()

	loggingConfig := "<root>=WARN;juju.log.test=DEBUG"
	s.setLoggingConfig(c, loggingConfig)
	// One change noticing the new version
	wc.AssertOneChange()
	// Setting the version to the same value doesn't trigger a change
	s.setLoggingConfig(c, loggingConfig)
	wc.AssertNoChange()

	loggingConfig = loggingConfig + ";wibble=DEBUG"
	s.setLoggingConfig(c, loggingConfig)
	wc.AssertOneChange()
}
