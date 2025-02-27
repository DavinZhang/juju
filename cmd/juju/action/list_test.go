// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package action_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/juju/cmd/v3"
	"github.com/juju/cmd/v3/cmdtesting"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	actionapi "github.com/DavinZhang/juju/api/action"
	"github.com/DavinZhang/juju/cmd/juju/action"
)

type ListSuite struct {
	BaseActionSuite
	wrappedCommand cmd.Command
	command        *action.ListCommand
}

var _ = gc.Suite(&ListSuite{})

func (s *ListSuite) SetUpTest(c *gc.C) {
	s.BaseActionSuite.SetUpTest(c)
	s.wrappedCommand, s.command = action.NewListCommandForTest(s.store)
}

func (s *ListSuite) TestInit(c *gc.C) {
	tests := []struct {
		should               string
		args                 []string
		expectedApp          string
		expectedOutputSchema bool
		expectedErr          string
	}{{
		should:      "fail with missing application name",
		args:        []string{},
		expectedErr: "no application name specified",
	}, {
		should:      "fail with invalid application name",
		args:        []string{invalidApplicationId},
		expectedErr: "invalid application name \"" + invalidApplicationId + "\"",
	}, {
		should:      "fail with too many args",
		args:        []string{"two", "things"},
		expectedErr: "unrecognized args: \\[\"things\"\\]",
	}, {
		should:      "init properly with valid application name",
		args:        []string{validApplicationId},
		expectedApp: validApplicationId,
	}, {
		should:      "schema with tabular output",
		args:        []string{"--format=tabular", "--schema", validApplicationId},
		expectedErr: "full schema not compatible with tabular output",
	}, {
		should:               "init properly with valid application name and --schema",
		args:                 []string{"--format=yaml", "--schema", validApplicationId},
		expectedOutputSchema: true,
		expectedApp:          validApplicationId,
	}, {
		should:               "default to yaml output when --schema option is specified",
		args:                 []string{"--schema", validApplicationId},
		expectedOutputSchema: true,
		expectedApp:          validApplicationId,
	}}

	for i, t := range tests {
		for _, modelFlag := range s.modelFlags {
			c.Logf("test %d should %s: juju actions defined %s", i,
				t.should, strings.Join(t.args, " "))
			s.wrappedCommand, s.command = action.NewListCommandForTest(s.store)
			args := append([]string{modelFlag, "admin"}, t.args...)
			err := cmdtesting.InitCommand(s.wrappedCommand, args)
			if t.expectedErr == "" {
				c.Check(err, jc.ErrorIsNil)
				c.Check(s.command.ApplicationName(), gc.Equals, t.expectedApp)
				c.Check(s.command.FullSchema(), gc.Equals, t.expectedOutputSchema)
			} else {
				c.Check(err, gc.ErrorMatches, t.expectedErr)
			}
		}
	}
}

func (s *ListSuite) TestRun(c *gc.C) {
	simpleOutput := `
Action          Description
kill            Kill the database.
no-description  No description
no-params       An action with no parameters.
snapshot        Take a snapshot of the database.
`[1:]

	tests := []struct {
		should           string
		expectFullSchema bool
		expectNoResults  bool
		expectMessage    string
		withArgs         []string
		withAPIErr       string
		withCharmActions map[string]actionapi.ActionSpec
		expectedErr      string
	}{{
		should:      "pass back API error correctly",
		withArgs:    []string{validApplicationId},
		withAPIErr:  "an API error",
		expectedErr: "an API error",
	}, {
		should:           "get short results properly",
		withArgs:         []string{validApplicationId},
		withCharmActions: someCharmActions,
	}, {
		should:           "get full schema results properly",
		withArgs:         []string{"--format=yaml", "--schema", validApplicationId},
		expectFullSchema: true,
		withCharmActions: someCharmActions,
	}, {
		should:          "work properly when no results found",
		withArgs:        []string{validApplicationId},
		expectNoResults: true,
		expectMessage:   fmt.Sprintf("No actions defined for %s.\n", validApplicationId),
	}, {
		should:           "get tabular default output when --schema is NOT specified",
		withArgs:         []string{"--format=default", validApplicationId},
		withCharmActions: someCharmActions,
	}, {
		should:           "get full schema default output (YAML) when --schema is specified",
		withArgs:         []string{"--format=default", "--schema", validApplicationId},
		expectFullSchema: true,
		withCharmActions: someCharmActions,
	}}

	for i, t := range tests {
		for _, modelFlag := range s.modelFlags {
			func() {
				c.Logf("test %d should %s", i, t.should)

				fakeClient := &fakeAPIClient{charmActions: t.withCharmActions}
				if t.withAPIErr != "" {
					fakeClient.apiErr = errors.New(t.withAPIErr)
				}
				restore := s.patchAPIClient(fakeClient)
				defer restore()

				args := append([]string{modelFlag, "admin"}, t.withArgs...)
				s.wrappedCommand, s.command = action.NewListCommandForTest(s.store)
				ctx, err := cmdtesting.RunCommand(c, s.wrappedCommand, args...)

				if t.expectedErr != "" || t.withAPIErr != "" {
					c.Check(err, gc.ErrorMatches, t.expectedErr)
				} else {
					c.Assert(err, gc.IsNil)
					result := ctx.Stdout.(*bytes.Buffer).Bytes()
					if t.expectFullSchema {
						checkFullSchema(c, t.withCharmActions, result)
					} else if t.expectNoResults {
						c.Check(cmdtesting.Stderr(ctx), gc.Matches, t.expectMessage)
					} else {
						c.Check(cmdtesting.Stdout(ctx), gc.Equals, simpleOutput)
					}
				}

			}()
		}
	}
}

func checkFullSchema(c *gc.C, expected map[string]actionapi.ActionSpec, actual []byte) {
	expectedOutput := make(map[string]interface{})
	for k, v := range expected {
		expectedOutput[k] = v.Params
	}
	c.Check(string(actual), jc.YAMLEquals, expectedOutput)
}
