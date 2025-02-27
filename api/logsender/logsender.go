// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package logsender implements the API for storing log
// messages on the API server.
package logsender

import (
	"io"
	"net/url"

	"github.com/juju/errors"

	"github.com/DavinZhang/juju/api/base"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/version"
)

// LogWriter is the interface that allows sending log
// messages to the server for storage.
type LogWriter interface {
	// WriteLog writes the given log record.
	WriteLog(*params.LogRecord) error

	io.Closer
}

// API provides access to the LogSender API.
type API struct {
	connector base.StreamConnector
}

// NewAPI creates a new client-side logsender API.
func NewAPI(connector base.StreamConnector) *API {
	return &API{connector: connector}
}

// LogWriter returns a new log writer interface value
// which must be closed when finished with.
func (api *API) LogWriter() (LogWriter, error) {
	attrs := make(url.Values)
	// TODO(wallyworld) - remove in juju 4
	attrs.Set("jujuclientversion", version.Current.String())
	// Version 1 does ping/pong handling.
	attrs.Set("version", "1")
	conn, err := api.connector.ConnectStream("/logsink", attrs)
	if err != nil {
		return nil, errors.Annotatef(err, "cannot connect to /logsink")
	}
	logWriter := writer{conn}
	go logWriter.readLoop()
	return logWriter, nil
}

type writer struct {
	conn base.Stream
}

// readLoop is necessary for the client to process websocket control messages.
// Close() is safe to call concurrently.
func (w writer) readLoop() {
	for {
		if _, _, err := w.conn.NextReader(); err != nil {
			w.conn.Close()
			break
		}
	}
}

func (w writer) WriteLog(m *params.LogRecord) error {
	// Note: due to the fire-and-forget nature of the
	// logsink API, it is possible that when the
	// connection dies, any logs that were "in-flight"
	// will not be recorded on the server side.
	if err := w.conn.WriteJSON(m); err != nil {
		return errors.Annotatef(err, "cannot send log message")
	}
	return nil
}

func (w writer) Close() error {
	return w.conn.Close()
}
