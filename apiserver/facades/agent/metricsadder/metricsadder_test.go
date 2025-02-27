// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package metricsadder_test

import (
	"time"

	"github.com/juju/errors"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v2"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/apiserver/common"
	"github.com/DavinZhang/juju/apiserver/facades/agent/metricsadder"
	"github.com/DavinZhang/juju/apiserver/params"
	apiservertesting "github.com/DavinZhang/juju/apiserver/testing"
	jujutesting "github.com/DavinZhang/juju/juju/testing"
	"github.com/DavinZhang/juju/state"
	jujuFactory "github.com/DavinZhang/juju/testing/factory"
)

var _ = gc.Suite(&metricsAdderSuite{})

type metricsAdderSuite struct {
	jujutesting.JujuConnSuite

	authorizer apiservertesting.FakeAuthorizer
	resources  *common.Resources

	machine0       *state.Machine
	machine1       *state.Machine
	mysqlService   *state.Application
	mysql          *state.Application
	mysqlUnit      *state.Unit
	meteredService *state.Application
	meteredCharm   *state.Charm
	meteredUnit    *state.Unit

	adder metricsadder.MetricsAdder
}

func (s *metricsAdderSuite) SetUpTest(c *gc.C) {
	s.JujuConnSuite.SetUpTest(c)
	s.machine0 = s.Factory.MakeMachine(c, &jujuFactory.MachineParams{
		Series: "quantal",
		Jobs:   []state.MachineJob{state.JobHostUnits, state.JobManageModel},
	})
	s.machine1 = s.Factory.MakeMachine(c, &jujuFactory.MachineParams{
		Series: "quantal",
		Jobs:   []state.MachineJob{state.JobHostUnits},
	})
	mysqlCharm := s.Factory.MakeCharm(c, &jujuFactory.CharmParams{
		Name: "mysql",
	})
	s.mysql = s.Factory.MakeApplication(c, &jujuFactory.ApplicationParams{
		Name:  "mysql",
		Charm: mysqlCharm,
	})
	s.mysqlUnit = s.Factory.MakeUnit(c, &jujuFactory.UnitParams{
		Application: s.mysql,
		Machine:     s.machine0,
	})

	s.meteredCharm = s.Factory.MakeCharm(c, &jujuFactory.CharmParams{
		Name: "metered",
		URL:  "cs:quantal/metered",
	})
	s.meteredService = s.Factory.MakeApplication(c, &jujuFactory.ApplicationParams{
		Charm: s.meteredCharm,
	})
	s.meteredUnit = s.Factory.MakeUnit(c, &jujuFactory.UnitParams{
		Application: s.meteredService,
		SetCharmURL: true,
		Machine:     s.machine1,
	})

	// Create a FakeAuthorizer so we can check permissions,
	// set up assuming unit 0 has logged in.
	s.authorizer = apiservertesting.FakeAuthorizer{
		Tag: names.NewMachineTag("1"),
	}

	// Create the resource registry separately to track invocations to
	// Register.
	s.resources = common.NewResources()
	s.AddCleanup(func(_ *gc.C) { s.resources.StopAll() })

	adder, err := metricsadder.NewMetricsAdderAPI(s.State, s.resources, s.authorizer)
	c.Assert(err, jc.ErrorIsNil)
	s.adder = adder
}

func (s *metricsAdderSuite) TestAddMetricsBatch(c *gc.C) {
	metrics := []params.Metric{{
		Key: "pings", Value: "5", Time: time.Now().UTC(),
	}, {
		Key: "pongs", Value: "6", Time: time.Now().UTC(), Labels: map[string]string{"foo": "bar"},
	}}
	uuid := utils.MustNewUUID().String()

	result, err := s.adder.AddMetricBatches(params.MetricBatchParams{
		Batches: []params.MetricBatchParam{{
			Tag: s.meteredUnit.Tag().String(),
			Batch: params.MetricBatch{
				UUID:     uuid,
				CharmURL: s.meteredCharm.URL().String(),
				Created:  time.Now(),
				Metrics:  metrics,
			}}}},
	)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, gc.DeepEquals, params.ErrorResults{
		Results: []params.ErrorResult{{nil}},
	})

	batches, err := s.State.AllMetricBatches()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(batches, gc.HasLen, 1)
	batch := batches[0]
	c.Assert(batch.UUID(), gc.Equals, uuid)
	c.Assert(batch.CharmURL(), gc.Equals, s.meteredCharm.URL().String())
	c.Assert(batch.Unit(), gc.Equals, s.meteredUnit.Name())
	storedMetrics := batch.Metrics()
	c.Assert(storedMetrics, gc.HasLen, 2)
	c.Assert(storedMetrics[0].Key, gc.Equals, metrics[0].Key)
	c.Assert(storedMetrics[0].Value, gc.Equals, metrics[0].Value)
	c.Assert(storedMetrics[0].Labels, gc.DeepEquals, metrics[0].Labels)
	c.Assert(storedMetrics[1].Key, gc.Equals, metrics[1].Key)
	c.Assert(storedMetrics[1].Value, gc.Equals, metrics[1].Value)
	c.Assert(storedMetrics[1].Labels, gc.DeepEquals, metrics[1].Labels)
}

