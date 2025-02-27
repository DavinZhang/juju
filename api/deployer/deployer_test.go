// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package deployer_test

import (
	stdtesting "testing"

	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/api"
	"github.com/DavinZhang/juju/api/deployer"
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

type deployerSuite struct {
	testing.JujuConnSuite

	stateAPI api.Connection

	// These are raw State objects. Use them for setup and assertions, but
	// should never be touched by the API calls themselves
	machine     *state.Machine
	app0        *state.Application
	app1        *state.Application
	principal   *state.Unit
	subordinate *state.Unit

	st *deployer.State
}

var _ = gc.Suite(&deployerSuite{})

func (s *deployerSuite) SetUpTest(c *gc.C) {
	s.JujuConnSuite.SetUpTest(c)
	s.stateAPI, s.machine = s.OpenAPIAsNewMachine(c, state.JobManageModel, state.JobHostUnits)
	err := s.machine.SetProviderAddresses(network.NewSpaceAddress("0.1.2.3"))
	c.Assert(err, jc.ErrorIsNil)

	// Create the needed applications and relate them.
	s.app0 = s.AddTestingApplication(c, "mysql", s.AddTestingCharm(c, "mysql"))
	s.app1 = s.AddTestingApplication(c, "logging", s.AddTestingCharm(c, "logging"))
	eps, err := s.State.InferEndpoints("mysql", "logging")
	c.Assert(err, jc.ErrorIsNil)
	rel, err := s.State.AddRelation(eps...)
	c.Assert(err, jc.ErrorIsNil)

	// Create principal and subordinate units and assign them.
	s.principal, err = s.app0.AddUnit(state.AddUnitParams{})
	c.Assert(err, jc.ErrorIsNil)
	err = s.principal.AssignToMachine(s.machine)
	c.Assert(err, jc.ErrorIsNil)
	relUnit, err := rel.Unit(s.principal)
	c.Assert(err, jc.ErrorIsNil)
	err = relUnit.EnterScope(nil)
	c.Assert(err, jc.ErrorIsNil)
	s.subordinate, err = s.State.Unit("logging/0")
	c.Assert(err, jc.ErrorIsNil)

	// Create the deployer facade.
	s.st = deployer.NewState(s.stateAPI)
	c.Assert(s.st, gc.NotNil)
}

// Note: This is really meant as a unit-test, this isn't a test that
// should need all of the setup we have for this test suite
func (s *deployerSuite) TestNew(c *gc.C) {
	deployer := deployer.NewState(s.stateAPI)
	c.Assert(deployer, gc.NotNil)
}

func (s *deployerSuite) assertUnauthorized(c *gc.C, err error) {
	c.Assert(err, gc.ErrorMatches, "permission denied")
	c.Assert(err, jc.Satisfies, params.IsCodeUnauthorized)
}

func (s *deployerSuite) TestWatchUnitsWrongMachine(c *gc.C) {
	// Try with a non-existent machine tag.
	machine, err := s.st.Machine(names.NewMachineTag("42"))
	c.Assert(err, jc.ErrorIsNil)
	w, err := machine.WatchUnits()
	s.assertUnauthorized(c, err)
	c.Assert(w, gc.IsNil)
}

func (s *deployerSuite) TestWatchUnits(c *gc.C) {
	// TODO(dfc) fix state.Machine to return a MachineTag
	machine, err := s.st.Machine(s.machine.Tag().(names.MachineTag))
	c.Assert(err, jc.ErrorIsNil)
	w, err := machine.WatchUnits()
	c.Assert(err, jc.ErrorIsNil)
	wc := watchertest.NewStringsWatcherC(c, w, s.BackingState.StartSync)
	defer wc.AssertStops()

	// Initial event.
	wc.AssertChange("mysql/0", "logging/0")
	wc.AssertNoChange()

	// Change something other than the lifecycle and make sure it's
	// not detected.
	err = s.subordinate.SetPassword("foo")
	c.Assert(err, gc.ErrorMatches, "password is only 3 bytes long, and is not a valid Agent password")
	wc.AssertNoChange()

	err = s.subordinate.SetPassword("foo-12345678901234567890")
	c.Assert(err, jc.ErrorIsNil)
	wc.AssertNoChange()

	// Make the subordinate dead and check it's detected.
	err = s.subordinate.EnsureDead()
	c.Assert(err, jc.ErrorIsNil)
	wc.AssertChange("logging/0")
	wc.AssertNoChange()
}

func (s *deployerSuite) TestUnit(c *gc.C) {
	// Try getting a missing unit and an invalid tag.
	unit, err := s.st.Unit(names.NewUnitTag("foo/42"))
	s.assertUnauthorized(c, err)
	c.Assert(unit, gc.IsNil)

	// Try getting a unit we're not responsible for.
	// First create a new machine and deploy another unit there.
	machine, err := s.State.AddMachine("quantal", state.JobHostUnits)
	c.Assert(err, jc.ErrorIsNil)
	principal1, err := s.app0.AddUnit(state.AddUnitParams{})
	c.Assert(err, jc.ErrorIsNil)
	err = principal1.AssignToMachine(machine)
	c.Assert(err, jc.ErrorIsNil)
	unit, err = s.st.Unit(principal1.Tag().(names.UnitTag))
	s.assertUnauthorized(c, err)
	c.Assert(unit, gc.IsNil)

	// Get the principal and subordinate we're responsible for.
	unit, err = s.st.Unit(s.principal.Tag().(names.UnitTag))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(unit.Name(), gc.Equals, "mysql/0")
	unit, err = s.st.Unit(s.subordinate.Tag().(names.UnitTag))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(unit.Name(), gc.Equals, "logging/0")
}

func (s *deployerSuite) TestUnitLifeRefresh(c *gc.C) {
	unit, err := s.st.Unit(s.subordinate.Tag().(names.UnitTag))
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(unit.Life(), gc.Equals, life.Alive)

	// Now make it dead and check again, then refresh and check.
	err = s.subordinate.EnsureDead()
	c.Assert(err, jc.ErrorIsNil)
	err = s.subordinate.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(s.subordinate.Life(), gc.Equals, state.Dead)
	c.Assert(unit.Life(), gc.Equals, life.Alive)
	err = unit.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(unit.Life(), gc.Equals, life.Dead)
}

func (s *deployerSuite) TestUnitRemove(c *gc.C) {
	unit, err := s.st.Unit(s.principal.Tag().(names.UnitTag))
	c.Assert(err, jc.ErrorIsNil)

	// It fails because the entity is still alive.
	// And EnsureDead will fail because there is a subordinate.
	err = unit.Remove()
	c.Assert(err, gc.ErrorMatches, `cannot remove entity "unit-mysql-0": still alive`)
	c.Assert(params.ErrCode(err), gc.Equals, "")

	// With the subordinate it also fails due to it being alive.
	unit, err = s.st.Unit(s.subordinate.Tag().(names.UnitTag))
	c.Assert(err, jc.ErrorIsNil)
	err = unit.Remove()
	c.Assert(err, gc.ErrorMatches, `cannot remove entity "unit-logging-0": still alive`)
	c.Assert(params.ErrCode(err), gc.Equals, "")

	// Make it dead first and try again.
	err = s.subordinate.EnsureDead()
	c.Assert(err, jc.ErrorIsNil)
	err = unit.Remove()
	c.Assert(err, jc.ErrorIsNil)

	// Verify it's gone.
	err = unit.Refresh()
	s.assertUnauthorized(c, err)
	unit, err = s.st.Unit(s.subordinate.Tag().(names.UnitTag))
	s.assertUnauthorized(c, err)
	c.Assert(unit, gc.IsNil)
}

func (s *deployerSuite) TestUnitSetPassword(c *gc.C) {
	unit, err := s.st.Unit(s.principal.Tag().(names.UnitTag))
	c.Assert(err, jc.ErrorIsNil)

	// Change the principal's password and verify.
	err = unit.SetPassword("foobar-12345678901234567890")
	c.Assert(err, jc.ErrorIsNil)
	err = s.principal.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(s.principal.PasswordValid("foobar-12345678901234567890"), jc.IsTrue)

	// Then the subordinate.
	unit, err = s.st.Unit(s.subordinate.Tag().(names.UnitTag))
	c.Assert(err, jc.ErrorIsNil)
	err = unit.SetPassword("phony-12345678901234567890")
	c.Assert(err, jc.ErrorIsNil)
	err = s.subordinate.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(s.subordinate.PasswordValid("phony-12345678901234567890"), jc.IsTrue)
}

func (s *deployerSuite) TestUnitSetStatus(c *gc.C) {
	unit, err := s.st.Unit(s.principal.Tag().(names.UnitTag))
	c.Assert(err, jc.ErrorIsNil)
	err = unit.SetStatus(status.Blocked, "waiting", map[string]interface{}{"foo": "bar"})
	c.Assert(err, jc.ErrorIsNil)

	stateUnit, err := s.BackingState.Unit(unit.Name())
	c.Assert(err, jc.ErrorIsNil)
	sInfo, err := stateUnit.Status()
	c.Assert(err, jc.ErrorIsNil)
	sInfo.Since = nil
	c.Assert(sInfo, jc.DeepEquals, status.StatusInfo{
		Status:  status.Blocked,
		Message: "waiting",
		Data:    map[string]interface{}{"foo": "bar"},
	})
}
