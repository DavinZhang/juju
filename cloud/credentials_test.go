// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package cloud_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v2"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/cloud"
	"github.com/DavinZhang/juju/testing"
)

type credentialsSuite struct {
	testing.FakeJujuXDGDataHomeSuite
}

var _ = gc.Suite(&credentialsSuite{})

func (s *credentialsSuite) TestMarshalAccessKey(c *gc.C) {
	creds := map[string]cloud.CloudCredential{
		"aws": {
			DefaultCredential: "default-cred",
			DefaultRegion:     "us-west-2",
			AuthCredentials: map[string]cloud.Credential{
				"peter": cloud.NewCredential(cloud.AccessKeyAuthType, map[string]string{
					"access-key": "key",
					"secret-key": "secret",
				}),
				// TODO(wallyworld) - add anther credential once goyaml.v2 supports inline MapSlice.
				//"paul": &cloud.AccessKeyCredentials{
				//	Key: "paulkey",
				//	Secret: "paulsecret",
				//},
			},
		},
	}
	out, err := cloud.MarshalCredentials(creds)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(out), gc.Equals, `
credentials:
  aws:
    default-credential: default-cred
    default-region: us-west-2
    peter:
      auth-type: access-key
      access-key: key
      secret-key: secret
`[1:])
}

func (s *credentialsSuite) TestMarshalOpenstackAccessKey(c *gc.C) {
	creds := map[string]cloud.CloudCredential{
		"openstack": {
			DefaultCredential: "default-cred",
			DefaultRegion:     "region-a",
			AuthCredentials: map[string]cloud.Credential{
				"peter": cloud.NewCredential(cloud.AccessKeyAuthType, map[string]string{
					"access-key":  "key",
					"secret-key":  "secret",
					"tenant-name": "tenant",
				}),
			},
		},
	}
	out, err := cloud.MarshalCredentials(creds)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(out), gc.Equals, `
credentials:
  openstack:
    default-credential: default-cred
    default-region: region-a
    peter:
      auth-type: access-key
      access-key: key
      secret-key: secret
      tenant-name: tenant
`[1:])
}

func (s *credentialsSuite) TestMarshalOpenstackUserPass(c *gc.C) {
	creds := map[string]cloud.CloudCredential{
		"openstack": {
			DefaultCredential: "default-cred",
			DefaultRegion:     "region-a",
			AuthCredentials: map[string]cloud.Credential{
				"peter": cloud.NewCredential(cloud.UserPassAuthType, map[string]string{
					"username":    "user",
					"password":    "secret",
					"tenant-name": "tenant",
				}),
			},
		},
	}
	out, err := cloud.MarshalCredentials(creds)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(out), gc.Equals, `
credentials:
  openstack:
    default-credential: default-cred
    default-region: region-a
    peter:
      auth-type: userpass
      password: secret
      tenant-name: tenant
      username: user
`[1:])
}

func (s *credentialsSuite) TestMarshalAzureCredntials(c *gc.C) {
	creds := map[string]cloud.CloudCredential{
		"azure": {
			DefaultCredential: "default-cred",
			DefaultRegion:     "Central US",
			AuthCredentials: map[string]cloud.Credential{
				"peter": cloud.NewCredential(cloud.UserPassAuthType, map[string]string{
					"application-id":       "app-id",
					"application-password": "app-secret",
					"subscription-id":      "subscription-id",
					"tenant-id":            "tenant-id",
				}),
			},
		},
	}
	out, err := cloud.MarshalCredentials(creds)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(out), gc.Equals, `
credentials:
  azure:
    default-credential: default-cred
    default-region: Central US
    peter:
      auth-type: userpass
      application-id: app-id
      application-password: app-secret
      subscription-id: subscription-id
      tenant-id: tenant-id
`[1:])
}

