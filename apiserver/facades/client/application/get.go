// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"github.com/juju/charm/v9"
	"github.com/juju/schema"
	"gopkg.in/juju/environschema.v1"

	apiservererrors "github.com/DavinZhang/juju/apiserver/errors"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/caas"
	"github.com/DavinZhang/juju/core/application"
	"github.com/DavinZhang/juju/core/constraints"
	"github.com/DavinZhang/juju/core/model"
)

// Get returns the charm configuration for an application.
func (api *APIBase) Get(args params.ApplicationGet) (params.ApplicationGetResults, error) {
	return api.getConfig(args, describe)
}

// Get returns the charm configuration for an application.
func (api *APIBase) getConfig(
	args params.ApplicationGet,
	describe func(settings charm.Settings, config *charm.Config) map[string]interface{},
) (params.ApplicationGetResults, error) {
	if err := api.checkCanRead(); err != nil {
		return params.ApplicationGetResults{}, err
	}

	app, err := api.backend.Application(args.ApplicationName)
	if err != nil {
		return params.ApplicationGetResults{}, err
	}

	// We need a guard on the API server-side for direct API callers such as
	// python-libjuju. Always default to the master branch.
	if args.BranchName == "" {
		args.BranchName = model.GenerationMaster
	}
	settings, err := app.CharmConfig(args.BranchName)
	if err != nil {
		return params.ApplicationGetResults{}, err
	}

	ch, _, err := app.Charm()
	if err != nil {
		return params.ApplicationGetResults{}, err
	}
	configInfo := describe(settings, ch.Config())
	appConfig, err := app.ApplicationConfig()
	if err != nil {
		return params.ApplicationGetResults{}, err
	}

	providerSchema, providerDefaults, err := applicationConfigSchema(api.modelType)
	if err != nil {
		return params.ApplicationGetResults{}, err
	}
	appConfigInfo := describeAppConfig(appConfig, providerSchema, caas.ConfigDefaults(providerDefaults))
	var cons constraints.Value
	if app.IsPrincipal() {
		cons, err = app.Constraints()
		if err != nil {
			return params.ApplicationGetResults{}, err
		}
	}
	endpoints, err := app.EndpointBindings()
	if err != nil {
		return params.ApplicationGetResults{}, err
	}

	allSpaceInfosLookup, err := api.backend.AllSpaceInfos()
	if err != nil {
		return params.ApplicationGetResults{}, apiservererrors.ServerError(err)
	}

	bindingMap, err := endpoints.MapWithSpaceNames(allSpaceInfosLookup)
	if err != nil {
		return params.ApplicationGetResults{}, err
	}

	return params.ApplicationGetResults{
		Application:       args.ApplicationName,
		Charm:             ch.Meta().Name,
		CharmConfig:       configInfo,
		ApplicationConfig: appConfigInfo,
		Constraints:       cons,
		Series:            app.Series(),
		Channel:           string(app.Channel()),
		EndpointBindings:  bindingMap,
	}, nil
}

func describeAppConfig(
	appConfig application.ConfigAttributes,
	schema environschema.Fields,
	defaults schema.Defaults,
) map[string]interface{} {
	results := make(map[string]interface{})
	for name, field := range schema {
		defaultValue := defaults[name]
		info := map[string]interface{}{
			"description": field.Description,
			"type":        field.Type,
			"source":      "unset",
		}
		set := false
		if value := appConfig[name]; value != nil && defaultValue != value {
			set = true
			info["value"] = value
			info["source"] = "user"
		}
		if defaultValue != nil {
			info["default"] = defaultValue
			if !set {
				info["value"] = defaultValue
				info["source"] = "default"
			}
		}
		results[name] = info
	}
	return results
}

func describe(settings charm.Settings, config *charm.Config) map[string]interface{} {
	results := make(map[string]interface{})
	for name, option := range config.Options {
		info := map[string]interface{}{
			"description": option.Description,
			"type":        option.Type,
			"source":      "unset",
		}
		set := false
		if value := settings[name]; value != nil && option.Default != value {
			set = true
			info["value"] = value
			info["source"] = "user"
		}
		if option.Default != nil {
			info["default"] = option.Default
			if !set {
				info["value"] = option.Default
				info["source"] = "default"
			}
		}
		results[name] = info
	}
	return results
}
