// Copyright 2015 Canonical Ltd.
// Copyright 2015 Cloudbase Solutions
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juju/utils/featureflag"
	"golang.org/x/sys/windows/svc"

	"github.com/juju/juju/cmd/service"
	"github.com/juju/juju/juju/names"
	"github.com/juju/juju/juju/osenv"
)

// FLAGSFROMENVIRONMENT can control whether we read featureflags from the
// environment or from the registry. This is only needed because we build the
// jujud binary in uniter tests and we cannot mock the registry out easily.
// Once uniter tests are fixed this should be removed.
var FLAGSFROMENVIRONMENT string

func init() {
	if FLAGSFROMENVIRONMENT == "true" {
		featureflag.SetFlagsFromEnvironment(osenv.JujuFeatureFlagEnvKey)
	} else {
		featureflag.SetFlagsFromRegistry(osenv.JujuRegistryKey, osenv.JujuFeatureFlagEnvKey)
	}
}

func main() {
    MainWrapper(os.Args)
}