func (s *credentialsSuite) TestMarshalOAuth1(c *gc.C) {
	creds := map[string]cloud.CloudCredential{
		"maas": {
			DefaultCredential: "default-cred",
			DefaultRegion:     "region-default",
			AuthCredentials: map[string]cloud.Credential{
				"peter": cloud.NewCredential(cloud.OAuth1AuthType, map[string]string{
					"consumer-key":    "consumer-key",
					"consumer-secret": "consumer-secret",
					"access-token":    "access-token",
					"token-secret":    "token-secret",
				}),
			},
		},
	}
	out, err := cloud.MarshalCredentials(creds)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(out), gc.Equals, `
credentials:
  maas:
    default-credential: default-cred
    default-region: region-default
    peter:
      auth-type: oauth1
      access-token: access-token
      consumer-key: consumer-key
      consumer-secret: consumer-secret
      token-secret: token-secret
`[1:])
}

func (s *credentialsSuite) TestMarshalOAuth2(c *gc.C) {
	creds := map[string]cloud.CloudCredential{
		"google": {
			DefaultCredential: "default-cred",
			DefaultRegion:     "West US",
			AuthCredentials: map[string]cloud.Credential{
				"peter": cloud.NewCredential(cloud.OAuth2AuthType, map[string]string{
					"client-id":    "client-id",
					"client-email": "client-email",
					"private-key":  "secret",
				}),
			},
		},
	}
	out, err := cloud.MarshalCredentials(creds)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(out), gc.Equals, `
credentials:
  google:
    default-credential: default-cred
    default-region: West US
    peter:
      auth-type: oauth2
      client-email: client-email
      client-id: client-id
      private-key: secret
`[1:])
}

func (s *credentialsSuite) TestParseCredentials(c *gc.C) {
	s.testParseCredentials(c, []byte(`
credentials:
  aws:
    default-credential: peter
    default-region: us-east-2
    peter:
      auth-type: access-key
      access-key: key
      secret-key: secret
  aws-china:
    default-credential: zhu8jie
    zhu8jie:
      auth-type: access-key
      access-key: key
      secret-key: secret
    sun5kong:
      auth-type: access-key
      access-key: quay
      secret-key: sekrit
  aws-gov:
    default-region: us-gov-west-1
    supersekrit:
      auth-type: access-key
      access-key: super
      secret-key: sekrit
`[1:]), map[string]cloud.CloudCredential{
		"aws": {
			DefaultCredential: "peter",
			DefaultRegion:     "us-east-2",
			AuthCredentials: map[string]cloud.Credential{
				"peter": cloud.NewCredential(cloud.AccessKeyAuthType, map[string]string{
					"access-key": "key",
					"secret-key": "secret",
				}),
			},
		},
		"aws-china": {
			DefaultCredential: "zhu8jie",
			AuthCredentials: map[string]cloud.Credential{
				"zhu8jie": cloud.NewCredential(cloud.AccessKeyAuthType, map[string]string{
					"access-key": "key",
					"secret-key": "secret",
				}),
				"sun5kong": cloud.NewCredential(cloud.AccessKeyAuthType, map[string]string{
					"access-key": "quay",
					"secret-key": "sekrit",
				}),
			},
		},
		"aws-gov": {
			DefaultRegion: "us-gov-west-1",
			AuthCredentials: map[string]cloud.Credential{
				"supersekrit": cloud.NewCredential(cloud.AccessKeyAuthType, map[string]string{
					"access-key": "super",
					"secret-key": "sekrit",
				}),
			},
		},
	})
}

func (s *credentialsSuite) TestParseCredentialsUnknownAuthType(c *gc.C) {
	// Unknown auth-type is not validated by ParseCredentials.
	// Validation is deferred to FinalizeCredential.
	s.testParseCredentials(c, []byte(`
credentials:
  cloud-name:
    credential-name:
      auth-type: woop
`[1:]), map[string]cloud.CloudCredential{
		"cloud-name": {
			AuthCredentials: map[string]cloud.Credential{
				"credential-name": cloud.NewCredential("woop", nil),
			},
		},
	})
}

