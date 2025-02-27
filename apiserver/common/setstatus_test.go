// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common_test

import (
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/apiserver/common"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/core/status"
	"github.com/DavinZhang/juju/state"
	"github.com/DavinZhang/juju/testing/factory"
)

type statusSetterSuite struct {
	statusBaseSuite
	setter *common.StatusSetter
}

var _ = gc.Suite(&statusSetterSuite{})

func (s *statusSetterSuite) SetUpTest(c *gc.C) {
	s.statusBaseSuite.SetUpTest(c)

	s.setter = common.NewStatusSetter(s.State, func() (common.AuthFunc, error) {
		return s.authFunc, nil
	})
}

func (s *statusSetterSuite) TestUnauthorized(c *gc.C) {
	tag := names.NewMachineTag("42")
	s.badTag = tag
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    tag.String(),
		Status: status.Executing.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, jc.Satisfies, params.IsCodeUnauthorized)
}

func (s *statusSetterSuite) TestNotATag(c *gc.C) {
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    "not a tag",
		Status: status.Executing.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, gc.ErrorMatches, `"not a tag" is not a valid tag`)
}

func (s *statusSetterSuite) TestNotFound(c *gc.C) {
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    names.NewMachineTag("42").String(),
		Status: status.Down.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, jc.Satisfies, params.IsCodeNotFound)
}

func (s *statusSetterSuite) TestSetMachineStatus(c *gc.C) {
	machine := s.Factory.MakeMachine(c, nil)
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    machine.Tag().String(),
		Status: status.Started.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, gc.IsNil)

	err = machine.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	machineStatus, err := machine.Status()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(machineStatus.Status, gc.Equals, status.Started)
}

func (s *statusSetterSuite) TestSetUnitStatus(c *gc.C) {
	// The status has to be a valid workload status, because get status
	// on the unit returns the workload status not the agent status as it
	// does on a machine.
	unit := s.Factory.MakeUnit(c, &factory.UnitParams{Status: &status.StatusInfo{
		Status: status.Maintenance,
	}})
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    unit.Tag().String(),
		Status: status.Active.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, gc.IsNil)

	err = unit.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	unitStatus, err := unit.Status()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(unitStatus.Status, gc.Equals, status.Active)
}

func (s *statusSetterSuite) TestSetServiceStatus(c *gc.C) {
	// Calls to set the status of a service should be going through the
	// ServiceStatusSetter that checks for leadership, so permission denied
	// here.
	service := s.Factory.MakeApplication(c, &factory.ApplicationParams{Status: &status.StatusInfo{
		Status: status.Maintenance,
	}})
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    service.Tag().String(),
		Status: status.Active.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, jc.Satisfies, params.IsCodeUnauthorized)

	err = service.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	serviceStatus, err := service.Status()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(serviceStatus.Status, gc.Equals, status.Maintenance)
}

func (s *statusSetterSuite) TestBulk(c *gc.C) {
	s.badTag = names.NewMachineTag("42")
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    s.badTag.String(),
		Status: status.Active.String(),
	}, {
		Tag:    "bad-tag",
		Status: status.Active.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 2)
	c.Assert(result.Results[0].Error, jc.Satisfies, params.IsCodeUnauthorized)
	c.Assert(result.Results[1].Error, gc.ErrorMatches, `"bad-tag" is not a valid tag`)
}

type serviceStatusSetterSuite struct {
	statusBaseSuite
	setter *common.ApplicationStatusSetter
}

var _ = gc.Suite(&serviceStatusSetterSuite{})

func (s *serviceStatusSetterSuite) SetUpTest(c *gc.C) {
	s.statusBaseSuite.SetUpTest(c)

	s.setter = common.NewApplicationStatusSetter(s.State, func() (common.AuthFunc, error) {
		return s.authFunc, nil
	}, s.leadershipChecker)
}

func (s *serviceStatusSetterSuite) TestUnauthorized(c *gc.C) {
	// Machines are unauthorized since they are not units
	tag := names.NewUnitTag("foo/0")
	s.badTag = tag
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    tag.String(),
		Status: status.Active.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, jc.Satisfies, params.IsCodeUnauthorized)
}

func (s *serviceStatusSetterSuite) TestNotATag(c *gc.C) {
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    "not a tag",
		Status: status.Active.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, gc.ErrorMatches, `"not a tag" is not a valid tag`)
}

func (s *serviceStatusSetterSuite) TestNotFound(c *gc.C) {
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    names.NewUnitTag("foo/0").String(),
		Status: status.Active.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, jc.Satisfies, params.IsCodeNotFound)
}

