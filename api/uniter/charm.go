// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package uniter

import (
	"fmt"

	"github.com/juju/charm/v9"
	"github.com/juju/errors"

	"github.com/DavinZhang/juju/apiserver/params"
)

// This module implements a subset of the interface provided by
// state.Charm, as needed by the uniter API.

// Charm represents the state of a charm in the model.
type Charm struct {
	st   *State
	curl *charm.URL
}

// String returns the charm URL as a string.
func (c *Charm) String() string {
	return c.curl.String()
}

// URL returns the URL that identifies the charm.
func (c *Charm) URL() *charm.URL {
	return c.curl
}

// ArchiveSha256 returns the SHA256 digest of the charm archive
// (bundle) bytes.
//
// NOTE: This differs from state.Charm.BundleSha256() by returning an
// error as well, because it needs to make an API call. It's also
// renamed to avoid confusion with juju deployment bundles.
//
// TODO(dimitern): 2013-09-06 bug 1221834
// Cache the result after getting it once for the same charm URL,
// because it's immutable.
func (c *Charm) ArchiveSha256() (string, error) {
	var results params.StringResults
	args := params.CharmURLs{
		URLs: []params.CharmURL{{URL: c.curl.String()}},
	}
	err := c.st.facade.FacadeCall("CharmArchiveSha256", args, &results)
	if err != nil {
		return "", err
	}
	if len(results.Results) != 1 {
		return "", fmt.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return "", result.Error
	}
	return result.Result, nil
}

// LXDProfileRequired returns true if this charm requires an
// lxd profile to be applied.
func (c *Charm) LXDProfileRequired() (bool, error) {
	var results params.BoolResults
	args := params.CharmURLs{
		URLs: []params.CharmURL{{URL: c.curl.String()}},
	}
	err := c.st.facade.FacadeCall("LXDProfileRequired", args, &results)
	if err != nil {
		return false, err
	}
	if len(results.Results) != 1 {
		return false, errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return false, result.Error
	}
	return result.Result, nil
}
