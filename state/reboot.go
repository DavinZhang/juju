// Copyright 2014 Canonical Ltd.
// Copyright 2014 Cloudbase Solutions SRL
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/juju/mgo/v2"
	"github.com/juju/mgo/v2/bson"
	"github.com/juju/mgo/v2/txn"

	"github.com/DavinZhang/juju/core/container"
)

var _ RebootFlagSetter = (*Machine)(nil)
var _ RebootActionGetter = (*Machine)(nil)

// RebootAction defines the action a machine should
// take when a hook needs to reboot
type RebootAction string

const (
	// ShouldDoNothing instructs a machine agent that no action
	// is required on its part
	ShouldDoNothing RebootAction = "noop"
	// ShouldReboot instructs a machine to reboot
	// this happens when a hook running on a machine, requests
	// a reboot
	ShouldReboot RebootAction = "reboot"
	// ShouldShutdown instructs a machine to shut down. This usually
	// happens when running inside a container, and a hook on the parent
	// machine requests a reboot
	ShouldShutdown RebootAction = "shutdown"
)

// rebootDoc will hold the reboot flag for a machine.
type rebootDoc struct {
	DocID     string `bson:"_id"`
	Id        string `bson:"machineid"`
	ModelUUID string `bson:"model-uuid"`
}

func (m *Machine) setFlag() error {
	if m.Life() == Dead {
		return mgo.ErrNotFound
	}
	ops := []txn.Op{
		assertModelActiveOp(m.st.ModelUUID()),
		{
			C:      machinesC,
			Id:     m.doc.DocID,
			Assert: notDeadDoc,
		}, {
			C:      rebootC,
			Id:     m.doc.DocID,
			Insert: &rebootDoc{Id: m.Id()},
		},
	}
	err := m.st.db().RunTransaction(ops)
	if err == txn.ErrAborted {
		if err := checkModelActive(m.st); err != nil {
			return errors.Trace(err)
		}
		return mgo.ErrNotFound
	} else if err != nil {
		return errors.Errorf("failed to set reboot flag: %v", err)
	}
	return nil
}

func removeRebootDocOp(st *State, machineId string) txn.Op {
	op := txn.Op{
		C:      rebootC,
		Id:     st.docID(machineId),
		Remove: true,
	}
	return op
}

func (m *Machine) clearFlag() error {
	reboot, closer := m.st.db().GetCollection(rebootC)
	defer closer()

	docID := m.doc.DocID
	count, err := reboot.FindId(docID).Count()
	if err != nil {
		return errors.Trace(err)
	}
	if count == 0 {
		return nil
	}
	ops := []txn.Op{removeRebootDocOp(m.st, m.Id())}
	err = m.st.db().RunTransaction(ops)
	if err != nil {
		return errors.Errorf("failed to clear reboot flag: %v", err)
	}
	return nil
}

// SetRebootFlag sets the reboot flag of a machine to a boolean value. It will also
// do a lazy create of a reboot document if needed; i.e. If a document
// does not exist yet for this machine, it will create it.
func (m *Machine) SetRebootFlag(flag bool) error {
	if flag {
		return m.setFlag()
	}
	return m.clearFlag()
}

// GetRebootFlag returns the reboot flag for this machine.
func (m *Machine) GetRebootFlag() (bool, error) {
	rebootCol, closer := m.st.db().GetCollection(rebootC)
	defer closer()

	count, err := rebootCol.FindId(m.doc.DocID).Count()
	if err != nil {
		return false, fmt.Errorf("failed to get reboot flag: %v", err)
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (m *Machine) machinesToCareAboutRebootsFor() []string {
	var possibleIds []string
	for currentId := m.Id(); currentId != ""; {
		possibleIds = append(possibleIds, currentId)
		currentId = container.ParentId(currentId)
	}
	return possibleIds
}

// ShouldRebootOrShutdown check if the current node should reboot or shutdown
// If we are a container, and our parent needs to reboot, this should return:
// ShouldShutdown
func (m *Machine) ShouldRebootOrShutdown() (RebootAction, error) {
	rebootCol, closer := m.st.db().GetCollection(rebootC)
	defer closer()

	machines := m.machinesToCareAboutRebootsFor()

	docs := []rebootDoc{}
	sel := bson.D{{"machineid", bson.D{{"$in", machines}}}}
	if err := rebootCol.Find(sel).All(&docs); err != nil {
		return ShouldDoNothing, errors.Trace(err)
	}

	iNeedReboot := false
	for _, val := range docs {
		if val.Id != m.doc.Id {
			return ShouldShutdown, nil
		}
		iNeedReboot = true
	}
	if iNeedReboot {
		return ShouldReboot, nil
	}
	return ShouldDoNothing, nil
}

type RebootFlagSetter interface {
	SetRebootFlag(flag bool) error
}

type RebootActionGetter interface {
	ShouldRebootOrShutdown() (RebootAction, error)
}
