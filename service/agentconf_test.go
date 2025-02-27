// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// The test cases in this file do not pertain to a specific command.

package service_test

import (
	"bytes"
	"os"
	"path"

	"github.com/juju/errors"
	"github.com/juju/names/v4"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v2/arch"
	"github.com/juju/version/v2"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/agent"
	agenttools "github.com/DavinZhang/juju/agent/tools"
	"github.com/DavinZhang/juju/service"
	"github.com/DavinZhang/juju/service/common"
	svctesting "github.com/DavinZhang/juju/service/common/testing"
	"github.com/DavinZhang/juju/testing"
	coretest "github.com/DavinZhang/juju/tools"
	jujuversion "github.com/DavinZhang/juju/version"
)

type agentConfSuite struct {
	testing.BaseSuite

	agentConf           agent.Config
	dataDir             string
	machineName         string
	unitNames           []string
	systemdDir          string
	systemdMultiUserDir string
	systemdDataDir      string // service.SystemdDataDir
	manager             service.SystemdServiceManager

	services    []*svctesting.FakeService
	serviceData *svctesting.FakeServiceData
}

func (s *agentConfSuite) SetUpSuite(c *gc.C) {
	s.BaseSuite.SetUpSuite(c)
}

func (s *agentConfSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)

	s.dataDir = c.MkDir()
	s.systemdDir = path.Join(s.dataDir, "etc", "systemd", "system")
	s.systemdMultiUserDir = path.Join(s.systemdDir, "multi-user.target.wants")
	c.Assert(os.MkdirAll(s.systemdMultiUserDir, os.ModeDir|os.ModePerm), jc.ErrorIsNil)
	s.systemdDataDir = path.Join(s.dataDir, "lib", "systemd", "system")

	s.machineName = "machine-0"
	s.unitNames = []string{"unit-ubuntu-0", "unit-mysql-0"}

	s.manager = service.NewServiceManager(
		func() bool { return true },
		s.newService,
	)

	s.assertSetupAgentsForTest(c)
	s.setUpAgentConf(c)
	s.setUpServices(c)
	s.services[0].ResetCalls()
	s.setupTools(c, "trusty")
}

func (s *agentConfSuite) TearDownTest(c *gc.C) {
	s.serviceData = nil
	s.services = nil
	s.BaseSuite.TearDownTest(c)
}

var _ = gc.Suite(&agentConfSuite{})

func (s *agentConfSuite) setUpAgentConf(c *gc.C) {
	// Required for CopyAgentBinaries to evaluate the version of the agent.
	configParams := agent.AgentConfigParams{
		Paths:             agent.Paths{DataDir: s.dataDir},
		Tag:               names.NewMachineTag("0"),
		UpgradedToVersion: jujuversion.Current,
		APIAddresses:      []string{"localhost:17070"},
		CACert:            testing.CACert,
		Password:          "fake",
		Controller:        testing.ControllerTag,
		Model:             testing.ModelTag,
	}

	agentConf, err := agent.NewAgentConfig(configParams)
	c.Assert(err, jc.ErrorIsNil)

	err = agentConf.Write()
	c.Assert(err, jc.ErrorIsNil)

	s.agentConf = agentConf
}

func (s *agentConfSuite) setUpServices(c *gc.C) {
	s.addService(c, "jujud-"+s.machineName)
	s.PatchValue(&service.ListServices, s.listServices)
}