func (s *credentialsSuite) testParseCredentials(c *gc.C, input []byte, expect map[string]cloud.CloudCredential) {
	output, err := cloud.ParseCredentials(input)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(output, jc.DeepEquals, expect)
}

func (s *credentialsSuite) TestParseCredentialsMissingAuthType(c *gc.C) {
	s.testParseCredentialsError(c, []byte(`
credentials:
  cloud-name:
    credential-name:
      doesnt: really-matter
`[1:]), "credentials.cloud-name.credential-name: missing auth-type")
}

func (s *credentialsSuite) TestParseCredentialsNonStringValue(c *gc.C) {
	s.testParseCredentialsError(c, []byte(`
credentials:
  cloud-name:
    credential-name:
      non-string-value: 123
`[1:]), `credentials\.cloud-name\.credential-name\.non-string-value: expected string, got int\(123\)`)
}

func (s *credentialsSuite) testParseCredentialsError(c *gc.C, input []byte, expect string) {
	_, err := cloud.ParseCredentials(input)
	c.Assert(err, gc.ErrorMatches, expect)
}

func (s *credentialsSuite) TestFinalizeCredential(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.UserPassAuthType,
		map[string]string{
			"key": "value",
		},
	)
	schema := cloud.CredentialSchema{{
		"key",
		cloud.CredentialAttr{
			Description: "key credential",
			Hidden:      true,
		},
	}}
	_, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.UserPassAuthType: schema,
	}, readFileNotSupported)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *credentialsSuite) TestFinalizeCredentialFileAttr(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.UserPassAuthType,
		map[string]string{
			"key-file": "path",
			"quay":     "value",
		},
	)
	schema := cloud.CredentialSchema{{
		"key",
		cloud.CredentialAttr{
			Description: "key credential",
			Hidden:      true,
			FileAttr:    "key-file",
		},
	}, {
		"quay", cloud.CredentialAttr{FileAttr: "quay-file"},
	}}
	readFile := func(s string) ([]byte, error) {
		c.Assert(s, gc.Equals, "path")
		return []byte("file-value"), nil
	}
	newCred, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.UserPassAuthType: schema,
	}, readFile)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(newCred.Attributes(), jc.DeepEquals, map[string]string{
		"key":  "file-value",
		"quay": "value",
	})
}

func (s *credentialsSuite) TestFinalizeCredentialFileEmpty(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.UserPassAuthType,
		map[string]string{
			"key-file": "path",
		},
	)
	schema := cloud.CredentialSchema{{
		"key",
		cloud.CredentialAttr{
			Description: "key credential",
			Hidden:      true,
			FileAttr:    "key-file",
		},
	}}
	readFile := func(string) ([]byte, error) {
		return nil, nil
	}
	_, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.UserPassAuthType: schema,
	}, readFile)
	c.Assert(err, gc.ErrorMatches, `empty file for "key" not valid`)
}

func (s *credentialsSuite) TestFinalizeCredentialFileAttrNeither(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.UserPassAuthType,
		map[string]string{},
	)
	schema := cloud.CredentialSchema{{
		"key",
		cloud.CredentialAttr{
			Description: "key credential",
			Hidden:      true,
			FileAttr:    "key-file",
		},
	}}
	_, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.UserPassAuthType: schema,
	}, readFileNotSupported)
	c.Assert(err, gc.ErrorMatches, `either "key" or "key-file" must be specified`)
}

func (s *credentialsSuite) TestFinalizeCredentialFileAttrBoth(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.UserPassAuthType,
		map[string]string{
			"key":      "value",
			"key-file": "path",
		},
	)
	schema := cloud.CredentialSchema{{
		"key",
		cloud.CredentialAttr{
			Description: "key credential",
			Hidden:      true,
			FileAttr:    "key-file",
		},
	}}
	_, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.UserPassAuthType: schema,
	}, readFileNotSupported)
	c.Assert(err, gc.ErrorMatches, `specifying both "key" and "key-file" not valid`)
}

