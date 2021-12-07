// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package imagemanager

import (
	"github.com/juju/names/v4"

	"github.com/DavinZhang/juju/state"
	"github.com/DavinZhang/juju/state/imagestorage"
)

type stateInterface interface {
	ImageStorage() imagestorage.Storage
	ControllerTag() names.ControllerTag
}

type stateShim struct {
	*state.State
}

func (s stateShim) ImageStorage() imagestorage.Storage {
	return s.State.ImageStorage()
}
