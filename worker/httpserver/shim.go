// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package httpserver

import (
	"github.com/juju/worker/v3"

	"github.com/DavinZhang/juju/controller"
	"github.com/DavinZhang/juju/state"
)

// NewWorkerShim calls through to NewWorker, and exists only
// to adapt to the signature of ManifoldConfig.NewWorker.
func NewWorkerShim(config Config) (worker.Worker, error) {
	return NewWorker(config)
}

// GetControllerConfig gets the controller config from a *State - it
// exists so we can test the manifold without a StateSuite.
func GetControllerConfig(st *state.State) (controller.Config, error) {
	return st.ControllerConfig()
}