func (s *credentialsSuite) TestFinalizeCredentialInvalid(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.UserPassAuthType,
		map[string]string{},
	)
	schema := cloud.CredentialSchema{{
		"key",
		cloud.CredentialAttr{
			Description: "key credential",
			Hidden:      true,
		},
	}}
	_, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.UserPassAuthType: schema,
	}, readFileNotSupported)
	c.Assert(err, gc.ErrorMatches, "key: expected string, got nothing")
}

func (s *credentialsSuite) TestFinalizeCredentialNotSupported(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.OAuth2AuthType,
		map[string]string{},
	)
	_, err := cloud.FinalizeCredential(
		cred, map[cloud.AuthType]cloud.CredentialSchema{}, readFileNotSupported,
	)
	c.Assert(err, jc.Satisfies, errors.IsNotSupported)
	c.Assert(err, gc.ErrorMatches, `auth-type "oauth2" not supported`)
}

func readFileNotSupported(f string) ([]byte, error) {
	return nil, errors.NotSupportedf("reading file %q", f)
}

func (s *credentialsSuite) TestFinalizeCredentialMandatoryFieldMissing(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.UserPassAuthType,
		map[string]string{
			"password": "secret",
			"domain":   "domain",
		},
	)
	schema := cloud.CredentialSchema{
		{"username", cloud.CredentialAttr{Optional: false}},
		{"password", cloud.CredentialAttr{Hidden: true}},
		{"domain", cloud.CredentialAttr{}},
	}
	_, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.UserPassAuthType: schema,
	}, nil)
	c.Assert(err, gc.ErrorMatches, "username: expected string, got nothing")
}

func (s *credentialsSuite) TestFinalizeCredentialMandatoryFieldFromFile(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.UserPassAuthType,
		map[string]string{
			"key-file": "path",
		},
	)
	schema := cloud.CredentialSchema{{
		"key",
		cloud.CredentialAttr{
			Description: "key credential",
			Optional:    false,
			FileAttr:    "key-file",
		},
	}}
	readFile := func(s string) ([]byte, error) {
		c.Assert(s, gc.Equals, "path")
		return []byte("file-value"), nil
	}
	newCred, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.UserPassAuthType: schema,
	}, readFile)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(newCred.Attributes(), jc.DeepEquals, map[string]string{
		"key": "file-value",
	})
}

func (s *credentialsSuite) TestFinalizeCredentialExtraField(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.UserPassAuthType,
		map[string]string{
			"username":   "user",
			"password":   "secret",
			"domain":     "domain",
			"access-key": "access-key",
		},
	)
	schema := cloud.CredentialSchema{
		{"username", cloud.CredentialAttr{Optional: false}},
		{"password", cloud.CredentialAttr{Hidden: true}},
		{"domain", cloud.CredentialAttr{}},
	}
	_, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.UserPassAuthType: schema,
	}, nil)
	c.Assert(err, gc.ErrorMatches, regexp.QuoteMeta(`unknown key "access-key" (value "access-key")`))
}

func (s *credentialsSuite) TestFinalizeCredentialInvalidChoice(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.UserPassAuthType,
		map[string]string{
			"username":  "user",
			"password":  "secret",
			"algorithm": "foo",
		},
	)
	schema := cloud.CredentialSchema{
		{"username", cloud.CredentialAttr{Optional: false}},
		{"password", cloud.CredentialAttr{Hidden: true}},
		{"algorithm", cloud.CredentialAttr{Options: []interface{}{"bar", "foobar"}}},
	}
	_, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.UserPassAuthType: schema,
	}, nil)
	c.Assert(err, gc.ErrorMatches, regexp.QuoteMeta(`algorithm: expected one of [bar foobar], got "foo"`))
}

