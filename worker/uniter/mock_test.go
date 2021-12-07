// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package uniter_test

import (
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/worker/uniter/operation"
	"github.com/DavinZhang/juju/worker/uniter/remotestate"
	"github.com/DavinZhang/juju/worker/uniter/resolver"
	"github.com/DavinZhang/juju/worker/uniter/storage"
	"github.com/juju/names/v4"
)

type dummyStorageAccessor struct {
	storage.StorageAccessor
}

func (*dummyStorageAccessor) UnitStorageAttachments(_ names.UnitTag) ([]params.StorageAttachmentId, error) {
	return nil, nil
}

type nopResolver struct{}

func (nopResolver) NextOp(resolver.LocalState, remotestate.Snapshot, operation.Factory) (operation.Operation, error) {
	return nil, resolver.ErrNoOperation
}
