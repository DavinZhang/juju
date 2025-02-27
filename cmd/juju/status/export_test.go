// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package status

import (
	"github.com/juju/cmd/v3"

	"github.com/DavinZhang/juju/cmd/juju/storage"
	"github.com/DavinZhang/juju/cmd/modelcmd"
)

func NewTestStatusHistoryCommand(api HistoryAPI) cmd.Command {
	return &statusHistoryCommand{api: api}
}

func NewTestStatusCommand(statusapi statusAPI, storageapi storage.StorageListAPI, clock Clock) cmd.Command {
	return modelcmd.Wrap(
		&statusCommand{statusAPI: statusapi, storageAPI: storageapi, clock: clock})
}
