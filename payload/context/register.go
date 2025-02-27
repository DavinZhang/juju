// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package context

import (
	"github.com/juju/charm/v9"
	"github.com/juju/cmd/v3"
	"github.com/juju/errors"

	jujucmd "github.com/DavinZhang/juju/cmd"
	"github.com/DavinZhang/juju/payload"
)

// RegisterCmdName is the name of the payload register command.
const RegisterCmdName = "payload-register"

// NewRegisterCmd returns a new RegisterCmd that wraps the given context.
func NewRegisterCmd(ctx HookContext) (*RegisterCmd, error) {
	return &RegisterCmd{hookContextFunc: componentHookContext(ctx)}, nil
}

// RegisterCmd is a command that registers a payload with juju.
type RegisterCmd struct {
	cmd.CommandBase

	hookContextFunc func() (Component, error)
	typ             string
	class           string
	id              string
	labels          []string
}

// TODO(ericsnow) Change "tags" to "labels" in the help text?

// Info implements cmd.Command.
func (c RegisterCmd) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:    RegisterCmdName,
		Args:    "<type> <class> <id> [tags...]",
		Purpose: "register a charm payload with juju",
		Doc: `
"payload-register" is used while a hook is running to let Juju know that a
payload has been started. The information used to start the payload must be
provided when "register" is run.

The payload class must correspond to one of the payloads defined in
the charm's metadata.yaml.

		`,
	})
}

// Init implements cmd.Command.
func (c *RegisterCmd) Init(args []string) error {
	if len(args) < 3 {
		return errors.Errorf("missing required arguments")
	}
	c.typ = args[0]
	c.class = args[1]
	c.id = args[2]
	c.labels = args[3:]
	return nil
}

// Run implements cmd.Command.
func (c *RegisterCmd) Run(ctx *cmd.Context) error {
	if err := c.validate(ctx); err != nil {
		return errors.Trace(err)
	}
	pl := payload.Payload{
		PayloadClass: charm.PayloadClass{
			Name: c.class,
			Type: c.typ,
		},
		ID:     c.id,
		Status: payload.StateRunning,
		Labels: c.labels,
		Unit:   "a-application/0", // TODO(ericsnow) eliminate this!
	}
	hctx, err := c.hookContextFunc()
	if err != nil {
		return errors.Trace(err)
	}
	if err := hctx.Track(pl); err != nil {
		return errors.Trace(err)
	}

	// We flush to state immediately so that status reflects the
	// payload correctly.
	if err := hctx.Flush(); err != nil {
		return errors.Trace(err)
	}

	// TODO(ericsnow) Print out the full ID.

	return nil
}

func (c *RegisterCmd) validate(ctx *cmd.Context) error {
	meta, err := readMetadata(ctx)
	if err != nil {
		return errors.Trace(err)
	}

	found := false
	for _, class := range meta.PayloadClasses {
		if c.class == class.Name {
			if c.typ != class.Type {
				return errors.Errorf("incorrect type %q for payload %q, expected %q", c.typ, class.Name, class.Type)
			}
			found = true
		}
	}
	if !found {
		return errors.Errorf("payload %q not found in metadata.yaml", c.class)
	}
	return nil
}
