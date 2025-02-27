// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package keymanager_test

import (
	"strings"

	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v2/ssh"
	sshtesting "github.com/juju/utils/v2/ssh/testing"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/api/keymanager"
	keymanagerserver "github.com/DavinZhang/juju/apiserver/facades/client/keymanager"
	keymanagertesting "github.com/DavinZhang/juju/apiserver/facades/client/keymanager/testing"
	"github.com/DavinZhang/juju/apiserver/params"
	jujutesting "github.com/DavinZhang/juju/juju/testing"
	"github.com/DavinZhang/juju/rpc"
)

type keymanagerSuite struct {
	jujutesting.JujuConnSuite

	keymanager *keymanager.Client
}

var _ = gc.Suite(&keymanagerSuite{})

func (s *keymanagerSuite) SetUpTest(c *gc.C) {
	s.JujuConnSuite.SetUpTest(c)
	s.keymanager = keymanager.NewClient(s.APIState)
	c.Assert(s.keymanager, gc.NotNil)

}

func (s *keymanagerSuite) setAuthorisedKeys(c *gc.C, keys string) {
	err := s.Model.UpdateModelConfig(map[string]interface{}{"authorized-keys": keys}, nil)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *keymanagerSuite) TestListKeys(c *gc.C) {
	key1 := sshtesting.ValidKeyOne.Key + " user@host"
	key2 := sshtesting.ValidKeyTwo.Key
	s.setAuthorisedKeys(c, strings.Join([]string{key1, key2}, "\n"))

	keyResults, err := s.keymanager.ListKeys(ssh.Fingerprints, s.AdminUserTag(c).Name())
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(len(keyResults), gc.Equals, 1)
	result := keyResults[0]
	c.Assert(result.Error, gc.IsNil)
	c.Assert(result.Result, gc.DeepEquals,
		[]string{sshtesting.ValidKeyOne.Fingerprint + " (user@host)", sshtesting.ValidKeyTwo.Fingerprint})
}

func (s *keymanagerSuite) TestListKeysErrors(c *gc.C) {
	c.Skip("the user name isn't checked for existence yet")
	keyResults, err := s.keymanager.ListKeys(ssh.Fingerprints, "invalid")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(len(keyResults), gc.Equals, 1)
	result := keyResults[0]
	c.Assert(result.Error, gc.ErrorMatches, `permission denied`)
}

func clientError(message string) *params.Error {
	return &params.Error{
		Message: message,
		Code:    "",
	}
}

func (s *keymanagerSuite) assertModelKeys(c *gc.C, expected []string) {
	modelConfig, err := s.Model.ModelConfig()
	c.Assert(err, jc.ErrorIsNil)
	keys := modelConfig.AuthorizedKeys()
	c.Assert(keys, gc.Equals, strings.Join(expected, "\n"))
}

func (s *keymanagerSuite) TestAddKeys(c *gc.C) {
	key1 := sshtesting.ValidKeyOne.Key + " user@host"
	s.setAuthorisedKeys(c, key1)

	newKeys := []string{sshtesting.ValidKeyTwo.Key, sshtesting.ValidKeyThree.Key, "invalid"}
	errResults, err := s.keymanager.AddKeys(s.AdminUserTag(c).Name(), newKeys...)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(errResults, gc.DeepEquals, []params.ErrorResult{
		{Error: nil},
		{Error: nil},
		{Error: clientError("invalid ssh key: invalid")},
	})
	s.assertModelKeys(c, append([]string{key1}, newKeys[:2]...))
}

func (s *keymanagerSuite) TestAddSystemKeyForbidden(c *gc.C) {
	key1 := sshtesting.ValidKeyOne.Key + " user@host"
	s.setAuthorisedKeys(c, key1)

	newKey := sshtesting.ValidKeyTwo.Key
	_, err := s.keymanager.AddKeys("juju-system-key", newKey)
	c.Assert(errors.Cause(err), gc.DeepEquals, &rpc.RequestError{
		Message: "permission denied",
		Code:    "unauthorized access",
	})
	s.assertModelKeys(c, []string{key1})
}

func (s *keymanagerSuite) TestDeleteKeys(c *gc.C) {
	key1 := sshtesting.ValidKeyOne.Key + " user@host"
	key2 := sshtesting.ValidKeyTwo.Key
	key3 := sshtesting.ValidKeyThree.Key
	initialKeys := []string{key1, key2, key3, "invalid"}
	s.setAuthorisedKeys(c, strings.Join(initialKeys, "\n"))

	errResults, err := s.keymanager.DeleteKeys(s.AdminUserTag(c).Name(), sshtesting.ValidKeyTwo.Fingerprint, "user@host", "missing")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(errResults, gc.DeepEquals, []params.ErrorResult{
		{Error: nil},
		{Error: nil},
		{Error: clientError("invalid ssh key: missing")},
	})
	s.assertModelKeys(c, []string{key3, "invalid"})
}

func (s *keymanagerSuite) TestImportKeys(c *gc.C) {
	s.PatchValue(&keymanagerserver.RunSSHImportId, keymanagertesting.FakeImport)

	key1 := sshtesting.ValidKeyOne.Key + " user@host"
	s.setAuthorisedKeys(c, key1)

	keyIds := []string{"lp:validuser", "invalid-key"}
	errResults, err := s.keymanager.ImportKeys(s.AdminUserTag(c).Name(), keyIds...)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(errResults, gc.DeepEquals, []params.ErrorResult{
		{Error: nil},
		{Error: clientError("invalid ssh key id: invalid-key")},
	})
	s.assertModelKeys(c, []string{key1, sshtesting.ValidKeyThree.Key})
}

func (s *keymanagerSuite) assertInvalidUserOperation(c *gc.C, test func(user string, keys []string) error) {
	key1 := sshtesting.ValidKeyOne.Key + " user@host"
	s.setAuthorisedKeys(c, key1)

	// Run the required test code and check the error.
	keys := []string{sshtesting.ValidKeyTwo.Key, sshtesting.ValidKeyThree.Key}
	err := test("invalid", keys)
	c.Assert(err, gc.ErrorMatches, `permission denied`)

	// No model changes.
	s.assertModelKeys(c, []string{key1})
}

func (s *keymanagerSuite) TestAddKeysInvalidUser(c *gc.C) {
	c.Skip("no user validation done yet")
	s.assertInvalidUserOperation(c, func(user string, keys []string) error {
		_, err := s.keymanager.AddKeys(user, keys...)
		return err
	})
}

func (s *keymanagerSuite) TestDeleteKeysInvalidUser(c *gc.C) {
	c.Skip("no user validation done yet")
	s.assertInvalidUserOperation(c, func(user string, keys []string) error {
		_, err := s.keymanager.DeleteKeys(user, keys...)
		return err
	})
}

func (s *keymanagerSuite) TestImportKeysInvalidUser(c *gc.C) {
	c.Skip("no user validation done yet")
	s.assertInvalidUserOperation(c, func(user string, keys []string) error {
		_, err := s.keymanager.ImportKeys(user, keys...)
		return err
	})
}

func (s *keymanagerSuite) TestExposesBestAPIVersion(c *gc.C) {
	c.Check(s.keymanager.BestAPIVersion(), gc.Equals, 1)
}
