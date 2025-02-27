// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasunitprovisioner

import (
	apicaasunitprovisioner "github.com/DavinZhang/juju/api/caasunitprovisioner"
	charmscommon "github.com/DavinZhang/juju/api/common/charms"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/caas"
	"github.com/DavinZhang/juju/core/application"
	"github.com/DavinZhang/juju/core/life"
	"github.com/DavinZhang/juju/core/status"
	"github.com/DavinZhang/juju/core/watcher"
)

// Client provides an interface for interacting with the
// CAASUnitProvisioner API. Subsets of this should be passed
// to the CAASUnitProvisioner worker.
type Client interface {
	ApplicationGetter
	ApplicationUpdater
	ProvisioningInfoGetter
	LifeGetter
	UnitUpdater
	ProvisioningStatusSetter
	CharmGetter
}

// ApplicationGetter provides an interface for
// watching for the lifecycle state changes
// (including addition) of applications in the
// model, and fetching their details.
type ApplicationGetter interface {
	WatchApplications() (watcher.StringsWatcher, error)
	WatchApplication(appName string) (watcher.NotifyWatcher, error)
	ApplicationConfig(string) (application.ConfigAttributes, error)
	DeploymentMode(string) (caas.DeploymentMode, error)
	WatchApplicationScale(string) (watcher.NotifyWatcher, error)
	ApplicationScale(string) (int, error)
}

// ApplicationUpdater provides an interface for updating
// Juju applications from changes in the cloud.
type ApplicationUpdater interface {
	UpdateApplicationService(arg params.UpdateApplicationServiceArg) error
	ClearApplicationResources(appName string) error
}

// ProvisioningInfoGetter provides an interface for
// watching and getting the pod spec and other info
// needed to provision an application.
type ProvisioningInfoGetter interface {
	ProvisioningInfo(appName string) (*apicaasunitprovisioner.ProvisioningInfo, error)
	WatchPodSpec(appName string) (watcher.NotifyWatcher, error)
}

// LifeGetter provides an interface for getting the
// lifecycle state value for an application or unit.
type LifeGetter interface {
	Life(string) (life.Value, error)
}

// UnitUpdater provides an interface for updating
// Juju units from changes in the cloud.
type UnitUpdater interface {
	UpdateUnits(arg params.UpdateApplicationUnits) (*params.UpdateApplicationUnitsInfo, error)
}

// ProvisioningStatusSetter provides an interface for
// setting status information.
type ProvisioningStatusSetter interface {
	// SetOperatorStatus sets the status for the application operator.
	SetOperatorStatus(appName string, status status.Status, message string, data map[string]interface{}) error
}

type CharmGetter interface {
	ApplicationCharmInfo(appName string) (*charmscommon.CharmInfo, error)
}
