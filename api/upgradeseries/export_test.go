// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upgradeseries

import (
	"github.com/DavinZhang/juju/api/base"
	"github.com/DavinZhang/juju/api/common"
	"github.com/juju/names/v4"
)

func NewStateFromCaller(caller base.FacadeCaller, authTag names.Tag) *Client {
	return &Client{
		facade:           caller,
		authTag:          authTag,
		UpgradeSeriesAPI: common.NewUpgradeSeriesAPI(caller, authTag),
	}
}
