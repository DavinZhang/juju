// Copyright 2019 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upgradedatabase

import (
	"time"

	"github.com/juju/errors"
	"github.com/juju/version/v2"

	"github.com/DavinZhang/juju/core/status"
	"github.com/DavinZhang/juju/mongo"
	"github.com/DavinZhang/juju/state"
)

// Logger represents the methods required to emit log messages.
type Logger interface {
	Debugf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Errorf(message string, args ...interface{})
}

// Clock provides an interface for dealing with clocks.
type Clock interface {
	// After waits for the duration to elapse and then sends the
	// current time on the returned channel.
	After(time.Duration) <-chan time.Time
}

// UpgradeInfo describes methods for working with the state representation of
// an upgrade-in-progress.
type UpgradeInfo interface {
	// Status returns the current status of the upgrade.
	Status() state.UpgradeStatus

	// SetStatus sets the current status of the upgrade.
	SetStatus(status state.UpgradeStatus) error

	// Watch returns a watcher that notifies of changes to the UpgradeInfo.
	Watch() state.NotifyWatcher

	// Refresh refreshes the UpgradeInfo from state.
	Refresh() error
}

// Pool describes methods required by the upgradeDB worker,
// supplied by a state pool.
type Pool interface {
	// IsPrimary returns true if the Mongo primary is
	// running on the machine with the input ID.
	IsPrimary(string) (bool, error)

	// SetStatus updates the status of the machine with the input ID.
	SetStatus(string, status.Status, string) error

	// EnsureUpgradeInfo ensures that a document exists in the "upgradeInfo"
	// collection for coordinating the current upgrade.
	EnsureUpgradeInfo(string, version.Number, version.Number) (UpgradeInfo, error)

	// Close closes the state pool.
	Close() error
}

type pool struct {
	*state.StatePool
}

// IsPrimary (Pool) returns true if the Mongo primary is
// running on the controller with the input ID.
func (p *pool) IsPrimary(controllerId string) (bool, error) {
	st := p.SystemState()

	// For IAAS models, controllers are machines.
	// For CAAS models, until we support HA, there is only one Mongo
	// and it is the primary.
	hasMachine, err := p.hasMachine()
	if err != nil {
		return false, errors.Trace(err)
	}
	// TODO(CAAS) - bug 1849030 support HA
	if !hasMachine {
		return true, nil
	}

	machine, err := st.Machine(controllerId)
	if err != nil {
		return false, errors.Trace(err)
	}
	isPrimary, err := mongo.IsMaster(st.MongoSession(), machine)
	return isPrimary, errors.Trace(err)
}

func (p *pool) hasMachine() (bool, error) {
	model, err := p.SystemState().Model()
	if err != nil {
		return false, errors.Trace(err)
	}
	return model.Type() == state.ModelTypeIAAS, nil
}

// SetStatus (Pool) updates the status of the machine with the input ID.
func (p *pool) SetStatus(controllerId string, sts status.Status, msg string) error {
	hasMachine, err := p.hasMachine()
	if err != nil {
		return errors.Trace(err)
	}
	if !hasMachine {
		// TODO(CAAS) - bug 1849030 support HA
		// Nothing we can do for now because we do not have any machine for CAAS controller.
		return nil
	}
	machine, err := p.SystemState().Machine(controllerId)
	if err != nil {
		return errors.Trace(err)
	}

	now := time.Now()
	return errors.Trace(machine.SetStatus(status.StatusInfo{
		Status:  sts,
		Message: msg,
		Since:   &now,
	}))
}

// EnsureUpgradeInfo (Pool) ensures that a document exists in the "upgradeInfo"
// collection for coordinating the current upgrade.
func (p *pool) EnsureUpgradeInfo(controllerId string, fromVersion, toVersion version.Number) (UpgradeInfo, error) {
	info, err := p.SystemState().EnsureUpgradeInfo(controllerId, fromVersion, toVersion)
	return info, errors.Trace(err)
}
