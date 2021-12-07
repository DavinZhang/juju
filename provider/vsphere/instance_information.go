// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package vsphere

import (
	"github.com/juju/errors"

	"github.com/DavinZhang/juju/core/constraints"
	"github.com/DavinZhang/juju/environs"
	"github.com/DavinZhang/juju/environs/context"
	"github.com/DavinZhang/juju/environs/instances"
)

var _ environs.InstanceTypesFetcher = (*environ)(nil)

// InstanceTypes implements InstanceTypesFetcher
func (env *environ) InstanceTypes(ctx context.ProviderCallContext, c constraints.Value) (instances.InstanceTypesWithCostMetadata, error) {
	result := instances.InstanceTypesWithCostMetadata{}
	return result, errors.NotSupportedf("InstanceTypes")
}
