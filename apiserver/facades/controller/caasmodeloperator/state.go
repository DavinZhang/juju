// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasmodeloperator

import (
	"github.com/juju/names/v4"

	"github.com/DavinZhang/juju/apiserver/common"
	"github.com/DavinZhang/juju/controller"
	"github.com/DavinZhang/juju/environs/config"
	"github.com/DavinZhang/juju/state"
)

// CAASModelOperatorState provides the subset of model state required by the
// model operator provisioner.
type CAASModelOperatorState interface {
	FindEntity(tag names.Tag) (state.Entity, error)
	Model() (Model, error)
	ModelUUID() string
}

// CAASModelOperatorState provides the subset of controller state required by the
// model operator provisioner.
type CAASControllerState interface {
	common.APIAddressAccessor
	ControllerConfig() (controller.Config, error)
}

type Model interface {
	ModelConfig() (*config.Config, error)
}

type stateShim struct {
	*state.State
}

func (s stateShim) Model() (Model, error) {
	model, err := s.State.Model()
	if err != nil {
		return nil, err
	}
	return model.CAASModel()
}
