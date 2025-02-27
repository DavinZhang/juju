// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package provider_test

import (
	"path/filepath"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v2"
	gc "gopkg.in/check.v1"

	k8s "github.com/DavinZhang/juju/caas/kubernetes"
	"github.com/DavinZhang/juju/caas/kubernetes/provider"
	"github.com/DavinZhang/juju/cloud"
	"github.com/DavinZhang/juju/environs"
	envtesting "github.com/DavinZhang/juju/environs/testing"
)

type credentialsSuite struct {
	testing.FakeHomeSuite
	provider environs.EnvironProvider
}

var _ = gc.Suite(&credentialsSuite{})

func (s *credentialsSuite) SetUpTest(c *gc.C) {
	s.FakeHomeSuite.SetUpTest(c)

	var err error
	s.provider, err = environs.Provider("kubernetes")
	c.Assert(err, jc.ErrorIsNil)
}

func (s *credentialsSuite) TestCredentialSchemas(c *gc.C) {
	envtesting.AssertProviderAuthTypes(c, s.provider, "userpass", "oauth2", "clientcertificate", "oauth2withcert", "certificate")
}

func (s *credentialsSuite) TestCredentialsValid(c *gc.C) {
	envtesting.AssertProviderCredentialsValid(c, s.provider, "userpass", map[string]string{
		"username": "fred",
		"password": "secret",
	})
}

func (s *credentialsSuite) TestHiddenAttributes(c *gc.C) {
	envtesting.AssertProviderCredentialsAttributesHidden(c, s.provider, "userpass", "password")
	envtesting.AssertProviderCredentialsAttributesHidden(c, s.provider, "oauth2", "Token")
	envtesting.AssertProviderCredentialsAttributesHidden(c, s.provider, "clientcertificate", "ClientKeyData")
	envtesting.AssertProviderCredentialsAttributesHidden(c, s.provider, "oauth2withcert", "ClientKeyData", "Token")
	envtesting.AssertProviderCredentialsAttributesHidden(c, s.provider, "certificate", "Token")
}

var singleConfigYAML = `
apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://1.1.1.1:8888
    certificate-authority-data: QQ==
  name: the-cluster
contexts:
- context:
    cluster: the-cluster
    user: the-user
  name: the-context
current-context: the-context
preferences: {}
users:
- name: the-user
  user:
    password: thepassword
    username: theuser
`

func (s *credentialsSuite) TestDetectCredentials(c *gc.C) {
	kubeConfig := filepath.Join(utils.Home(), "config")
	s.PatchEnvironment("KUBECONFIG", kubeConfig)
	s.Home.AddFiles(c, testing.TestFile{
		Name: "config",
		Data: singleConfigYAML,
	})
	creds, err := s.provider.DetectCredentials("")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(creds.DefaultRegion, gc.Equals, "")
	expected := cloud.NewNamedCredential(
		"the-user", cloud.UserPassAuthType, map[string]string{
			"username": "theuser",
			"password": "thepassword",
		}, false,
	)
	c.Assert(creds.AuthCredentials["the-user"], jc.DeepEquals, expected)
}

func (s *credentialsSuite) TestRegisterCredentialsNotMicrok8s(c *gc.C) {
	p := provider.NewProviderCredentials(credentialGetterFunc(builtinCloudRet{}))
	credentials, err := p.RegisterCredentials(cloud.Cloud{})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(credentials, gc.HasLen, 0)
}

func (s *credentialsSuite) TestRegisterCredentialsMicrok8s(c *gc.C) {
	p := provider.NewProviderCredentials(
		credentialGetterFunc(
			builtinCloudRet{
				cloud:      defaultK8sCloud,
				credential: getDefaultCredential(),
				err:        nil,
			},
		),
	)
	credentials, err := p.RegisterCredentials(defaultK8sCloud)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(credentials, gc.HasLen, 1)
	c.Assert(credentials[k8s.K8sCloudMicrok8s], gc.DeepEquals, &cloud.CloudCredential{
		DefaultCredential: k8s.K8sCloudMicrok8s,
		AuthCredentials: map[string]cloud.Credential{
			k8s.K8sCloudMicrok8s: getDefaultCredential(),
		},
	})
}
