// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasfirewallersidecar

import (
	"strings"

	"github.com/juju/errors"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/catacomb"

	"github.com/DavinZhang/juju/caas"
	"github.com/DavinZhang/juju/core/watcher"
)

type applicationWorker struct {
	catacomb       catacomb.Catacomb
	controllerUUID string
	modelUUID      string
	appName        string

	firewallerAPI CAASFirewallerAPI

	broker         CAASBroker
	portMutator    PortMutator
	serviceUpdater ServiceUpdater

	appWatcher   watcher.NotifyWatcher
	portsWatcher watcher.StringsWatcher

	lifeGetter LifeGetter

	initial           bool
	previouslyExposed bool

	logger Logger
}

func newApplicationWorker(
	controllerUUID string,
	modelUUID string,
	appName string,
	firewallerAPI CAASFirewallerAPI,
	broker CAASBroker,
	lifeGetter LifeGetter,
	logger Logger,
) (worker.Worker, error) {
	w := &applicationWorker{
		controllerUUID: controllerUUID,
		modelUUID:      modelUUID,
		appName:        appName,
		firewallerAPI:  firewallerAPI,
		broker:         broker,
		lifeGetter:     lifeGetter,
		initial:        true,
		logger:         logger,
	}
	if err := catacomb.Invoke(catacomb.Plan{
		Site: &w.catacomb,
		Work: w.loop,
	}); err != nil {
		return nil, errors.Trace(err)
	}
	return w, nil
}

// Kill is part of the worker.Worker interface.
func (w *applicationWorker) Kill() {
	w.catacomb.Kill(nil)
}

// Wait is part of the worker.Worker interface.
func (w *applicationWorker) Wait() error {
	return w.catacomb.Wait()
}

func (w *applicationWorker) setUp() (err error) {
	w.appWatcher, err = w.firewallerAPI.WatchApplication(w.appName)
	if err != nil {
		return errors.Trace(err)
	}
	if err := w.catacomb.Add(w.appWatcher); err != nil {
		return errors.Trace(err)
	}

	w.portsWatcher, err = w.firewallerAPI.WatchOpenedPorts()
	if err != nil {
		return errors.Trace(err)
	}
	if err := w.catacomb.Add(w.portsWatcher); err != nil {
		return errors.Trace(err)
	}

	// TODO(sidecar): support deployment other than statefulset
	app := w.broker.Application(w.appName, caas.DeploymentStateful)
	w.portMutator = app
	w.serviceUpdater = app

	// TODO(sidecar):
	/*
		if ports, err := w.firewallerAPI.OpenedPorts(w.appName); err != nil {
			return errors.Annotatef(err, "failed to get initial openned ports for application")
		}
		w.currentPorts = newPortRanges(ports)
	*/
	return nil
}

func (w *applicationWorker) loop() (err error) {
	defer func() {
		// If the application has been deleted, we can return nil.
		if errors.IsNotFound(err) {
			w.logger.Debugf("sidecar caas firewaller application %v has been removed", w.appName)
			err = nil
		}
	}()

	if err = w.setUp(); err != nil {
		return errors.Trace(err)
	}

	for {
		select {
		case <-w.catacomb.Dying():
			return w.catacomb.ErrDying()
		case _, ok := <-w.appWatcher.Changes():
			if !ok {
				return errors.New("application watcher closed")
			}
			// We know this is a v2 charm at this point, because this child
			// worker is only ever started for v2 charms.
			if err := w.onApplicationChanged(); err != nil {
				if strings.Contains(err.Error(), "unexpected EOF") {
					return nil
				}
				return errors.Trace(err)
			}
		case _, ok := <-w.portsWatcher.Changes():
			if !ok {
				return errors.New("application watcher closed")
			}
			// TODO(sidecar): implement portWatcher to return application names having port changes,
			/*
				if !sets.NewString(changes...).Contains(w.appName){continue}
			*/
			if err := w.onPortChanged(); err != nil {
				return errors.Trace(err)
			}
		}
	}
}

func (w *applicationWorker) onPortChanged() (err error) {
	// TODO(sidecar):
	/*
		ports, err := w.firewallerAPI.OpenedPorts(w.appName)
		if err != nil {
			return err
		}
		changedPortRanges := newPortRanges(ports)
		if w.current.equal(changedPortRanges){
			logger.Tracef("no port changes for app %q", w.appName)
			return nil
		}
		w.currentPorts = changedPortRanges
		return w.portMutator.UpdatePorts(w.currentPorts.toServicePorts(), false)
	*/
	return nil
}

func (w *applicationWorker) onApplicationChanged() (err error) {
	defer func() {
		// Not found could be because the app got removed or there's
		// no container service created yet as the app is still being set up.
		if errors.IsNotFound(err) {
			// Perhaps the app got removed while we were processing.
			if _, err2 := w.lifeGetter.Life(w.appName); err2 != nil {
				err = err2
				return
			}
			// Ignore not found error because the ip could be not ready yet at this stage.
			w.logger.Warningf("processing change for application %q, %v", w.appName, err)
			err = nil
		}
	}()

	exposed, err := w.firewallerAPI.IsExposed(w.appName)
	if err != nil {
		return errors.Trace(err)
	}
	if !w.initial && exposed == w.previouslyExposed {
		return nil
	}

	w.initial = false
	w.previouslyExposed = exposed
	if exposed {
		return errors.Trace(exposeService(w.serviceUpdater))
	}
	return errors.Trace(unExposeService(w.serviceUpdater))
}

func exposeService(app ServiceUpdater) error {
	// TODO(sidecar): implement expose once it's modelled.
	// app.UpdateService()
	return nil
}

func unExposeService(app ServiceUpdater) error {
	// TODO(sidecar): implement un-expose once it's modelled.
	// app.UpdateService()
	return nil
}
