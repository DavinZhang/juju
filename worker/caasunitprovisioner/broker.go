// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasunitprovisioner

import (
	"github.com/DavinZhang/juju/caas"
	"github.com/DavinZhang/juju/core/application"
	"github.com/DavinZhang/juju/core/watcher"
	"github.com/juju/names/v4"
)

type ContainerBroker interface {
	Provider() caas.ContainerEnvironProvider
	WatchOperator(string) (watcher.NotifyWatcher, error)
	Operator(string) (*caas.Operator, error)

	WatchUnits(appName string, mode caas.DeploymentMode) (watcher.NotifyWatcher, error)
	Units(appName string, mode caas.DeploymentMode) ([]caas.Unit, error)
	AnnotateUnit(appName string, mode caas.DeploymentMode, podName string, unit names.UnitTag) error
}

type ServiceBroker interface {
	Provider() caas.ContainerEnvironProvider
	EnsureService(appName string, statusCallback caas.StatusCallbackFunc, params *caas.ServiceParams, numUnits int, config application.ConfigAttributes) error
	DeleteService(appName string) error
	UnexposeService(appName string) error

	GetService(appName string, mode caas.DeploymentMode, includeClusterIP bool) (*caas.Service, error)
	WatchService(appName string, mode caas.DeploymentMode) (watcher.NotifyWatcher, error)
}
