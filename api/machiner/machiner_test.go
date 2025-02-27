// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package machiner_test

import (
	stdtesting "testing"

	"github.com/juju/errors"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/api"
	"github.com/DavinZhang/juju/api/machiner"
	apitesting "github.com/DavinZhang/juju/api/testing"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/core/life"
	"github.com/DavinZhang/juju/core/network"
	"github.com/DavinZhang/juju/core/status"
	"github.com/DavinZhang/juju/core/watcher/watchertest"
	"github.com/DavinZhang/juju/juju/testing"
	"github.com/DavinZhang/juju/state"
	coretesting "github.com/DavinZhang/juju/testing"
)

func TestAll(t *stdtesting.T) {
	coretesting.MgoTestPackage(t)
}

type machinerSuite struct {
	testing.JujuConnSuite
	*apitesting.APIAddresserTests

	st      api.Connection
	machine *state.Machine

	machiner *machiner.State
}

var _ = gc.Suite(&machinerSuite{})

func (s *machinerSuite) SetUpTest(c *gc.C) {
	s.JujuConnSuite.SetUpTest(c)
	m, err := s.State.AddMachine("quantal", state.JobManageModel)
	c.Assert(err, jc.ErrorIsNil)
	err = m.SetProviderAddresses(network.NewSpaceAddress("10.0.0.1"))
	c.Assert(err, jc.ErrorIsNil)

	s.st, s.machine = s.OpenAPIAsNewMachine(c)
	// Create the machiner API facade.
	s.machiner = machiner.NewState(s.st)
	c.Assert(s.machiner, gc.NotNil)
	waitForModelWatchersIdle := func(c *gc.C) {
		s.JujuConnSuite.WaitForModelWatchersIdle(c, s.BackingState.ModelUUID())
	}
	s.APIAddresserTests = apitesting.NewAPIAddresserTests(s.machiner, s.StatePool.SystemState(), s.BackingState, waitForModelWatchersIdle)
}

func (s *machinerSuite) TestMachineAndMachineTag(c *gc.C) {
	machine, err := s.machiner.Machine(names.NewMachineTag("42"))
	c.Assert(err, gc.ErrorMatches, ".*permission denied")
	c.Assert(err, jc.Satisfies, params.IsCodeUnauthorized)
	c.Assert(machine, gc.IsNil)

	machine1 := names.NewMachineTag("1")
	machine, err = s.machiner.Machine(machine1)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(machine.Tag(), gc.Equals, machine1)
}

func (s *machinerSuite) TestSetStatus(c *gc.C) {
	machine, err := s.machiner.Machine(names.NewMachineTag("1"))
	c.Assert(err, jc.ErrorIsNil)

	statusInfo, err := s.machine.Status()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(statusInfo.Status, gc.Equals, status.Pending)
	c.Assert(statusInfo.Message, gc.Equals, "")
	c.Assert(statusInfo.Data, gc.HasLen, 0)

	err = machine.SetStatus(status.Started, "blah", nil)
	c.Assert(err, jc.ErrorIsNil)

	statusInfo, err = s.machine.Status()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(statusInfo.Status, gc.Equals, status.Started)
	c.Assert(statusInfo.Message, gc.Equals, "blah")
	c.Assert(statusInfo.Data, gc.HasLen, 0)
	c.Assert(statusInfo.Since, gc.NotNil)
}

func (s *machinerSuite) TestEnsureDead(c *gc.C) {
	c.Assert(s.machine.Life(), gc.Equals, state.Alive)

	machine, err := s.machiner.Machine(names.NewMachineTag("1"))
	c.Assert(err, jc.ErrorIsNil)

	err = machine.EnsureDead()
	c.Assert(err, jc.ErrorIsNil)

	err = s.machine.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(s.machine.Life(), gc.Equals, state.Dead)

	err = machine.EnsureDead()
	c.Assert(err, jc.ErrorIsNil)
	err = s.machine.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(s.machine.Life(), gc.Equals, state.Dead)

	err = s.machine.Remove()
	c.Assert(err, jc.ErrorIsNil)
	err = s.machine.Refresh()
	c.Assert(err, jc.Satisfies, errors.IsNotFound)

	err = machine.EnsureDead()
	c.Assert(err, gc.ErrorMatches, "machine 1 not found")
	c.Assert(err, jc.Satisfies, params.IsCodeNotFound)
}

