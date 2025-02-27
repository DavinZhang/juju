// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package controller_test

import (
	"testing"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/jujuclient"
	"github.com/DavinZhang/juju/jujuclient/jujuclienttesting"
	coretesting "github.com/DavinZhang/juju/testing"
)

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}

type baseControllerSuite struct {
	coretesting.FakeJujuXDGDataHomeSuite
	store                                     jujuclient.ClientStore
	controllersYaml, modelsYaml, accountsYaml string
	expectedOutput, expectedErr               string
}

func (s *baseControllerSuite) SetUpTest(c *gc.C) {
	s.FakeJujuXDGDataHomeSuite.SetUpTest(c)
	s.controllersYaml = testControllersYaml
	s.modelsYaml = testModelsYaml
	s.accountsYaml = testAccountsYaml
	s.store = jujuclienttesting.MinimalStore()
}

func (s *baseControllerSuite) createTestClientStore(c *gc.C) *jujuclient.MemStore {
	controllers, err := jujuclient.ParseControllers([]byte(s.controllersYaml))
	c.Assert(err, jc.ErrorIsNil)

	models, err := jujuclient.ParseModels([]byte(s.modelsYaml))
	c.Assert(err, jc.ErrorIsNil)

	accounts, err := jujuclient.ParseAccounts([]byte(s.accountsYaml))
	c.Assert(err, jc.ErrorIsNil)

	store := jujuclient.NewMemStore()
	store.Controllers = controllers.Controllers
	store.CurrentControllerName = controllers.CurrentController
	store.Models = models
	store.Accounts = accounts
	s.store = store
	return store
}

const testControllersYaml = `
controllers:
  aws-test:
    uuid: this-is-the-aws-test-uuid
    api-endpoints: [this-is-aws-test-of-many-api-endpoints]
    cloud: aws
    region: us-east-1
    model-count: 2
    machine-count: 5
    agent-version: 2.0.1
    ca-cert: this-is-aws-test-ca-cert
  mallards:
    uuid: this-is-another-uuid
    api-endpoints: [this-is-another-of-many-api-endpoints, this-is-one-more-of-many-api-endpoints]
    cloud: mallards
    region: mallards1
    ca-cert: this-is-another-ca-cert
  mark-test-prodstack:
    uuid: this-is-a-uuid
    api-endpoints: [this-is-one-of-many-api-endpoints]
    cloud: prodstack
    ca-cert: this-is-a-ca-cert
  k8s-controller:
    uuid: this-is-a-k8s-uuid
    api-endpoints: [this-is-one-of-many-k8s-api-endpoints]
    cloud: microk8s
    region: localhost
    type: kubernetes
    ca-cert: this-is-a-k8s-ca-cert
    machine-count: 3
    agent-version: 6.6.6
current-controller: mallards
`

const testModelsYaml = `
controllers:
  aws-test:
    models:
      controller:
        uuid: ghi
        type: iaas
    current-model: admin/controller
  mallards:
    models:
      model0:
        uuid: abc
        type: iaas
      my-model:
        uuid: def
        type: iaas
    current-model: admin/my-model
  k8s-controller:
    models:
      controller:
        uuid: xyz
        type: caas
      my-k8s-model:
        uuid: def
        type: caas
    current-model: admin/my-k8s-model
`

const testAccountsYaml = `
controllers:
  aws-test:
    user: admin
    password: hun+er2
  mark-test-prodstack:
    user: admin
    password: hunter2
  mallards:
    user: admin
    password: hunter2
    last-known-access: superuser
  k8s-controller:
    user: admin
    password: hunter2
    last-known-access: superuser
`
