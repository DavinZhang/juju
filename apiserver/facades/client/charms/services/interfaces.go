// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package services

import (
	"io"

	"github.com/juju/charm/v9"

	"github.com/DavinZhang/juju/controller"
	"github.com/DavinZhang/juju/environs/config"
	"github.com/DavinZhang/juju/state"
)

// StateBackend describes an API for accessing/mutating information in state.
type StateBackend interface {
	ControllerConfig() (controller.Config, error)
	UpdateUploadedCharm(info state.CharmInfo) (UploadedCharm, error)
	PrepareCharmUpload(curl *charm.URL) (UploadedCharm, error)
	ModelUUID() string
}

// UploadedCharm represents a charm whose upload status can be queried.
type UploadedCharm interface {
	IsUploaded() bool
}

// ModelBackend describes an API for accessing model-specific details.
type ModelBackend interface {
	Config() (*config.Config, error)
}

// Storage describes an API for storing and deleting blobs.
type Storage interface {
	Put(path string, r io.Reader, length int64) error
	Remove(path string) error
}