func (s *credentialsSuite) TestFinalizeCredentialFilePath(c *gc.C) {
	dir := c.MkDir()
	filename := filepath.Join(dir, "filename")
	err := ioutil.WriteFile(filename, []byte{}, 0600)
	c.Assert(err, jc.ErrorIsNil)

	cred := cloud.NewCredential(
		cloud.JSONFileAuthType,
		map[string]string{
			"file": filename,
		},
	)
	schema := cloud.CredentialSchema{{
		"file", cloud.CredentialAttr{FilePath: true},
	}}

	readFile := func(path string) ([]byte, error) {
		c.Assert(path, gc.Equals, filename)
		return []byte("file-contents"), nil
	}

	newCred, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.JSONFileAuthType: schema,
	}, readFile)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(newCred.Attributes(), jc.DeepEquals, map[string]string{
		"file": "file-contents",
	})
}

func (s *credentialsSuite) TestFinalizeCredentialRelativeFilePath(c *gc.C) {
	absFilename := filepath.Join(utils.Home(), "filename")
	err := ioutil.WriteFile(absFilename, []byte{}, 0600)
	c.Assert(err, jc.ErrorIsNil)

	cred := cloud.NewCredential(
		cloud.JSONFileAuthType,
		map[string]string{
			"file": "~/filename",
		},
	)
	schema := cloud.CredentialSchema{{
		"file", cloud.CredentialAttr{FilePath: true},
	}}
	readFile := func(path string) ([]byte, error) {
		c.Assert(path, gc.Equals, absFilename)
		return []byte("file-contents"), nil
	}
	newCred, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.JSONFileAuthType: schema,
	}, readFile)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(newCred.Attributes(), jc.DeepEquals, map[string]string{
		"file": "file-contents",
	})
}

func (s *credentialsSuite) TestFinalizeCredentialInvalidFilePath(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.JSONFileAuthType,
		map[string]string{
			"file": filepath.Join(c.MkDir(), "somefile"),
		},
	)
	schema := cloud.CredentialSchema{{
		"file", cloud.CredentialAttr{FilePath: true},
	}}
	_, err := cloud.FinalizeCredential(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.JSONFileAuthType: schema,
	}, nil)
	c.Assert(err, gc.ErrorMatches, "invalid file path: .*")
}

func (s *credentialsSuite) TestRemoveSecrets(c *gc.C) {
	cred := cloud.NewCredential(
		cloud.UserPassAuthType,
		map[string]string{
			"username": "user",
			"password": "secret",
		},
	)
	c.Assert(cred.Revoked, jc.IsFalse)
	schema := cloud.CredentialSchema{{
		"username", cloud.CredentialAttr{},
	}, {
		"password", cloud.CredentialAttr{Hidden: true},
	}}
	sanitisedCred, err := cloud.RemoveSecrets(cred, map[cloud.AuthType]cloud.CredentialSchema{
		cloud.UserPassAuthType: schema,
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(sanitisedCred.Attributes(), jc.DeepEquals, map[string]string{
		"username": "user",
	})
}

func (s *credentialsSuite) TestValidateFileAttrValue(c *gc.C) {
	_, err := cloud.ValidateFileAttrValue("/xyz/nothing.blah")
	c.Assert(err, gc.ErrorMatches, "invalid file path: /xyz/nothing.blah")

	absPathNewFile := filepath.Join(utils.Home(), "new-creds.json")
	err = ioutil.WriteFile(absPathNewFile, []byte("abc"), 0600)
	c.Assert(err, jc.ErrorIsNil)

	absPath, err := cloud.ValidateFileAttrValue("~/new-creds.json")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(absPath, gc.Equals, absPathNewFile)

	_, err = cloud.ValidateFileAttrValue(utils.Home())
	c.Assert(err, gc.ErrorMatches, fmt.Sprintf("file path must be a file: %s", utils.Home()))
}