func (s *agentConfSuite) addService(c *gc.C, name string) {
	svc, err := s.newService(name, common.Conf{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(svc.Install(), jc.ErrorIsNil)
	c.Assert(svc.Start(), jc.ErrorIsNil)
}

func (s *agentConfSuite) listServices() ([]string, error) {
	return s.serviceData.InstalledNames(), nil
}

func (s *agentConfSuite) newService(name string, _ common.Conf) (service.Service, error) {
	for _, svc := range s.services {
		if svc.Name() == name {
			return svc, nil
		}
	}
	if s.serviceData == nil {
		s.serviceData = svctesting.NewFakeServiceData()
	}
	svc := &svctesting.FakeService{
		FakeServiceData: s.serviceData,
		Service: common.Service{
			Name: name,
			Conf: common.Conf{},
		},
		DataDir: s.dataDir,
	}
	s.services = append(s.services, svc)
	return svc, nil
}

func (s *agentConfSuite) setupTools(c *gc.C, series string) {
	files := []*testing.TarFile{
		testing.NewTarFile("jujud", 0755, "jujuc executable"),
	}
	data, checksum := testing.TarGz(files...)
	testTools := &coretest.Tools{
		URL: "http://foo/bar1",
		Version: version.Binary{
			Number:  jujuversion.Current,
			Arch:    arch.HostArch(),
			Release: series,
		},
		Size:   int64(len(data)),
		SHA256: checksum,
	}
	err := agenttools.UnpackTools(s.dataDir, testTools, bytes.NewReader(data))
	c.Assert(err, jc.ErrorIsNil)
}

func (s *agentConfSuite) assertSetupAgentsForTest(c *gc.C) {
	agentsDir := path.Join(s.dataDir, "agents")
	err := os.MkdirAll(path.Join(agentsDir, s.machineName), os.ModeDir|os.ModePerm)
	c.Assert(err, jc.ErrorIsNil)
	for _, unit := range s.unitNames {
		err = os.Mkdir(path.Join(agentsDir, unit), os.ModeDir|os.ModePerm)
		c.Assert(err, jc.ErrorIsNil)
	}
}

func (s *agentConfSuite) TestFindAgents(c *gc.C) {
	machineAgent, unitAgents, errAgents, err := s.manager.FindAgents(s.dataDir)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(machineAgent, gc.Equals, s.machineName)
	c.Assert(unitAgents, jc.SameContents, s.unitNames)
	c.Assert(errAgents, gc.HasLen, 0)
}

func (s *agentConfSuite) TestFindAgentsUnexpectedTagType(c *gc.C) {
	unexpectedAgent := names.NewApplicationTag("failme").String()
	unexpectedAgentDir := path.Join(s.dataDir, "agents", unexpectedAgent)
	err := os.MkdirAll(unexpectedAgentDir, os.ModeDir|os.ModePerm)
	c.Assert(err, jc.ErrorIsNil)

	machineAgent, unitAgents, unexpectedAgents, err := s.manager.FindAgents(s.dataDir)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(machineAgent, gc.Equals, s.machineName)
	c.Assert(unitAgents, jc.SameContents, s.unitNames)
	c.Assert(unexpectedAgents, gc.DeepEquals, []string{unexpectedAgent})
}

func (s *agentConfSuite) TestCreateAgentConfDesc(c *gc.C) {
	conf, err := s.manager.CreateAgentConf("machine-2", s.dataDir)
	c.Assert(err, jc.ErrorIsNil)
	// Spot check Conf
	c.Assert(conf.Desc, gc.Equals, "juju agent for machine-2")
}

func (s *agentConfSuite) TestCreateAgentConfLogPath(c *gc.C) {
	conf, err := s.manager.CreateAgentConf("machine-2", s.dataDir)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(conf.Logfile, gc.Equals, "/var/log/juju/machine-2.log")
}

func (s *agentConfSuite) TestCreateAgentConfFailAgentKind(c *gc.C) {
	_, err := s.manager.CreateAgentConf("application-fail", s.dataDir)
	c.Assert(err, gc.ErrorMatches, `agent "application-fail" is neither a machine nor a unit`)
}

func (s *agentConfSuite) agentUnitNames() []string {
	unitAgents := make([]string, len(s.unitNames))
	for i, name := range s.unitNames {
		unitAgents[i] = "jujud-" + name
	}
	return unitAgents
}

func (s *agentConfSuite) TestWriteSystemdAgent(c *gc.C) {
	err := s.manager.WriteSystemdAgent(
		s.machineName, s.systemdDataDir, s.systemdMultiUserDir)

	c.Assert(err, jc.ErrorIsNil)
	s.assertServicesCalls(c, "WriteService", 1)
}

func (s *agentConfSuite) TestWriteSystemdAgentSystemdNotRunning(c *gc.C) {
	s.manager = service.NewServiceManager(
		func() bool { return false },
		s.newService,
	)

	err := s.manager.WriteSystemdAgent(
		s.machineName, s.systemdDataDir, s.systemdMultiUserDir)

	c.Assert(err, jc.ErrorIsNil)
	s.assertServicesCalls(c, "WriteService", 1)
}

func (s *agentConfSuite) TestWriteSystemdAgentDBusErrManualLink(c *gc.C) {
	// nil errors are for calls to RemoveOldService.
	err := errors.New("no such method 'LinkUnitFiles'")
	s.services[0].SetErrors(nil, err)

	err = s.manager.WriteSystemdAgent(
		s.machineName, s.systemdDataDir, s.systemdMultiUserDir)

	c.Assert(err, jc.ErrorIsNil)

	// This exhibits the same characteristics as for Systemd not running (above).
	s.assertServicesCalls(c, "RemoveOldService", 1)
	s.assertServicesCalls(c, "WriteService", 1)
}

func (s *agentConfSuite) TestWriteSystemdAgentWriteServiceFail(c *gc.C) {
	// Return an error for the machine agent.
	s.services[0].SetErrors(nil, errors.New("fail me"))

	err := s.manager.WriteSystemdAgent(
		s.machineName, s.systemdDataDir, s.systemdMultiUserDir)

	c.Assert(err, gc.ErrorMatches, "fail me")
	s.assertServicesCalls(c, "RemoveOldService", 1)
	s.assertServicesCalls(c, "WriteService", 1)
}

func (s *agentConfSuite) assertToolsCopySymlink(c *gc.C, series string) {
	// Check tools changes.
	ver := version.Binary{
		Number:  jujuversion.Current,
		Arch:    arch.HostArch(),
		Release: series,
	}
	jujuTools, err := agenttools.ReadTools(s.dataDir, ver)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(jujuTools.Version, gc.DeepEquals, ver)

	link := path.Join(s.dataDir, "tools", s.machineName)
	linkResult, err := os.Readlink(link)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(linkResult, gc.Equals, path.Join(s.dataDir, "tools", ver.String()))
}

func (s *agentConfSuite) assertServicesCalls(c *gc.C, svc string, expectedCnt int) {
	// Call list shared by the services
	calls := s.services[0].Calls()
	serviceCount := 0
	for _, call := range calls {
		if call.FuncName == svc {
			serviceCount += 1
		}
	}
	c.Assert(serviceCount, gc.Equals, expectedCnt)
}
