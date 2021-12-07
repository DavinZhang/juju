// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package applicationscaler

import (
	"github.com/DavinZhang/juju/api/applicationscaler"
	"github.com/DavinZhang/juju/api/base"
	"github.com/DavinZhang/juju/api/watcher"
)

// NewFacade creates a Facade from a base.APICaller.
// It's a sensible value for ManifoldConfig.NewFacade.
func NewFacade(apiCaller base.APICaller) (Facade, error) {
	return applicationscaler.NewAPI(
		apiCaller,
		watcher.NewStringsWatcher,
	), nil
}
