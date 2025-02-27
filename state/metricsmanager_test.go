// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state_test

import (
	"time"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/state"
	testing "github.com/DavinZhang/juju/state/testing"
	coretesting "github.com/DavinZhang/juju/testing"
)

type metricsManagerSuite struct {
	testing.StateSuite
}

var _ = gc.Suite(&metricsManagerSuite{})

func (s *metricsManagerSuite) TestDefaultsWritten(c *gc.C) {
	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.LastSuccessfulSend(), gc.DeepEquals, time.Time{})
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, 0)
	c.Assert(mm.GracePeriod(), gc.Equals, 24*7*time.Hour)
	c.Assert(mm.ModelStatus().Code, gc.Equals, state.MeterNotSet)
}

func (s *metricsManagerSuite) TestNewMetricsManager(c *gc.C) {
	state.SetBeforeHooks(c, s.State, func() {
		_, err := s.State.MetricsManager()
		c.Assert(err, jc.ErrorIsNil)
	})
	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.LastSuccessfulSend(), gc.DeepEquals, time.Time{})
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, 0)
	c.Assert(mm.GracePeriod(), gc.Equals, 24*7*time.Hour)
}

func (s *metricsManagerSuite) TestMetricsManagerCreatesThenReturns(c *gc.C) {
	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	err = mm.IncrementConsecutiveErrors()
	c.Assert(err, jc.ErrorIsNil)
	mm2, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, mm2.ConsecutiveErrors())
}

func (s *metricsManagerSuite) TestSetLastSuccesfulSend(c *gc.C) {
	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	err = mm.IncrementConsecutiveErrors()
	c.Assert(err, jc.ErrorIsNil)
	now := coretesting.ZeroTime().Round(time.Second).UTC()
	err = mm.SetLastSuccessfulSend(now)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.LastSuccessfulSend(), gc.DeepEquals, now)
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, 0)

	m, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(m.LastSuccessfulSend().Equal(now), jc.IsTrue)
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, 0)
}

func (s *metricsManagerSuite) TestIncrementConsecutiveErrors(c *gc.C) {
	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, 0)
	err = mm.IncrementConsecutiveErrors()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.ConsecutiveErrors(), gc.Equals, 1)

	m, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(m.ConsecutiveErrors(), gc.Equals, 1)
}

func (s *metricsManagerSuite) TestSetGracePeriod(c *gc.C) {
	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.GracePeriod(), gc.Equals, time.Hour*24*7)
	err = mm.SetGracePeriod(time.Hour)
	c.Assert(err, jc.ErrorIsNil)

	m, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(m.GracePeriod(), gc.Equals, time.Hour)
}

func (s *metricsManagerSuite) TestNegativeGracePeriod(c *gc.C) {
	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(mm.GracePeriod(), gc.Equals, time.Hour*24*7)
	err = mm.SetGracePeriod(-time.Hour)
	c.Assert(err, gc.ErrorMatches, "grace period can't be negative")
}

func (s *metricsManagerSuite) TestMeterStatus(c *gc.C) {
	err := s.State.SetClockForTesting(s.Clock)
	c.Assert(err, jc.ErrorIsNil)

	mm, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	status := mm.MeterStatus()
	c.Assert(status.Code, gc.Equals, state.MeterGreen)
	c.Assert(mm.ModelStatus().Code, gc.Equals, state.MeterNotSet)
	now := coretesting.NonZeroTime()
	err = mm.SetLastSuccessfulSend(now)
	c.Assert(err, jc.ErrorIsNil)
	for i := 0; i < 3; i++ {
		err := mm.IncrementConsecutiveErrors()
		c.Assert(err, jc.ErrorIsNil)
	}
	status = mm.MeterStatus()
	c.Assert(status.Code, gc.Equals, state.MeterAmber)
	err = mm.SetLastSuccessfulSend(now.Add(-(24 * 7 * time.Hour)))
	c.Assert(err, jc.ErrorIsNil)

	for i := 0; i < 3; i++ {
		err := mm.IncrementConsecutiveErrors()
		c.Assert(err, jc.ErrorIsNil)
	}
	status = mm.MeterStatus()
	c.Assert(status.Code, gc.Equals, state.MeterRed)

	// if we create a new metrics manager, it will pick up
	// model meter status from mongo (RED).
	m, err := s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(m.ModelStatus().Code, gc.Equals, state.MeterRed)

	err = mm.SetLastSuccessfulSend(now)
	c.Assert(err, jc.ErrorIsNil)
	status = mm.MeterStatus()
	c.Assert(status.Code, gc.Equals, state.MeterGreen)

	// if we create a new metrics manager, it will pick up
	// model meter status from mongo (GREEN).
	m, err = s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(m.ModelStatus().Code, gc.Equals, state.MeterGreen)

	err = mm.SetGracePeriod(time.Hour)
	c.Assert(err, jc.ErrorIsNil)

	for i := 0; i < 3; i++ {
		err := mm.IncrementConsecutiveErrors()
		c.Assert(err, jc.ErrorIsNil)
	}

	s.Clock.Advance(24 * time.Hour)

	status = mm.MeterStatus()
	c.Assert(status.Code, gc.Equals, state.MeterRed)

	m, err = s.State.MetricsManager()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(m.MeterStatus(), jc.DeepEquals, status)
}
