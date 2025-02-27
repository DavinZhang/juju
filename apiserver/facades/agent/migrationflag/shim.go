// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package migrationflag

import (
	"github.com/juju/errors"

	"github.com/DavinZhang/juju/apiserver/facade"
	"github.com/DavinZhang/juju/core/migration"
	"github.com/DavinZhang/juju/state"
)

// NewFacade wraps New to express the supplied *state.State as a Backend.
func NewFacade(st *state.State, resources facade.Resources, auth facade.Authorizer) (*Facade, error) {
	facade, err := New(&backend{st}, resources, auth)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return facade, nil
}

// backend implements Backend by wrapping a *state.State.
type backend struct {
	st *state.State
}

// ModelUUID is part of the Backend interface.
func (shim *backend) ModelUUID() string {
	return shim.st.ModelUUID()
}

// WatchMigrationPhase is part of the Backend interface.
func (shim *backend) WatchMigrationPhase() state.NotifyWatcher {
	return shim.st.WatchMigrationStatus()
}

// MigrationPhase is part of the Backend interface.
func (shim *backend) MigrationPhase() (migration.Phase, error) {
	mig, err := shim.st.LatestMigration()
	if errors.IsNotFound(err) {
		return migration.NONE, nil
	} else if err != nil {
		return migration.UNKNOWN, errors.Trace(err)
	}
	phase, err := mig.Phase()
	if err != nil {
		return migration.UNKNOWN, errors.Trace(err)
	}
	return phase, nil
}
