// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package sender

import (
	"time"

	"github.com/juju/errors"
	"github.com/juju/loggo"
	"github.com/juju/names/v4"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/dependency"

	"github.com/DavinZhang/juju/agent"
	"github.com/DavinZhang/juju/api/base"
	"github.com/DavinZhang/juju/api/metricsadder"
	jworker "github.com/DavinZhang/juju/worker"
	"github.com/DavinZhang/juju/worker/metrics/spool"
	"github.com/DavinZhang/juju/worker/uniter"
)

var (
	logger               = loggo.GetLogger("juju.worker.metrics.sender")
	newMetricAdderClient = func(apiCaller base.APICaller) metricsadder.MetricsAdderClient {
		return metricsadder.NewClient(apiCaller)
	}
	period = time.Minute * 5
)

// ManifoldConfig defines configuration of a metric sender manifold.
type ManifoldConfig struct {
	AgentName       string
	APICallerName   string
	MetricSpoolName string
}

// Manifold creates a metric sender manifold.
func Manifold(config ManifoldConfig) dependency.Manifold {
	return dependency.Manifold{
		Inputs: []string{
			config.AgentName,
			config.APICallerName,
			config.MetricSpoolName,
		},
		Start: func(context dependency.Context) (worker.Worker, error) {
			var apicaller base.APICaller
			var factory spool.MetricFactory
			err := context.Get(config.APICallerName, &apicaller)
			if err != nil {
				return nil, errors.Trace(err)
			}
			err = context.Get(config.MetricSpoolName, &factory)
			if err != nil {
				return nil, errors.Trace(err)
			}
			var agent agent.Agent
			if err := context.Get(config.AgentName, &agent); err != nil {
				return nil, err
			}
			agentConfig := agent.CurrentConfig()
			tag := agentConfig.Tag()
			unitTag, ok := tag.(names.UnitTag)
			if !ok {
				return nil, errors.Errorf("expected a unit tag, got %v", tag)
			}
			paths := uniter.NewWorkerPaths(agentConfig.DataDir(), unitTag, "metrics-send", nil)

			client := newMetricAdderClient(apicaller)

			s, err := newSender(client, factory, paths.State.BaseDir, unitTag.String())
			if err != nil {
				return nil, errors.Trace(err)
			}
			return spool.NewPeriodicWorker(s.Do, period, jworker.NewTimer, s.stop), nil
		},
	}
}
