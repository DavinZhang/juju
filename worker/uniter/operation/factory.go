// Copyright 2014-2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package operation

import (
	corecharm "github.com/juju/charm/v9"
	"github.com/juju/errors"
	"github.com/juju/names/v4"

	"github.com/DavinZhang/juju/api/uniter"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/worker/common/charmrunner"
	"github.com/DavinZhang/juju/worker/uniter/charm"
	"github.com/DavinZhang/juju/worker/uniter/hook"
	"github.com/DavinZhang/juju/worker/uniter/remotestate"
	"github.com/DavinZhang/juju/worker/uniter/runner"
)

// FactoryParams holds all the necessary parameters for a new operation factory.
type FactoryParams struct {
	Deployer       charm.Deployer
	RunnerFactory  runner.Factory
	Callbacks      Callbacks
	State          *uniter.State
	Abort          <-chan struct{}
	MetricSpoolDir string
	Logger         Logger
}

// NewFactory returns a Factory that creates Operations backed by the supplied
// parameters.
func NewFactory(params FactoryParams) Factory {
	return &factory{
		config: params,
	}
}

type factory struct {
	config FactoryParams
}

// newDeploy is the common code for creating arbitrary deploy operations.
func (f *factory) newDeploy(kind Kind, charmURL *corecharm.URL, revert, resolved bool) (Operation, error) {
	if f.config.Deployer == nil {
		return nil, errors.New("deployer required")
	}
	if charmURL == nil {
		return nil, errors.New("charm url required")
	} else if kind != Install && kind != Upgrade {
		return nil, errors.Errorf("unknown deploy kind: %s", kind)
	}
	var op Operation = &deploy{
		kind:      kind,
		charmURL:  charmURL,
		revert:    revert,
		resolved:  resolved,
		callbacks: f.config.Callbacks,
		deployer:  f.config.Deployer,
		abort:     f.config.Abort,
	}
	return op, nil
}

// NewInstall is part of the Factory interface.
func (f *factory) NewInstall(charmURL *corecharm.URL) (Operation, error) {
	return f.newDeploy(Install, charmURL, false, false)
}

// NewUpgrade is part of the Factory interface.
func (f *factory) NewUpgrade(charmURL *corecharm.URL) (Operation, error) {
	return f.newDeploy(Upgrade, charmURL, false, false)
}

// NewRemoteInit is part of the Factory interface.
func (f *factory) NewRemoteInit(runningStatus remotestate.ContainerRunningStatus) (Operation, error) {
	return &remoteInit{
		callbacks:     f.config.Callbacks,
		abort:         f.config.Abort,
		runningStatus: runningStatus,
	}, nil
}

func (f *factory) NewSkipRemoteInit(retry bool) (Operation, error) {
	return &skipRemoteInit{retry}, nil
}

func (f *factory) NewNoOpFinishUpgradeSeries() (Operation, error) {
	return &noOpFinishUpgradeSeries{&skipOperation{}}, nil
}

// NewRevertUpgrade is part of the Factory interface.
func (f *factory) NewRevertUpgrade(charmURL *corecharm.URL) (Operation, error) {
	return f.newDeploy(Upgrade, charmURL, true, false)
}

// NewResolvedUpgrade is part of the Factory interface.
func (f *factory) NewResolvedUpgrade(charmURL *corecharm.URL) (Operation, error) {
	return f.newDeploy(Upgrade, charmURL, false, true)
}

// NewRunHook is part of the Factory interface.
func (f *factory) NewRunHook(hookInfo hook.Info) (Operation, error) {
	if err := hookInfo.Validate(); err != nil {
		return nil, err
	}
	return &runHook{
		info:          hookInfo,
		callbacks:     f.config.Callbacks,
		runnerFactory: f.config.RunnerFactory,
		logger:        f.config.Logger,
	}, nil
}

// NewSkipHook is part of the Factory interface.
func (f *factory) NewSkipHook(hookInfo hook.Info) (Operation, error) {
	hookOp, err := f.NewRunHook(hookInfo)
	if err != nil {
		return nil, err
	}
	return &skipOperation{hookOp}, nil
}

// NewAction is part of the Factory interface.
func (f *factory) NewAction(actionId string) (Operation, error) {
	if !names.IsValidAction(actionId) {
		return nil, errors.Errorf("invalid action id %q", actionId)
	}

	tag := names.NewActionTag(actionId)
	action, err := f.config.State.Action(tag)
	if params.IsCodeNotFoundOrCodeUnauthorized(err) {
		return nil, charmrunner.ErrActionNotAvailable
	} else if params.IsCodeActionNotAvailable(err) {
		return nil, charmrunner.ErrActionNotAvailable
	} else if err != nil {
		return nil, errors.Trace(err)
	}

	return &runAction{
		action:        action,
		callbacks:     f.config.Callbacks,
		runnerFactory: f.config.RunnerFactory,
		logger:        f.config.Logger,
	}, nil
}

// NewFailAction is part of the factory interface.
func (f *factory) NewFailAction(actionId string) (Operation, error) {
	if !names.IsValidAction(actionId) {
		return nil, errors.Errorf("invalid action id %q", actionId)
	}
	return &failAction{
		actionId:  actionId,
		callbacks: f.config.Callbacks,
	}, nil
}

// NewCommands is part of the Factory interface.
func (f *factory) NewCommands(args CommandArgs, sendResponse CommandResponseFunc) (Operation, error) {
	if err := args.Validate(); err != nil {
		return nil, errors.Trace(err)
	} else if sendResponse == nil {
		return nil, errors.New("response sender required")
	}
	return &runCommands{
		args:          args,
		sendResponse:  sendResponse,
		callbacks:     f.config.Callbacks,
		runnerFactory: f.config.RunnerFactory,
		logger:        f.config.Logger,
	}, nil
}

// NewResignLeadership is part of the Factory interface.
func (f *factory) NewResignLeadership() (Operation, error) {
	return &resignLeadership{logger: f.config.Logger}, nil
}

// NewAcceptLeadership is part of the Factory interface.
func (f *factory) NewAcceptLeadership() (Operation, error) {
	return &acceptLeadership{}, nil
}
