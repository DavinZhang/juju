// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package featuretests

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	jujutesting "github.com/DavinZhang/juju/juju/testing"
	"github.com/DavinZhang/juju/testing/factory"
)

type cmdSetSeriesSuite struct {
	jujutesting.JujuConnSuite
}

func (s *cmdSetSeriesSuite) TestSetApplicationSeries(c *gc.C) {
	charm := s.Factory.MakeCharm(c, &factory.CharmParams{Name: "multi-series", URL: "local:quantal/multi-series-1"})
	app := s.Factory.MakeApplication(c, &factory.ApplicationParams{Charm: charm})
	s.Factory.MakeUnit(c, &factory.UnitParams{Application: app, SetCharmURL: true})
	c.Assert(app.Series(), gc.Equals, "quantal")
	runCommandExpectSuccess(c, "set-series", "multi-series", "trusty")
	err := app.Refresh()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(app.Series(), gc.Equals, "trusty")
}
