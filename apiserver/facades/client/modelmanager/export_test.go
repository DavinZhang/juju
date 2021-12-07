// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelmanager

import (
	"github.com/juju/names/v4"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/core/assumes"
	"github.com/DavinZhang/juju/environs"
	"github.com/DavinZhang/juju/state/stateenvirons"
)

func AuthCheck(c *gc.C, mm *ModelManagerAPI, user names.UserTag) bool {
	mm.authCheck(user)
	return mm.isAdmin
}

func MockSupportedFeatures(fs assumes.FeatureSet) {
	supportedFeaturesGetter = func(stateenvirons.Model, environs.NewEnvironFunc) (assumes.FeatureSet, error) {
		return fs, nil
	}
}

func ResetSupportedFeaturesGetter() {
	supportedFeaturesGetter = stateenvirons.SupportedFeatures
}
