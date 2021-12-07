// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lifeflag

import (
	"github.com/juju/errors"
	"github.com/juju/worker/v3"

	"github.com/DavinZhang/juju/api/base"
	"github.com/DavinZhang/juju/api/lifeflag"
	"github.com/DavinZhang/juju/api/watcher"
)

func NewWorker(config Config) (worker.Worker, error) {
	worker, err := New(config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return worker, nil
}

func NewFacade(apiCaller base.APICaller) (Facade, error) {
	facade := lifeflag.NewFacade(apiCaller, watcher.NewNotifyWatcher)
	return facade, nil
}
