// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package applicationscaler

import (
	"github.com/juju/errors"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/dependency"

	"github.com/DavinZhang/juju/api/base"
	"github.com/DavinZhang/juju/cmd/jujud/agent/engine"
)

// ManifoldConfig holds dependencies and configuration for an
// applicationscaler worker.
type ManifoldConfig struct {
	APICallerName string
	NewFacade     func(base.APICaller) (Facade, error)
	NewWorker     func(Config) (worker.Worker, error)
}

// start is a method on ManifoldConfig because that feels a bit cleaner
// than closing over config in Manifold.
func (config ManifoldConfig) start(apiCaller base.APICaller) (worker.Worker, error) {
	facade, err := config.NewFacade(apiCaller)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return config.NewWorker(Config{
		Facade: facade,
	})
}

// Manifold returns a dependency.Manifold that runs an applicationscaler worker.
func Manifold(config ManifoldConfig) dependency.Manifold {
	return engine.APIManifold(
		engine.APIManifoldConfig{config.APICallerName},
		config.start,
	)
}
