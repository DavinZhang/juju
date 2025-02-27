// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package testing

import (
	"github.com/juju/testing"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/state/backups"
)

// BaseSuite is the  base suite for backups testing.
type BaseSuite struct {
	testing.IsolationSuite
	// Meta is a Metadata with standard test values.
	Meta *backups.Metadata
}

func (s *BaseSuite) SetUpTest(c *gc.C) {
	s.IsolationSuite.SetUpTest(c)
	s.Meta = NewMetadata()
}
