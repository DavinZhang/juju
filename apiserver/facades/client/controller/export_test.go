// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package controller

import (
	"github.com/DavinZhang/juju/apiserver/facade"
	"github.com/DavinZhang/juju/core/migration"
	"github.com/DavinZhang/juju/state"
)

type patcher interface {
	PatchValue(destination, source interface{})
}

func SetPrecheckResult(p patcher, err error) {
	p.PatchValue(&runMigrationPrechecks, func(*state.State, *state.State, *migration.TargetInfo, facade.Presence) error {
		return err
	})
}

func NewControllerAPIForTest(backend Backend) *ControllerAPI {
	return &ControllerAPI{state: backend}
}
