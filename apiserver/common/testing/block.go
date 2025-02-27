// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package testing

import (
	"fmt"

	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/api"
	"github.com/DavinZhang/juju/api/block"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/core/model"
)

// BlockHelper helps manage blocks for apiserver tests.
// It provides easy access to switch blocks on
// as well as test whether operations are blocked or not.
type BlockHelper struct {
	apiState api.Connection
	client   *block.Client
}

// NewBlockHelper creates a block switch used in testing
// to manage desired juju blocks.
func NewBlockHelper(st api.Connection) BlockHelper {
	return BlockHelper{
		apiState: st,
		client:   block.NewClient(st),
	}
}

// on switches on desired block and
// asserts that no errors were encountered.
func (s BlockHelper) on(c *gc.C, blockType model.BlockType, msg string) {
	c.Assert(s.client.SwitchBlockOn(fmt.Sprintf("%v", blockType), msg), gc.IsNil)
}

// BlockAllChanges blocks all operations that could change the model.
func (s BlockHelper) BlockAllChanges(c *gc.C, msg string) {
	s.on(c, model.BlockChange, msg)
}

// BlockRemoveObject blocks all operations that remove
// machines, services, units or relations.
func (s BlockHelper) BlockRemoveObject(c *gc.C, msg string) {
	s.on(c, model.BlockRemove, msg)
}

func (s BlockHelper) Close() {
	s.client.Close()
	s.apiState.Close()
}

// BlockDestroyModel blocks destroy-model.
func (s BlockHelper) BlockDestroyModel(c *gc.C, msg string) {
	s.on(c, model.BlockDestroy, msg)
}

// AssertBlocked checks if given error is
// related to switched block.
func (s BlockHelper) AssertBlocked(c *gc.C, err error, msg string) {
	c.Assert(params.IsCodeOperationBlocked(err), jc.IsTrue, gc.Commentf("error: %#v", err))
	c.Assert(errors.Cause(err), gc.DeepEquals, &params.Error{
		Message: msg,
		Code:    "operation is blocked",
	})
}
