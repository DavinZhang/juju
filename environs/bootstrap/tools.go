// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package bootstrap

import (
	"fmt"
	"os"

	"github.com/juju/errors"
	"github.com/juju/utils/v2/arch"
	"github.com/juju/version/v2"

	"github.com/DavinZhang/juju/core/constraints"
	coreos "github.com/DavinZhang/juju/core/os"
	coreseries "github.com/DavinZhang/juju/core/series"
	"github.com/DavinZhang/juju/environs"
	envtools "github.com/DavinZhang/juju/environs/tools"
	coretools "github.com/DavinZhang/juju/tools"
	jujuversion "github.com/DavinZhang/juju/version"
)

var (
	findTools = envtools.FindTools
)

func localToolsArch() string {
	toolsArch := os.Getenv("GOARCH")
	if toolsArch == "" {
		toolsArch = arch.HostArch()
	}
	return toolsArch
}

// validateUploadAllowed returns an error if an attempt to upload tools should
// not be allowed.
func validateUploadAllowed(env environs.ConfigGetter, toolsArch, toolsSeries *string, validator constraints.Validator) error {
	// Now check that the architecture and series for which we are setting up an
	// environment matches that from which we are bootstrapping.
	hostArch := localToolsArch()
	// We can't build tools for a different architecture if one is specified.
	if toolsArch != nil && *toolsArch != hostArch {
		return fmt.Errorf("cannot use agent built for %q using a machine running on %q", *toolsArch, hostArch)
	}
	hostOS := coreos.HostOS()
	if toolsSeries != nil {
		toolsSeriesOS, err := coreseries.GetOSFromSeries(*toolsSeries)
		if err != nil {
			return errors.Trace(err)
		}
		if !toolsSeriesOS.EquivalentTo(hostOS) {
			return errors.Errorf("cannot use agent built for %q using a machine running %q", *toolsSeries, hostOS)
		}
	}
	// If no architecture is specified, ensure the target provider supports instances matching our architecture.
	if _, err := validator.Validate(constraints.Value{Arch: &hostArch}); err != nil {
		return errors.Errorf(
			"model %q of type %s does not support instances running on %q",
			env.Config().Name(), env.Config().Type(), hostArch,
		)
	}
	return nil
}

// findPackagedTools returns a list of tools for in simplestreams.
func findPackagedTools(
	env environs.BootstrapEnviron,
	ss envtools.SimplestreamsFetcher,
	vers *version.Number,
	arch, series *string,
) (coretools.List, error) {
	// Look for tools in the environment's simplestreams search paths
	// for existing tools.

	// If the user hasn't asked for a specified tools version, see if
	// one is configured in the environment.
	if vers == nil {
		if agentVersion, ok := env.Config().AgentVersion(); ok {
			vers = &agentVersion
		}
	}
	logger.Infof("looking for bootstrap agent binaries: version=%v", vers)
	toolsList, findToolsErr := findBootstrapTools(env, ss, vers, arch, series)
	logger.Infof("found %d packaged agent binaries", len(toolsList))
	if findToolsErr != nil {
		return nil, findToolsErr
	}
	return toolsList, nil
}

// locallyBuildableTools returns the list of tools that
// can be built locally.
func locallyBuildableTools() (buildable coretools.List, _ version.Number, _ error) {
	buildNumber := jujuversion.Current
	// Increment the build number so we know it's a custom build.
	buildNumber.Build++
	if !coreos.HostOS().EquivalentTo(coreos.Ubuntu) {
		return buildable, buildNumber, nil
	}
	binary := version.Binary{
		Number:  buildNumber,
		Release: "ubuntu",
		Arch:    localToolsArch(),
	}
	buildable = append(buildable, &coretools.Tools{Version: binary})
	return buildable, buildNumber, nil
}

// findBootstrapTools returns a tools.List containing only those tools with
// which it would be reasonable to launch an environment's first machine,
// given the supplied constraints. If a specific agent version is not requested,
// all tools matching the current major.minor version are chosen.
func findBootstrapTools(env environs.BootstrapEnviron, ss envtools.SimplestreamsFetcher, vers *version.Number, arch, series *string) (list coretools.List, err error) {
	// Construct a tools filter.
	cliVersion := jujuversion.Current
	var filter coretools.Filter
	if arch != nil {
		filter.Arch = *arch
	}
	if series != nil {
		filter.OSType = coreseries.DefaultOSTypeNameFromSeries(*series)
	}
	if vers != nil {
		filter.Number = *vers
	}
	streams := envtools.PreferredStreams(vers, env.Config().Development(), env.Config().AgentStream())
	return findTools(ss, env, cliVersion.Major, cliVersion.Minor, streams, filter)
}
