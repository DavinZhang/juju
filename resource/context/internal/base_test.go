// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package internal_test

import (
	"io"
	"time"

	charmresource "github.com/juju/charm/v9/resource"
	"github.com/juju/testing"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/resource"
	"github.com/DavinZhang/juju/resource/resourcetesting"
)

func newCharmResource(c *gc.C, stub *testing.Stub, name, content string, resType charmresource.Type) (resource.Resource, io.ReadCloser) {
	opened := resourcetesting.NewResource(c, stub, name, "a-application", content)
	res := opened.Resource
	res.Type = resType
	if content != "" {
		return res, opened.ReadCloser
	}
	res.Username = ""
	res.Timestamp = time.Time{}
	return res, nil
}

func newResource(c *gc.C, stub *testing.Stub, name, content string) (resource.Resource, io.ReadCloser) {
	return newCharmResource(c, stub, name, content, charmresource.TypeFile)
}

func newDockerResource(c *gc.C, stub *testing.Stub, name, content string) (resource.Resource, io.ReadCloser) {
	return newCharmResource(c, stub, name, content, charmresource.TypeContainerImage)
}