func (s *machinerSuite) TestRefresh(c *gc.C) {
	machine, err := s.machiner.Machine(names.NewMachineTag("1"))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(machine.Life(), gc.Equals, life.Alive)

	err = machine.EnsureDead()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(machine.Life(), gc.Equals, life.Alive)

	err = machine.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(machine.Life(), gc.Equals, life.Dead)
}

func (s *machinerSuite) TestSetMachineAddresses(c *gc.C) {
	machine, err := s.machiner.Machine(names.NewMachineTag("1"))
	c.Assert(err, jc.ErrorIsNil)

	addr := s.machine.Addresses()
	c.Assert(addr, gc.HasLen, 0)

	setAddresses := []network.MachineAddress{
		network.NewMachineAddress("8.8.8.8"),
		network.NewMachineAddress("127.0.0.1"),
		network.NewMachineAddress("10.0.0.1"),
	}
	// Before setting, the addresses are sorted to put public on top,
	// cloud-local next, machine-local last.
	expectAddresses := network.NewSpaceAddresses(
		"8.8.8.8",
		"10.0.0.1",
		"127.0.0.1",
	)
	err = machine.SetMachineAddresses(setAddresses)
	c.Assert(err, jc.ErrorIsNil)

	err = s.machine.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(s.machine.MachineAddresses(), jc.DeepEquals, expectAddresses)
}

func (s *machinerSuite) TestSetEmptyMachineAddresses(c *gc.C) {
	machine, err := s.machiner.Machine(names.NewMachineTag("1"))
	c.Assert(err, jc.ErrorIsNil)

	setAddresses := []network.MachineAddress{
		network.NewMachineAddress("8.8.8.8"),
		network.NewMachineAddress("127.0.0.1"),
		network.NewMachineAddress("10.0.0.1"),
	}
	err = machine.SetMachineAddresses(setAddresses)
	c.Assert(err, jc.ErrorIsNil)
	err = s.machine.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(s.machine.MachineAddresses(), gc.HasLen, 3)

	err = machine.SetMachineAddresses(nil)
	c.Assert(err, jc.ErrorIsNil)
	err = s.machine.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(s.machine.MachineAddresses(), gc.HasLen, 0)
}

func (s *machinerSuite) TestWatch(c *gc.C) {
	machine, err := s.machiner.Machine(names.NewMachineTag("1"))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(machine.Life(), gc.Equals, life.Alive)
	s.WaitForModelWatchersIdle(c, s.State.ModelUUID())

	w, err := machine.Watch()
	c.Assert(err, jc.ErrorIsNil)
	wc := watchertest.NewNotifyWatcherC(c, w, s.BackingState.StartSync)
	defer wc.AssertStops()

	// Initial event.
	wc.AssertOneChange()

	// Change something other than the lifecycle and make sure it's
	// not detected.
	err = machine.SetStatus(status.Started, "not really", nil)
	c.Assert(err, jc.ErrorIsNil)
	wc.AssertNoChange()

	// Make the machine dead and check it's detected.
	err = machine.EnsureDead()
	c.Assert(err, jc.ErrorIsNil)
	wc.AssertOneChange()
}

func (s *machinerSuite) TestRecordAgentStartInformation(c *gc.C) {
	mTag := names.NewMachineTag("1")
	stMachine, err := s.State.Machine(mTag.Id())
	c.Assert(err, jc.ErrorIsNil)
	oldStartedAt := stMachine.AgentStartTime()

	machine, err := s.machiner.Machine(mTag)
	c.Assert(err, jc.ErrorIsNil)

	err = machine.RecordAgentStartInformation("thundering-herds")
	c.Assert(err, jc.ErrorIsNil)

	err = stMachine.Refresh()
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(stMachine.AgentStartTime(), gc.Not(gc.Equals), oldStartedAt, gc.Commentf("expected the agent start time to be updated"))
	c.Assert(stMachine.Hostname(), gc.Equals, "thundering-herds", gc.Commentf("expected for the recorded machine hostname to be updated"))
}
