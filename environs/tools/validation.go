// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package tools

import (
	"fmt"

	"github.com/juju/version/v2"

	"github.com/DavinZhang/juju/environs/simplestreams"
	jujuversion "github.com/DavinZhang/juju/version"
)

// ToolsMetadataLookupParams is used to query metadata for matching tools.
type ToolsMetadataLookupParams struct {
	simplestreams.MetadataLookupParams
	Version string
	Major   int
	Minor   int
}

// ValidateToolsMetadata attempts to load tools metadata for the specified cloud attributes and returns
// any tools versions found, or an error if the metadata could not be loaded.
func ValidateToolsMetadata(ss SimplestreamsFetcher, params *ToolsMetadataLookupParams) ([]string, *simplestreams.ResolveInfo, error) {
	if len(params.Sources) == 0 {
		return nil, nil, fmt.Errorf("required parameter sources not specified")
	}
	if params.Version == "" && params.Major == 0 {
		params.Version = jujuversion.Current.String()
	}
	var toolsConstraint *ToolsConstraint
	if params.Version == "" {
		toolsConstraint = NewGeneralToolsConstraint(params.Major, params.Minor, simplestreams.LookupParams{
			CloudSpec: simplestreams.CloudSpec{
				Region:   params.Region,
				Endpoint: params.Endpoint,
			},
			Stream:   params.Stream,
			Releases: []string{params.Release},
			Arches:   params.Architectures,
		})
	} else {
		versNum, err := version.Parse(params.Version)
		if err != nil {
			return nil, nil, err
		}
		toolsConstraint = NewVersionedToolsConstraint(versNum, simplestreams.LookupParams{
			CloudSpec: simplestreams.CloudSpec{
				Region:   params.Region,
				Endpoint: params.Endpoint,
			},
			Stream:   params.Stream,
			Releases: []string{params.Release},
			Arches:   params.Architectures,
		})
	}
	matchingTools, resolveInfo, err := Fetch(ss, params.Sources, toolsConstraint)
	if err != nil {
		return nil, resolveInfo, err
	}
	if len(matchingTools) == 0 {
		return nil, resolveInfo, fmt.Errorf("no matching agent binaries found for constraint %+v", toolsConstraint)
	}
	versions := make([]string, len(matchingTools))
	for i, tm := range matchingTools {
		vers := version.Binary{
			Number:  version.MustParse(tm.Version),
			Release: tm.Release,
			Arch:    tm.Arch,
		}
		versions[i] = vers.String()
	}
	return versions, resolveInfo, nil
}
