// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

//go:build !windows
// +build !windows

package unit

import (
	"github.com/juju/cmd/v3"
	"github.com/juju/names/v4"
	"github.com/juju/utils/v2/voyeur"

	"github.com/DavinZhang/juju/agent"
	"github.com/DavinZhang/juju/cmd/containeragent/utils"
	"github.com/DavinZhang/juju/cmd/jujud/agent/agentconf"
	"github.com/DavinZhang/juju/worker/logsender"
)

type (
	ManifoldsConfig    = manifoldsConfig
	ContainerUnitAgent = containerUnitAgent
)

type ContainerUnitAgentTest interface {
	cmd.Command
	DataDir() string
	SetAgentConf(cfg agentconf.AgentConf)
	ChangeConfig(change agent.ConfigMutator) error
	CurrentConfig() agent.Config
	Tag() names.UnitTag
	CharmModifiedVersion() int
	GetContainerNames() []string
}

func NewForTest(
	ctx *cmd.Context,
	bufferedLogger *logsender.BufferedLogWriter,
	configChangedVal *voyeur.Value,
	fileReaderWriter utils.FileReaderWriter,
	environment utils.Environment,
) ContainerUnitAgentTest {
	return &containerUnitAgent{
		ctx:              ctx,
		AgentConf:        agentconf.NewAgentConf(""),
		bufferedLogger:   bufferedLogger,
		configChangedVal: configChangedVal,
		fileReaderWriter: fileReaderWriter,
		environment:      environment,
	}
}

func (c *containerUnitAgent) SetAgentConf(cfg agentconf.AgentConf) {
	c.AgentConf = cfg
}

func (c *containerUnitAgent) GetContainerNames() []string {
	return c.containerNames
}

func (c *containerUnitAgent) DataDir() string {
	return c.AgentConf.DataDir()
}
