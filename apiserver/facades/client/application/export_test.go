// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"github.com/DavinZhang/juju/core/assumes"
	"github.com/DavinZhang/juju/environs"
	"github.com/DavinZhang/juju/state"
	"github.com/DavinZhang/juju/state/stateenvirons"
)

var (
	ParseSettingsCompatible = parseSettingsCompatible
	GetStorageState         = getStorageState
)

func GetState(st *state.State) Backend {
	return stateShim{st}
}

func GetModel(m *state.Model) Model {
	return modelShim{m}
}

func SetModelType(api *APIv13, modelType state.ModelType) {
	api.modelType = modelType
}

func MockSupportedFeatures(fs assumes.FeatureSet) {
	supportedFeaturesGetter = func(stateenvirons.Model, environs.NewEnvironFunc) (assumes.FeatureSet, error) {
		return fs, nil
	}
}

func ResetSupportedFeaturesGetter() {
	supportedFeaturesGetter = stateenvirons.SupportedFeatures
}
