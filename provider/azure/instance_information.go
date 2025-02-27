// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package azure

import (
	"github.com/juju/errors"

	"github.com/DavinZhang/juju/core/constraints"
	"github.com/DavinZhang/juju/environs"
	"github.com/DavinZhang/juju/environs/context"
	"github.com/DavinZhang/juju/environs/instances"
)

var _ environs.InstanceTypesFetcher = (*azureEnviron)(nil)

// InstanceTypes implements InstanceTypesFetcher
func (env *azureEnviron) InstanceTypes(ctx context.ProviderCallContext, c constraints.Value) (instances.InstanceTypesWithCostMetadata, error) {
	types, err := env.getInstanceTypes(ctx)
	if err != nil {
		return instances.InstanceTypesWithCostMetadata{}, errors.Trace(err)
	}
	result := make([]instances.InstanceType, len(types))
	i := 0
	for _, iType := range types {
		result[i] = iType
		i++
	}
	result, err = instances.MatchingInstanceTypes(result, "", c)
	if err != nil {
		return instances.InstanceTypesWithCostMetadata{}, errors.Trace(err)
	}

	return instances.InstanceTypesWithCostMetadata{
		InstanceTypes: result,
		CostUnit:      "",
		CostCurrency:  "USD"}, nil
}
