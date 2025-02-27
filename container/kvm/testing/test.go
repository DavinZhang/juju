// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Functions defined in this file should *ONLY* be used for testing.  These
// functions are exported for testing purposes only, and shouldn't be called
// from code that isn't in a test file.

package testing

import (
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/container"
	"github.com/DavinZhang/juju/container/kvm"
	"github.com/DavinZhang/juju/container/kvm/mock"
	"github.com/DavinZhang/juju/testing"
)

// TestSuite replaces the kvm factory that the manager uses with a mock
// implementation.
type TestSuite struct {
	testing.BaseSuite
	ContainerFactory mock.ContainerFactory
	ContainerDir     string
	RemovedDir       string
}

func (s *TestSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)
	s.ContainerDir = c.MkDir()
	s.PatchValue(&container.ContainerDir, s.ContainerDir)
	s.RemovedDir = c.MkDir()
	s.PatchValue(&container.RemovedContainerDir, s.RemovedDir)
	s.ContainerFactory = mock.MockFactory()
	s.PatchValue(&kvm.KVMObjectFactory, s.ContainerFactory)
}
