// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package kvm_test

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/DavinZhang/juju/core/network"

	"github.com/juju/loggo"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v2/arch"
	"github.com/juju/version/v2"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/cloudconfig/instancecfg"
	"github.com/DavinZhang/juju/container"
	"github.com/DavinZhang/juju/container/kvm"
	"github.com/DavinZhang/juju/core/constraints"
	"github.com/DavinZhang/juju/core/status"
	"github.com/DavinZhang/juju/environs/config"
	"github.com/DavinZhang/juju/environs/imagemetadata"
	"github.com/DavinZhang/juju/environs/instances"
	jujutesting "github.com/DavinZhang/juju/juju/testing"
	coretesting "github.com/DavinZhang/juju/testing"
	"github.com/DavinZhang/juju/tools"
	jujuversion "github.com/DavinZhang/juju/version"
)

type LiveSuite struct {
	coretesting.BaseSuite
	ContainerDir string
	RemovedDir   string
}

var _ = gc.Suite(&LiveSuite{})

func (s *LiveSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)
	// Skip if not linux
	if runtime.GOOS != "linux" {
		c.Skip("not running linux")
	}
	// Skip if virsh is not installed.
	if _, err := exec.LookPath("virsh"); err != nil {
		c.Skip("virsh not found")
	}
	// Skip if not running as root.
	if os.Getuid() != 0 {
		c.Skip("not running as root")
	}
	s.ContainerDir = c.MkDir()
	s.PatchValue(&container.ContainerDir, s.ContainerDir)
	s.RemovedDir = c.MkDir()
	s.PatchValue(&container.RemovedContainerDir, s.RemovedDir)
	loggo.GetLogger("juju.container").SetLogLevel(loggo.TRACE)
}

func (s *LiveSuite) newManager(c *gc.C, name string) container.Manager {
	manager, err := kvm.NewContainerManager(
		container.ManagerConfig{
			container.ConfigModelUUID: coretesting.ModelTag.Id(),
			container.ConfigLogDir:    c.MkDir(),
		})
	c.Assert(err, jc.ErrorIsNil)
	return manager
}

func assertNumberOfContainers(c *gc.C, manager container.Manager, count int) {
	containers, err := manager.ListContainers()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(containers, gc.HasLen, count)
}

func (s *LiveSuite) TestNoInitialContainers(c *gc.C) {
	manager := s.newManager(c, "test")
	assertNumberOfContainers(c, manager, 0)
}

func shutdownMachines(manager container.Manager) func(*gc.C) {
	return func(c *gc.C) {
		instances, err := manager.ListContainers()
		c.Assert(err, jc.ErrorIsNil)
		for _, instance := range instances {
			err := manager.DestroyContainer(instance.Id())
			c.Check(err, jc.ErrorIsNil)
		}
	}
}

func createContainer(c *gc.C, manager container.Manager, machineId string) instances.Instance {
	machineNonce := "fake-nonce"
	apiInfo := jujutesting.FakeAPIInfo(machineId)
	instanceConfig, err := instancecfg.NewInstanceConfig(coretesting.ControllerTag, machineId, machineNonce, imagemetadata.ReleasedStream, "quantal", apiInfo)
	c.Assert(err, jc.ErrorIsNil)

	nics := network.InterfaceInfos{{
		InterfaceName: "eth0",
		InterfaceType: network.EthernetDevice,
		ConfigType:    network.ConfigDHCP,
	}}
	net := container.BridgeNetworkConfig(0, nics)

	err = instanceConfig.SetTools(tools.List{
		&tools.Tools{
			Version: version.MustParseBinary("2.3.4-foo-bar"),
			URL:     "http://tools.testing.invalid/2.3.4-foo-bar.tgz",
		},
	})
	c.Assert(err, jc.ErrorIsNil)
	environConfig := dummyConfig(c)
	err = instancecfg.FinishInstanceConfig(instanceConfig, environConfig)
	c.Assert(err, jc.ErrorIsNil)
	callback := func(settableStatus status.Status, info string, data map[string]interface{}) error { return nil }
	inst, hardware, err := manager.CreateContainer(instanceConfig, constraints.Value{}, "precise", net, nil, callback)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(hardware, gc.NotNil)
	expected := fmt.Sprintf("arch=%s cores=1 mem=512M root-disk=8192M", arch.HostArch())
	c.Assert(hardware.String(), gc.Equals, expected)
	return inst
}

func (s *LiveSuite) TestShutdownMachines(c *gc.C) {
	manager := s.newManager(c, "test")
	createContainer(c, manager, "1/kvm/0")
	createContainer(c, manager, "1/kvm/1")
	assertNumberOfContainers(c, manager, 2)

	shutdownMachines(manager)(c)
	assertNumberOfContainers(c, manager, 0)
}

func (s *LiveSuite) TestManagerIsolation(c *gc.C) {
	firstManager := s.newManager(c, "first")
	s.AddCleanup(shutdownMachines(firstManager))

	createContainer(c, firstManager, "1/kvm/0")
	createContainer(c, firstManager, "1/kvm/1")

	secondManager := s.newManager(c, "second")
	s.AddCleanup(shutdownMachines(secondManager))

	createContainer(c, secondManager, "1/kvm/0")

	assertNumberOfContainers(c, firstManager, 2)
	assertNumberOfContainers(c, secondManager, 1)
}

func dummyConfig(c *gc.C) *config.Config {
	testConfig, err := config.New(config.UseDefaults, coretesting.FakeConfig())
	c.Assert(err, jc.ErrorIsNil)
	testConfig, err = testConfig.Apply(map[string]interface{}{
		"type":          "dummy",
		"controller":    false,
		"agent-version": jujuversion.Current.String(),
	})
	c.Assert(err, jc.ErrorIsNil)
	return testConfig
}
