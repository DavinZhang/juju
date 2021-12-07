// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasunitprovisioner

import (
	"github.com/juju/errors"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/dependency"

	"github.com/DavinZhang/juju/api/base"
	"github.com/DavinZhang/juju/caas"
)

// Logger represents the methods used by the worker to log details.
type Logger interface {
	Debugf(string, ...interface{})
	Warningf(string, ...interface{})
	Errorf(string, ...interface{})
	Tracef(string, ...interface{})
}

// ManifoldConfig defines a CAAS unit provisioner's dependencies.
type ManifoldConfig struct {
	APICallerName string
	BrokerName    string

	NewClient func(base.APICaller) Client
	NewWorker func(Config) (worker.Worker, error)
	Logger    Logger
}

// Validate is called by start to check for bad configuration.
func (config ManifoldConfig) Validate() error {
	if config.APICallerName == "" {
		return errors.NotValidf("empty APICallerName")
	}
	if config.BrokerName == "" {
		return errors.NotValidf("empty BrokerName")
	}
	if config.NewClient == nil {
		return errors.NotValidf("nil NewClient")
	}
	if config.NewWorker == nil {
		return errors.NotValidf("nil NewWorker")
	}
	if config.Logger == nil {
		return errors.NotValidf("nil Logger")
	}
	return nil
}

func (config ManifoldConfig) start(context dependency.Context) (worker.Worker, error) {
	if err := config.Validate(); err != nil {
		return nil, errors.Trace(err)
	}

	var apiCaller base.APICaller
	if err := context.Get(config.APICallerName, &apiCaller); err != nil {
		return nil, errors.Trace(err)
	}

	var broker caas.Broker
	if err := context.Get(config.BrokerName, &broker); err != nil {
		return nil, errors.Trace(err)
	}

	client := config.NewClient(apiCaller)
	w, err := config.NewWorker(Config{
		ApplicationGetter:  client,
		ApplicationUpdater: client,

		ServiceBroker:   broker,
		ContainerBroker: broker,

		ProvisioningInfoGetter:   client,
		ProvisioningStatusSetter: client,
		LifeGetter:               client,
		UnitUpdater:              client,
		CharmGetter:              client,

		Logger: config.Logger,
	})
	if err != nil {
		return nil, errors.Trace(err)
	}
	return w, nil
}

// Manifold creates a manifold that runs a CAAS unit provisioner. See the
// ManifoldConfig type for discussion about how this can/should evolve.
func Manifold(config ManifoldConfig) dependency.Manifold {
	return dependency.Manifold{
		Inputs: []string{
			config.APICallerName,
			config.BrokerName,
		},
		Start: config.start,
	}
}
