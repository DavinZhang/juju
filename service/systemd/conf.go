// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package systemd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/coreos/go-systemd/v22/unit"
	"github.com/juju/errors"
	"github.com/juju/utils/v2/shell"

	"github.com/DavinZhang/juju/core/paths"
	"github.com/DavinZhang/juju/service/common"
)

var limitMap = map[string]string{
	"as":         "LimitAS",
	"core":       "LimitCORE",
	"cpu":        "LimitCPU",
	"data":       "LimitDATA",
	"fsize":      "LimitFSIZE",
	"memlock":    "LimitMEMLOCK",
	"msgqueue":   "LimitMSGQUEUE",
	"nice":       "LimitNICE",
	"nofile":     "LimitNOFILE",
	"nproc":      "LimitNPROC",
	"rss":        "LimitRSS",
	"rtprio":     "LimitRTPRIO",
	"sigpending": "LimitSIGPENDING",
	"stack":      "LimitSTACK",
}

// TODO(ericsnow) Move normalize to common.Conf.Normalize.

type confRenderer interface {
	shell.Renderer
	shell.ScriptRenderer
}

// normalize adjusts the conf to more standardized content and
// returns a new Conf with that updated content. It also returns the
// content of any script file that should accompany the conf.
func normalize(name string, conf common.Conf, scriptPath string, renderer confRenderer) (common.Conf, []byte) {
	var data []byte

	var cmds []string
	if conf.Logfile != "" {
		filename := conf.Logfile
		cmds = append(cmds, "# Set up logging.")
		cmds = append(cmds, renderer.Touch(filename, nil)...)
		user, group := paths.SyslogUserGroup()
		cmds = append(cmds, renderer.Chown(filename, user, group)...)
		cmds = append(cmds, renderer.Chmod(filename, paths.LogfilePermission)...)
		cmds = append(cmds, renderer.RedirectOutput(filename)...)
		cmds = append(cmds, renderer.RedirectFD("out", "err")...)
		cmds = append(cmds,
			"",
			"# Run the script.",
		)
		// We leave conf.Logfile alone (it will be ignored during validation).
	}
	cmds = append(cmds, conf.ExecStart)

	if conf.ExtraScript != "" {
		cmds = append([]string{conf.ExtraScript}, cmds...)
		conf.ExtraScript = ""
	}
	if !isSimpleCommand(strings.Join(cmds, "\n")) {
		data = renderer.RenderScript(cmds)
		conf.ExecStart = scriptPath
	}

	if len(conf.Env) == 0 {
		conf.Env = nil
	}

	if len(conf.Limit) == 0 {
		conf.Limit = nil
	}

	if conf.Transient {
		// TODO(ericsnow) Handle Transient via systemd-run command?
		conf.ExecStopPost = commands{}.disable(name)
	}

	return conf, data
}

func isSimpleCommand(cmd string) bool {
	if strings.ContainsAny(cmd, "\n;|><&") {
		return false
	}

	return true
}

func validate(name string, conf common.Conf, renderer shell.Renderer) error {
	if name == "" {
		return errors.NotValidf("missing service name")
	}

	if err := conf.Validate(renderer); err != nil {
		return errors.Trace(err)
	}

	if conf.ExtraScript != "" {
		return errors.NotValidf("unexpected ExtraScript")
	}

	// We ignore Desc and Logfile.

	for k := range conf.Limit {
		if _, ok := limitMap[k]; !ok {
			return errors.NotValidf("conf.Limit key %q", k)
		}
	}

	return nil
}

// serialize returns the data that should be written to disk for the
// provided Conf, rendered in the systemd unit file format.
func serialize(name string, conf common.Conf, renderer shell.Renderer) ([]byte, error) {
	if err := validate(name, conf, renderer); err != nil {
		return nil, errors.Trace(err)
	}

	var unitOptions []*unit.UnitOption
	unitOptions = append(unitOptions, serializeUnit(conf)...)
	unitOptions = append(unitOptions, serializeService(conf)...)
	unitOptions = append(unitOptions, serializeInstall(conf)...)
	// Don't use unit.Serialize because it has map ordering issues.
	// Serialize copied locally, and outputs sections in alphabetical order.
	data, err := ioutil.ReadAll(UnitSerialize(unitOptions))
	return data, errors.Trace(err)
}

func serializeUnit(conf common.Conf) []*unit.UnitOption {
	var unitOptions []*unit.UnitOption

	if conf.Desc != "" {
		unitOptions = append(unitOptions, &unit.UnitOption{
			Section: "Unit",
			Name:    "Description",
			Value:   conf.Desc,
		})
	}

	after := []string{
		"syslog.target",
		"network.target",
		"systemd-user-sessions.service",
	}
	for _, name := range after {
		unitOptions = append(unitOptions, &unit.UnitOption{
			Section: "Unit",
			Name:    "After",
			Value:   name,
		})
	}

	if conf.AfterStopped != "" {
		unitOptions = append(unitOptions, &unit.UnitOption{
			Section: "Unit",
			Name:    "After",
			Value:   conf.AfterStopped,
		})
	}

	return unitOptions
}

