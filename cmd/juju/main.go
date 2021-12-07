// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/juju/cmd/v3"
	"github.com/juju/loggo"

	"github.com/DavinZhang/juju/cmd/juju/commands"
	components "github.com/DavinZhang/juju/component/all"
	_ "github.com/DavinZhang/juju/provider/all" // Import the providers.
)

var log = loggo.GetLogger("juju.cmd.juju")

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func init() {
	if err := components.RegisterForClient(); err != nil {
		log.Criticalf("unable to register client components: %v", err)
		os.Exit(1)
	}
}

func main() {
	_, err := loggo.ReplaceDefaultWriter(cmd.NewWarningWriter(os.Stderr))
	if err != nil {
		panic(err)
	}
	os.Exit(commands.Main(os.Args))
}
