// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package progress_test

import (
	"testing"

	gc "gopkg.in/check.v1"
)

//go:generate go run github.com/golang/mock/mockgen -package mocks -destination ./mocks/term_mock.go github.com/DavinZhang/juju/cmd/output/progress Terminal
//go:generate go run github.com/golang/mock/mockgen -package mocks -destination ./mocks/clock_mock.go github.com/juju/clock Clock

func Test(t *testing.T) { gc.TestingT(t) }
