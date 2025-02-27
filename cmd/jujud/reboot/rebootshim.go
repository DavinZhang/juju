// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package reboot

import (
	"github.com/juju/os/v2/series"

	"github.com/DavinZhang/juju/agent"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/container"
	"github.com/DavinZhang/juju/container/factory"
	"github.com/DavinZhang/juju/core/instance"
	"github.com/DavinZhang/juju/service"
	"github.com/DavinZhang/juju/service/common"
)

// rebootWaiterShim wraps the functions required by RebootWaiter
// to facilitate mock testing.
type rebootWaiterShim struct {
}

// HostSeries returns the series of the current host.
func (r rebootWaiterShim) HostSeries() (string, error) {
	return series.HostSeries()
}

// ListServices returns a list of names of services running
// on the current host.
func (r rebootWaiterShim) ListServices() ([]string, error) {
	return service.ListServices()
}

// NewService returns a new juju service object.
func (r rebootWaiterShim) NewService(name string, conf common.Conf, series string) (Service, error) {
	return service.NewService(name, conf, series)
}

// NewContainerManager return an object implementing Manager.
func (r rebootWaiterShim) NewContainerManager(containerType instance.ContainerType, conf container.ManagerConfig) (Manager, error) {
	return factory.NewContainerManager(containerType, conf)
}

// ScheduleAction schedules the reboot action based on the
// current operating system.
func (r rebootWaiterShim) ScheduleAction(action params.RebootAction, after int) error {
	return scheduleAction(action, after)
}

// agentConfigShim wraps the method required by a Model in
// the RebootWaiter.
type agentConfigShim struct {
	aCfg agent.Config
}

// Model return an object implementing Model.
func (a *agentConfigShim) Model() Model {
	return a.aCfg.Model()
}
