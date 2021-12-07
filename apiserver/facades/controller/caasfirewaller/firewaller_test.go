// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasfirewaller_test

import (
	"github.com/juju/charm/v9"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/worker/v3/workertest"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/apiserver/common"
	charmscommon "github.com/DavinZhang/juju/apiserver/common/charms"
	"github.com/DavinZhang/juju/apiserver/facade"
	"github.com/DavinZhang/juju/apiserver/facades/controller/caasfirewaller"
	"github.com/DavinZhang/juju/apiserver/params"
	apiservertesting "github.com/DavinZhang/juju/apiserver/testing"
	"github.com/DavinZhang/juju/core/life"
	"github.com/DavinZhang/juju/state"
	statetesting "github.com/DavinZhang/juju/state/testing"
	coretesting "github.com/DavinZhang/juju/testing"
)

type firewallerBaseSuite struct {
	coretesting.BaseSuite

	st                  *mockState
	applicationsChanges chan []string
	openPortsChanges    chan []string
	appExposedChanges   chan struct{}

	resources  *common.Resources
	authorizer *apiservertesting.FakeAuthorizer
	facade     facadeCommon

	newFunc func(c *gc.C, resources facade.Resources,
		authorizer facade.Authorizer,
		st *mockState,
	) (facadeCommon, error)
}

type firewallerLegacySuite struct {
	firewallerBaseSuite
}

var _ = gc.Suite(&firewallerLegacySuite{
	firewallerBaseSuite: firewallerBaseSuite{
		newFunc: func(c *gc.C, resources facade.Resources,
			authorizer facade.Authorizer,
			st *mockState,
		) (facadeCommon, error) {
			commonState := &mockCommonStateShim{st}
			commonCharmsAPI, err := charmscommon.NewCharmInfoAPI(commonState, authorizer)
			c.Assert(err, jc.ErrorIsNil)
			appCharmInfoAPI, err := charmscommon.NewApplicationCharmInfoAPI(commonState, authorizer)
			c.Assert(err, jc.ErrorIsNil)
			return caasfirewaller.NewFacadeLegacyForTest(
				resources,
				authorizer,
				st,
				commonCharmsAPI,
				appCharmInfoAPI,
			)
		},
	},
})

type firewallerSidecarSuite struct {
	firewallerBaseSuite

	facade facadeSidecar
}

var _ = gc.Suite(&firewallerSidecarSuite{
	firewallerBaseSuite: firewallerBaseSuite{
		newFunc: func(c *gc.C, resources facade.Resources,
			authorizer facade.Authorizer,
			st *mockState,
		) (facadeCommon, error) {
			commonState := &mockCommonStateShim{st}
			commonCharmsAPI, err := charmscommon.NewCharmInfoAPI(commonState, authorizer)
			c.Assert(err, jc.ErrorIsNil)
			appCharmInfoAPI, err := charmscommon.NewApplicationCharmInfoAPI(commonState, authorizer)
			c.Assert(err, jc.ErrorIsNil)
			return caasfirewaller.NewFacadeSidecarForTest(
				resources,
				authorizer,
				st,
				commonCharmsAPI,
				appCharmInfoAPI,
			)
		},
	},
})

func (s *firewallerSidecarSuite) SetUpTest(c *gc.C) {
	s.firewallerBaseSuite.SetUpTest(c)

	// charm.FormatV2.
	s.st.application.charm.manifest.Bases = []charm.Base{
		{
			Name: "ubuntu",
			Channel: charm.Channel{
				Risk:  "stable",
				Track: "20.04",
			},
		},
	}

	var ok bool
	s.facade, ok = s.firewallerBaseSuite.facade.(facadeSidecar)
	c.Assert(ok, jc.IsTrue)
}

func (s *firewallerSidecarSuite) TestWatchOpenedPorts(c *gc.C) {
	openPortsChanges := []string{"port1", "port2"}
	s.openPortsChanges <- openPortsChanges

	results, err := s.facade.WatchOpenedPorts(params.Entities{
		Entities: []params.Entity{{
			Tag: "model-deadbeef-0bad-400d-8000-4b1d0d06f00d",
		}},
	})
	c.Assert(err, jc.ErrorIsNil)
	result := results.Results[0]
	c.Assert(result.Error, gc.IsNil)
	c.Assert(result.StringsWatcherId, gc.Equals, "1")
	c.Assert(result.Changes, jc.DeepEquals, openPortsChanges)
}

type facadeCommon interface {
	IsExposed(args params.Entities) (params.BoolResults, error)
	ApplicationsConfig(args params.Entities) (params.ApplicationGetConfigResults, error)
	WatchApplications() (params.StringsWatchResult, error)
	Life(args params.Entities) (params.LifeResults, error)
	Watch(args params.Entities) (params.NotifyWatchResults, error)
	ApplicationCharmInfo(args params.Entity) (params.Charm, error)
}