func (s *serviceStatusSetterSuite) TestSetMachineStatus(c *gc.C) {
	machine := s.Factory.MakeMachine(c, nil)
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    machine.Tag().String(),
		Status: status.Active.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	// Can't call set service status on a machine.
	c.Assert(result.Results[0].Error, jc.Satisfies, params.IsCodeUnauthorized)
}

func (s *serviceStatusSetterSuite) TestSetServiceStatus(c *gc.C) {
	// TODO: the correct way to fix this is to have the authorizer on the
	// simple status setter to check to see if the unit (authTag) is a leader
	// and able to set the service status. However, that is for another day.
	service := s.Factory.MakeApplication(c, &factory.ApplicationParams{Status: &status.StatusInfo{
		Status: status.Maintenance,
	}})
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    service.Tag().String(),
		Status: status.Active.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	// Can't call set service status on a service. Weird I know, but the only
	// way is to go through the unit leader.
	c.Assert(result.Results[0].Error, jc.Satisfies, params.IsCodeUnauthorized)
}

func (s *serviceStatusSetterSuite) TestSetUnitStatusNotLeader(c *gc.C) {
	// If the unit isn't the leader, it can't set it.
	s.leadershipChecker.isLeader = false
	unit := s.Factory.MakeUnit(c, &factory.UnitParams{Status: &status.StatusInfo{
		Status: status.Maintenance,
	}})
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    unit.Tag().String(),
		Status: status.Active.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	status := result.Results[0]
	c.Assert(status.Error, gc.ErrorMatches, "not leader")
}

func (s *serviceStatusSetterSuite) TestSetUnitStatusIsLeader(c *gc.C) {
	service := s.Factory.MakeApplication(c, &factory.ApplicationParams{Status: &status.StatusInfo{
		Status: status.Maintenance,
	}})
	unit := s.Factory.MakeUnit(c, &factory.UnitParams{
		Application: service,
		Status: &status.StatusInfo{
			Status: status.Maintenance,
		}})
	// No need to claim leadership - the checker passed in in setup
	// always returns true.
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    unit.Tag().String(),
		Status: status.Active.String(),
	}}})

	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 1)
	c.Assert(result.Results[0].Error, gc.IsNil)

	err = service.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	unitStatus, err := service.Status()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(unitStatus.Status, gc.Equals, status.Active)
}

func (s *serviceStatusSetterSuite) TestBulk(c *gc.C) {
	s.badTag = names.NewMachineTag("42")
	machine := s.Factory.MakeMachine(c, nil)
	result, err := s.setter.SetStatus(params.SetStatus{[]params.EntityStatusArgs{{
		Tag:    s.badTag.String(),
		Status: status.Active.String(),
	}, {
		Tag:    machine.Tag().String(),
		Status: status.Active.String(),
	}, {
		Tag:    "bad-tag",
		Status: status.Active.String(),
	}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Results, gc.HasLen, 3)
	c.Assert(result.Results[0].Error, jc.Satisfies, params.IsCodeUnauthorized)
	c.Assert(result.Results[1].Error, jc.Satisfies, params.IsCodeUnauthorized)
	c.Assert(result.Results[2].Error, gc.ErrorMatches, `"bad-tag" is not a valid tag`)
}

type unitAgentFinderSuite struct{}

var _ = gc.Suite(&unitAgentFinderSuite{})

func (unitAgentFinderSuite) TestFindEntity(c *gc.C) {
	f := fakeEntityFinder{
		unit: fakeUnit{
			agent: &state.UnitAgent{},
		},
	}
	ua := &common.UnitAgentFinder{f}
	entity, err := ua.FindEntity(names.NewUnitTag("unit/0"))
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(entity, gc.DeepEquals, f.unit.agent)
}

func (unitAgentFinderSuite) TestFindEntityBadTag(c *gc.C) {
	ua := &common.UnitAgentFinder{fakeEntityFinder{}}
	_, err := ua.FindEntity(names.NewApplicationTag("foo"))
	c.Assert(err, gc.ErrorMatches, "unsupported tag.*")
}

func (unitAgentFinderSuite) TestFindEntityErr(c *gc.C) {
	f := fakeEntityFinder{err: errors.Errorf("boo")}
	ua := &common.UnitAgentFinder{f}
	_, err := ua.FindEntity(names.NewUnitTag("unit/0"))
	c.Assert(errors.Cause(err), gc.Equals, f.err)
}

type fakeEntityFinder struct {
	unit fakeUnit
	err  error
}

func (f fakeEntityFinder) FindEntity(tag names.Tag) (state.Entity, error) {
	return f.unit, f.err
}

type fakeUnit struct {
	state.Entity
	agent *state.UnitAgent
}

func (f fakeUnit) Agent() *state.UnitAgent {
	return f.agent
}
