// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelmanager

import (
	"github.com/juju/errors"

	"github.com/DavinZhang/juju/apiserver/common"
	"github.com/DavinZhang/juju/core/network"
	"github.com/DavinZhang/juju/environs/space"
	"github.com/DavinZhang/juju/state"
)

type spaceStateShim struct {
	common.ModelManagerBackend
}

func (s spaceStateShim) AllSpaces() ([]space.Space, error) {
	spaces, err := s.ModelManagerBackend.AllSpaces()
	if err != nil {
		return nil, errors.Trace(err)
	}

	results := make([]space.Space, len(spaces))
	for i, space := range spaces {
		results[i] = space
	}
	return results, nil
}

func (s spaceStateShim) AddSpace(name string, providerId network.Id, subnetIds []string, public bool) (space.Space, error) {
	result, err := s.ModelManagerBackend.AddSpace(name, providerId, subnetIds, public)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return result, nil
}

func (s spaceStateShim) ConstraintsBySpaceName(name string) ([]space.Constraints, error) {
	constraints, err := s.ModelManagerBackend.ConstraintsBySpaceName(name)
	if err != nil {
		return nil, errors.Trace(err)
	}

	results := make([]space.Constraints, len(constraints))
	for i, constraint := range constraints {
		results[i] = constraint
	}
	return results, nil
}

type statePoolShim struct {
	*state.StatePool
}

func (s statePoolShim) Get(uuid string) (State, error) {
	st, err := s.StatePool.Get(uuid)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return stateShim{
		PooledState: st,
	}, nil
}

type stateShim struct {
	*state.PooledState
}

func (s stateShim) Model() (Model, error) {
	model, err := s.PooledState.Model()
	if err != nil {
		return nil, errors.Trace(err)
	}
	return modelShim{
		Model: model,
	}, nil
}

type modelShim struct {
	*state.Model
}

func (s modelShim) IsControllerModel() bool {
	return s.Model.IsControllerModel()
}
