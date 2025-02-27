// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common

import (
	"github.com/juju/names/v4"

	apiservererrors "github.com/DavinZhang/juju/apiserver/errors"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/core/instance"
	"github.com/DavinZhang/juju/state"
)

// InstanceIdGetter implements a common InstanceId method for use by
// various facades.
type InstanceIdGetter struct {
	st         state.EntityFinder
	getCanRead GetAuthFunc
}

// NewInstanceIdGetter returns a new InstanceIdGetter. The GetAuthFunc
// will be used on each invocation of InstanceId to determine current
// permissions.
func NewInstanceIdGetter(st state.EntityFinder, getCanRead GetAuthFunc) *InstanceIdGetter {
	return &InstanceIdGetter{
		st:         st,
		getCanRead: getCanRead,
	}
}

func (ig *InstanceIdGetter) getInstanceId(tag names.Tag) (instance.Id, error) {
	entity0, err := ig.st.FindEntity(tag)
	if err != nil {
		return "", err
	}
	entity, ok := entity0.(state.InstanceIdGetter)
	if !ok {
		return "", apiservererrors.NotSupportedError(tag, "instance id")
	}
	return entity.InstanceId()
}

// InstanceId returns the provider specific instance id for each given
// machine or an CodeNotProvisioned error, if not set.
func (ig *InstanceIdGetter) InstanceId(args params.Entities) (params.StringResults, error) {
	result := params.StringResults{
		Results: make([]params.StringResult, len(args.Entities)),
	}
	canRead, err := ig.getCanRead()
	if err != nil {
		return result, err
	}
	for i, entity := range args.Entities {
		tag, err := names.ParseTag(entity.Tag)
		if err != nil {
			result.Results[i].Error = apiservererrors.ServerError(apiservererrors.ErrPerm)
			continue
		}
		err = apiservererrors.ErrPerm
		if canRead(tag) {
			var instanceId instance.Id
			instanceId, err = ig.getInstanceId(tag)
			if err == nil {
				result.Results[i].Result = string(instanceId)
			}
		}
		result.Results[i].Error = apiservererrors.ServerError(err)
	}
	return result, nil
}
