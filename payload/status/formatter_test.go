// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package status_test

import (
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/payload/status"
)

var _ = gc.Suite(&formatterSuite{})

type formatterSuite struct {
	testing.IsolationSuite
}

func (s *formatterSuite) TestFormatPayloadOkay(c *gc.C) {
	payload := status.NewPayload("spam", "a-application", 1, 0)
	payload.Labels = []string{"a-tag"}
	formatted := status.FormatPayload(payload)

	c.Check(formatted, jc.DeepEquals, status.FormattedPayload{
		Unit:    "a-application/0",
		Machine: "1",
		ID:      "idspam",
		Type:    "docker",
		Class:   "spam",
		Labels:  []string{"a-tag"},
		Status:  "running",
	})
}
