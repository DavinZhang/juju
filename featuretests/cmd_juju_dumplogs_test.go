// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package featuretests

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/juju/cmd/v3/cmdtesting"
	"github.com/juju/loggo"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/agent"
	"github.com/DavinZhang/juju/cmd/jujud/agent/agenttest"
	"github.com/DavinZhang/juju/cmd/jujud/dumplogs"
	"github.com/DavinZhang/juju/state"
	"github.com/DavinZhang/juju/testing/factory"
	"github.com/DavinZhang/juju/version"
)

type dumpLogsCommandSuite struct {
	agenttest.AgentSuite
}

func (s *dumpLogsCommandSuite) SetUpTest(c *gc.C) {
	s.AgentSuite.SetUpTest(c)
}

func (s *dumpLogsCommandSuite) TestRun(c *gc.C) {
	// Create a controller machine and an agent for it.
	m, password := s.Factory.MakeMachineReturningPassword(c, &factory.MachineParams{
		Jobs:  []state.MachineJob{state.JobManageModel},
		Nonce: agent.BootstrapNonce,
	})
	err := m.SetMongoPassword(password)
	c.Assert(err, jc.ErrorIsNil)

	s.PrimeStateAgent(c, m.Tag(), password)

	//  Create multiple environments and add some logs for each.
	st1 := s.Factory.MakeModel(c, nil)
	defer st1.Close()
	st2 := s.Factory.MakeModel(c, nil)
	defer st2.Close()
	states := []*state.State{s.State, st1, st2}

	t := time.Date(2015, 11, 4, 3, 2, 1, 0, time.UTC)
	for _, st := range states {
		w := state.NewDbLogger(st)
		defer w.Close()
		for i := 0; i < 3; i++ {
			err := w.Log([]state.LogRecord{{
				Time:     t,
				Entity:   "machine-42",
				Version:  version.Current,
				Module:   "module",
				Location: "location",
				Level:    loggo.INFO,
				Message:  fmt.Sprintf("%d", i),
				Labels:   []string{"http"},
			}})
			c.Assert(err, jc.ErrorIsNil)
		}
	}

	// Run the juju-dumplogs command
	command := dumplogs.NewCommand()
	context, err := cmdtesting.RunCommand(c, command, "--data-dir", s.DataDir())
	c.Assert(err, jc.ErrorIsNil)

	// Check the log file for each environment
	expectedLog := "machine-42: 2015-11-04 03:02:01 INFO module %d http"
	for _, st := range states {
		logName := context.AbsPath(fmt.Sprintf("%s.log", st.ModelUUID()))
		logFile, err := os.Open(logName)
		c.Assert(err, jc.ErrorIsNil)
		scanner := bufio.NewScanner(logFile)
		for i := 0; scanner.Scan(); i++ {
			c.Assert(scanner.Text(), gc.Equals, fmt.Sprintf(expectedLog, i))
		}
		c.Assert(scanner.Err(), jc.ErrorIsNil)
	}
}