type facadeSidecar interface {
	facadeCommon
	WatchOpenedPorts(args params.Entities) (params.StringsWatchResults, error)
}

func (s *firewallerBaseSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)

	s.applicationsChanges = make(chan []string, 1)
	s.appExposedChanges = make(chan struct{}, 1)
	s.openPortsChanges = make(chan []string, 1)
	appExposedWatcher := statetesting.NewMockNotifyWatcher(s.appExposedChanges)
	s.st = &mockState{
		application: mockApplication{
			life:    state.Alive,
			watcher: appExposedWatcher,
			charm: mockCharm{
				meta: &charm.Meta{
					Deployment: &charm.Deployment{},
				},
				manifest: &charm.Manifest{},
				url:      &charm.URL{Schema: "cs", Name: "gitlab", Revision: -1},
			},
		},
		applicationsWatcher: statetesting.NewMockStringsWatcher(s.applicationsChanges),
		openPortsWatcher:    statetesting.NewMockStringsWatcher(s.openPortsChanges),
		appExposedWatcher:   appExposedWatcher,
	}
	s.AddCleanup(func(c *gc.C) { workertest.DirtyKill(c, s.st.applicationsWatcher) })
	s.AddCleanup(func(c *gc.C) { workertest.DirtyKill(c, s.st.openPortsWatcher) })
	s.AddCleanup(func(c *gc.C) { workertest.DirtyKill(c, s.st.appExposedWatcher) })

	s.resources = common.NewResources()
	s.authorizer = &apiservertesting.FakeAuthorizer{
		Tag:        names.NewMachineTag("0"),
		Controller: true,
	}

	facade, err := s.newFunc(c, s.resources, s.authorizer, s.st)
	c.Assert(err, jc.ErrorIsNil)
	s.facade = facade
}

func (s *firewallerBaseSuite) TestPermission(c *gc.C) {
	s.authorizer = &apiservertesting.FakeAuthorizer{
		Tag: names.NewMachineTag("0"),
	}
	_, err := s.newFunc(c, s.resources, s.authorizer, s.st)
	c.Assert(err, gc.ErrorMatches, "permission denied")
}

func (s *firewallerBaseSuite) TestWatchApplications(c *gc.C) {
	applicationNames := []string{"db2", "hadoop"}
	s.applicationsChanges <- applicationNames
	result, err := s.facade.WatchApplications()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result.Error, gc.IsNil)
	c.Assert(result.StringsWatcherId, gc.Equals, "1")
	c.Assert(result.Changes, jc.DeepEquals, applicationNames)
}

func (s *firewallerBaseSuite) TestWatchApplication(c *gc.C) {
	s.appExposedChanges <- struct{}{}

	results, err := s.facade.Watch(params.Entities{
		Entities: []params.Entity{
			{Tag: "application-gitlab"},
			{Tag: "unit-gitlab-0"},
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results.Results, gc.HasLen, 2)
	c.Assert(results.Results[0].Error, gc.IsNil)
	c.Assert(results.Results[1].Error, jc.DeepEquals, &params.Error{
		Message: "permission denied",
		Code:    "unauthorized access",
	})

	c.Assert(results.Results[0].NotifyWatcherId, gc.Equals, "1")
	resource := s.resources.Get("1")
	c.Assert(resource, gc.Equals, s.st.appExposedWatcher)
}

func (s *firewallerBaseSuite) TestIsExposed(c *gc.C) {
	s.st.application.exposed = true
	results, err := s.facade.IsExposed(params.Entities{
		Entities: []params.Entity{
			{Tag: "application-gitlab"},
			{Tag: "unit-gitlab-0"},
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, params.BoolResults{
		Results: []params.BoolResult{{
			Result: true,
		}, {
			Error: &params.Error{
				Message: `"unit-gitlab-0" is not a valid application tag`,
			},
		}},
	})
}

func (s *firewallerBaseSuite) TestLife(c *gc.C) {
	results, err := s.facade.Life(params.Entities{
		Entities: []params.Entity{
			{Tag: "application-gitlab"},
			{Tag: "machine-0"},
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, params.LifeResults{
		Results: []params.LifeResult{{
			Life: life.Alive,
		}, {
			Error: &params.Error{
				Code:    "unauthorized access",
				Message: "permission denied",
			},
		}},
	})
}

func (s *firewallerBaseSuite) TestApplicationConfig(c *gc.C) {
	results, err := s.facade.ApplicationsConfig(params.Entities{
		Entities: []params.Entity{
			{Tag: "application-gitlab"},
			{Tag: "unit-gitlab-0"},
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results.Results, gc.HasLen, 2)
	c.Assert(results.Results[0].Error, gc.IsNil)
	c.Assert(results.Results[1].Error, jc.DeepEquals, &params.Error{
		Message: `"unit-gitlab-0" is not a valid application tag`,
	})
	c.Assert(results.Results[0].Config, jc.DeepEquals, map[string]interface{}{"foo": "bar"})
}
