// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package raft_test

import (
	stdtesting "testing"

	"github.com/DavinZhang/juju/testing"
)

//go:generate go run github.com/golang/mock/mockgen -package raft -destination raft_mock_test.go github.com/DavinZhang/juju/worker/raft Raft,ApplierMetrics
//go:generate go run github.com/golang/mock/mockgen -package raft -destination raftlease_mock_test.go github.com/DavinZhang/juju/core/raftlease NotifyTarget,FSMResponse
//go:generate go run github.com/golang/mock/mockgen -package raft -destination raft_future_mock_test.go github.com/hashicorp/raft ApplyFuture,ConfigurationFuture

func TestPackage(t *stdtesting.T) {
	testing.MgoTestPackage(t)
}
