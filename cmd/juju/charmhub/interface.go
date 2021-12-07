// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charmhub

import (
	"context"
	"net/url"

	apicharmhub "github.com/DavinZhang/juju/api/charmhub"
	"github.com/DavinZhang/juju/charmhub"
	"github.com/DavinZhang/juju/charmhub/transport"
)

// Printer defines an interface for printing out values.
type Printer interface {
	Print() error
}

// Log describes a log format function to output to.
type Log = func(format string, params ...interface{})

// InfoCommandAPI describes API methods required to execute the info command.
type InfoCommandAPI interface {
	Info(string, ...apicharmhub.InfoOption) (apicharmhub.InfoResponse, error)
	Close() error
}

// FindCommandAPI describes API methods required to execute the find command.
type FindCommandAPI interface {
	Find(string, ...apicharmhub.FindOption) ([]apicharmhub.FindResponse, error)
	Close() error
}

// DownloadCommandAPI describes API methods required to execute the download
// command.
type DownloadCommandAPI interface {
	Info(context.Context, string, ...charmhub.InfoOption) (transport.InfoResponse, error)
	Refresh(context.Context, charmhub.RefreshConfig) ([]transport.RefreshResponse, error)
	Download(context.Context, *url.URL, string, ...charmhub.DownloadOption) error
}

type ModelConfigGetter interface {
	ModelGet() (map[string]interface{}, error)
}