func iterSortedKeys(in map[string]string) []string {
	out := make([]string, 0, len(in))
	for k := range in {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// ServiceLimits converts the limits in conf to systemd unit options.
func ServiceLimits(conf common.Conf) []*unit.UnitOption {
	var unitOptions []*unit.UnitOption
	for _, k := range iterSortedKeys(conf.Limit) {
		v := conf.Limit[k]
		if v == "unlimited" {
			// In ulimit you pass 'unlimited', but this maps to "infinity" in systemd
			v = "infinity"
		}
		unitOptions = append(unitOptions, &unit.UnitOption{
			Section: "Service",
			Name:    limitMap[k],
			Value:   v,
		})
	}
	return unitOptions
}

func serializeService(conf common.Conf) []*unit.UnitOption {
	var unitOptions []*unit.UnitOption

	// TODO(ericsnow) Support "Type" (e.g. "forking")? For now we just
	// use the default, "simple".

	for k, v := range conf.Env {
		unitOptions = append(unitOptions, &unit.UnitOption{
			Section: "Service",
			Name:    "Environment",
			Value:   fmt.Sprintf(`"%s=%s"`, k, v),
		})
	}
	unitOptions = append(unitOptions, ServiceLimits(conf)...)

	if conf.ExecStart != "" {
		unitOptions = append(unitOptions, &unit.UnitOption{
			Section: "Service",
			Name:    "ExecStart",
			Value:   conf.ExecStart,
		})
	}

	// TODO(ericsnow) This should key off Conf.Restart, once added.
	if !conf.Transient {
		unitOptions = append(unitOptions, &unit.UnitOption{
			Section: "Service",
			Name:    "Restart",
			Value:   "on-failure",
		})
	}

	if conf.Timeout > 0 {
		unitOptions = append(unitOptions, &unit.UnitOption{
			Section: "Service",
			Name:    "TimeoutSec",
			Value:   strconv.Itoa(conf.Timeout),
		})
	}

	if conf.ExecStopPost != "" {
		unitOptions = append(unitOptions, &unit.UnitOption{
			Section: "Service",
			Name:    "ExecStopPost",
			Value:   conf.ExecStopPost,
		})
	}

	return unitOptions
}

func serializeInstall(conf common.Conf) []*unit.UnitOption {
	var unitOptions []*unit.UnitOption

	unitOptions = append(unitOptions, &unit.UnitOption{
		Section: "Install",
		Name:    "WantedBy",
		Value:   "multi-user.target",
	})

	return unitOptions
}

// deserialize parses the provided data (in the systemd unit file
// format) and populates a new Conf with the result.
func deserialize(data []byte, renderer shell.Renderer) (common.Conf, error) {
	opts, err := unit.Deserialize(bytes.NewBuffer(data))
	if err != nil {
		return common.Conf{}, errors.Trace(err)
	}
	return deserializeOptions(opts, renderer)
}

func deserializeOptions(opts []*unit.UnitOption, renderer shell.Renderer) (common.Conf, error) {
	var conf common.Conf

	for _, uo := range opts {
		switch uo.Section {
		case "Unit":
			switch uo.Name {
			case "Description":
				conf.Desc = uo.Value
			case "After":
				// Do nothing until we support it in common.Conf.
			default:
				return conf, errors.NotSupportedf("Unit directive %q", uo.Name)
			}
		case "Service":
			switch {
			case uo.Name == "ExecStart":
				conf.ExecStart = uo.Value
			case uo.Name == "Environment":
				if conf.Env == nil {
					conf.Env = make(map[string]string)
				}
				var value = uo.Value
				if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
					value = value[1 : len(value)-1]
				}
				parts := strings.SplitN(value, "=", 2)
				if len(parts) != 2 {
					return conf, errors.NotValidf("service environment value %q", uo.Value)
				}
				conf.Env[parts[0]] = parts[1]
			case strings.HasPrefix(uo.Name, "Limit"):
				if conf.Limit == nil {
					conf.Limit = make(map[string]string)
				}
				for k, v := range limitMap {
					if v == uo.Name {
						conf.Limit[k] = uo.Value
						break
					}
				}
			case uo.Name == "TimeoutSec":
				timeout, err := strconv.Atoi(uo.Value)
				if err != nil {
					return conf, errors.Trace(err)
				}
				conf.Timeout = timeout
			case uo.Name == "Type":
				// Do nothing until we support it in common.Conf.
			case uo.Name == "RemainAfterExit":
				// Do nothing until we support it in common.Conf.
			case uo.Name == "Restart":
				// Do nothing until we support it in common.Conf.
			default:
				return conf, errors.NotSupportedf("Service directive %q", uo.Name)
			}
		case "Install":
			switch uo.Name {
			case "WantedBy":
				if uo.Value != "multi-user.target" {
					return conf, errors.NotValidf("unit target %q", uo.Value)
				}
			default:
				return conf, errors.NotSupportedf("Install directive %q", uo.Name)
			}
		default:
			return conf, errors.NotSupportedf("section %q", uo.Name)
		}
	}

	err := validate("<>", conf, renderer)
	return conf, errors.Trace(err)
}
