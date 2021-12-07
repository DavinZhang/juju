// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package dummy

import (
	"github.com/DavinZhang/juju/storage"
	dummystorage "github.com/DavinZhang/juju/storage/provider/dummy"
)

func StorageProviders() storage.ProviderRegistry {
	return dummystorage.StorageProviders()
}