func (s *metricsAdderSuite) TestAddMetricsBatchNoCharmURL(c *gc.C) {
	metrics := []params.Metric{{Key: "pings", Value: "5", Time: time.Now().UTC()}}
	uuid := utils.MustNewUUID().String()

	result, err := s.adder.AddMetricBatches(params.MetricBatchParams{
		Batches: []params.MetricBatchParam{{
			Tag: s.meteredUnit.Tag().String(),
			Batch: params.MetricBatch{
				UUID:     uuid,
				CharmURL: s.meteredCharm.URL().String(),
				Created:  time.Now(),
				Metrics:  metrics,
			}}}})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, gc.DeepEquals, params.ErrorResults{
		Results: []params.ErrorResult{{nil}},
	})

	batches, err := s.State.AllMetricBatches()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(batches, gc.HasLen, 1)
	batch := batches[0]
	c.Assert(batch.UUID(), gc.Equals, uuid)
	c.Assert(batch.CharmURL(), gc.Equals, s.meteredCharm.URL().String())
	c.Assert(batch.Unit(), gc.Equals, s.meteredUnit.Name())
	storedMetrics := batch.Metrics()
	c.Assert(storedMetrics, gc.HasLen, 1)
	c.Assert(storedMetrics[0].Key, gc.Equals, metrics[0].Key)
	c.Assert(storedMetrics[0].Value, gc.Equals, metrics[0].Value)
}

func (s *metricsAdderSuite) TestAddMetricsBatchDiffTag(c *gc.C) {
	metrics := []params.Metric{{Key: "pings", Value: "5", Time: time.Now().UTC()}}
	uuid := utils.MustNewUUID().String()

	tests := []struct {
		about  string
		tag    string
		expect string
	}{{
		about:  "unknown unit",
		tag:    names.NewUnitTag("unknownservice/11").String(),
		expect: "unit \"unknownservice/11\" not found",
	}, {
		about:  "user tag",
		tag:    names.NewLocalUserTag("admin").String(),
		expect: `"user-admin" is not a valid unit tag`,
	}, {
		about:  "machine tag",
		tag:    names.NewMachineTag("0").String(),
		expect: `"machine-0" is not a valid unit tag`,
	}}

	for i, test := range tests {
		c.Logf("test %d: %s -> %s", i, test.about, test.tag)
		result, err := s.adder.AddMetricBatches(params.MetricBatchParams{
			Batches: []params.MetricBatchParam{{
				Tag: test.tag,
				Batch: params.MetricBatch{
					UUID:     uuid,
					CharmURL: s.meteredCharm.URL().String(),
					Created:  time.Now(),
					Metrics:  metrics,
				}}}})
		c.Assert(err, jc.ErrorIsNil)
		if test.expect == "" {
			c.Assert(result.OneError(), jc.ErrorIsNil)
		} else {
			c.Assert(result.OneError(), gc.ErrorMatches, test.expect)
		}

		batches, err := s.State.AllMetricBatches()
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(batches, gc.HasLen, 0)

		_, err = s.State.MetricBatch(uuid)
		c.Assert(err, jc.Satisfies, errors.IsNotFound)
	}
}

func (s *metricsAdderSuite) TestNewMetricsAdderAPIRefusesNonAgent(c *gc.C) {
	tests := []struct {
		tag           names.Tag
		controller    bool
		expectedError string
	}{
		// TODO(cmars): unit agent should get permission denied when callers are
		// moved to machine agent.
		{names.NewUnitTag("mysql/0"), false, ""},

		{names.NewLocalUserTag("admin"), true, "permission denied"},
		{names.NewMachineTag("0"), false, ""},
		{names.NewMachineTag("0"), true, ""},
	}
	for i, test := range tests {
		c.Logf("test %d", i)

		anAuthoriser := s.authorizer
		anAuthoriser.Controller = test.controller
		anAuthoriser.Tag = test.tag
		endPoint, err := metricsadder.NewMetricsAdderAPI(s.State, nil, anAuthoriser)
		if test.expectedError == "" {
			c.Assert(err, jc.ErrorIsNil)
			c.Assert(endPoint, gc.NotNil)
		} else {
			c.Assert(err, gc.ErrorMatches, test.expectedError)
			c.Assert(endPoint, gc.IsNil)
		}
	}
}
