// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package authenticationworker

import (
	"github.com/juju/errors"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/dependency"

	"github.com/DavinZhang/juju/agent"
	"github.com/DavinZhang/juju/api/base"
	"github.com/DavinZhang/juju/api/keyupdater"
	"github.com/DavinZhang/juju/cmd/jujud/agent/engine"
)

// ManifoldConfig defines the names of the manifolds on which a Manifold will depend.
type ManifoldConfig engine.AgentAPIManifoldConfig

// Manifold returns a dependency manifold that runs a authenticationworker worker,
// using the resource names defined in the supplied config.
func Manifold(config ManifoldConfig) dependency.Manifold {
	typedConfig := engine.AgentAPIManifoldConfig(config)

	return engine.AgentAPIManifold(typedConfig, newWorker)
}

func newWorker(a agent.Agent, apiCaller base.APICaller) (worker.Worker, error) {
	w, err := NewWorker(keyupdater.NewState(apiCaller), a.CurrentConfig())
	if err != nil {
		return nil, errors.Annotate(err, "cannot start ssh auth-keys updater worker")
	}
	return w, nil
}
