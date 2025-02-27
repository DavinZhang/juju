// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package commands

import (
	"github.com/DavinZhang/juju/api/keymanager"
	"github.com/DavinZhang/juju/cmd/modelcmd"
)

type SSHKeysBase struct {
	modelcmd.ModelCommandBase
	modelcmd.IAASOnlyCommand
}

// NewKeyManagerClient returns a keymanager client for the root api endpoint
// that the environment command returns.
func (c *SSHKeysBase) NewKeyManagerClient() (*keymanager.Client, error) {
	root, err := c.NewAPIRoot()
	if err != nil {
		return nil, err
	}
	return keymanager.NewClient(root), nil
}
