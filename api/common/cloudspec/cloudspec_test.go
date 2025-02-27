// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package cloudspec_test

import (
	"errors"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	apitesting "github.com/DavinZhang/juju/api/base/testing"
	"github.com/DavinZhang/juju/api/common/cloudspec"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/cloud"
	environscloudspec "github.com/DavinZhang/juju/environs/cloudspec"
	coretesting "github.com/DavinZhang/juju/testing"
)

var _ = gc.Suite(&CloudSpecSuite{})

type CloudSpecSuite struct {
	testing.IsolationSuite
}

func (s *CloudSpecSuite) TestNewCloudSpecAPI(c *gc.C) {
	api := cloudspec.NewCloudSpecAPI(nil, coretesting.ModelTag)
	c.Check(api, gc.NotNil)
}

func (s *CloudSpecSuite) TestCloudSpec(c *gc.C) {
	facadeCaller := apitesting.StubFacadeCaller{Stub: &testing.Stub{}}
	facadeCaller.FacadeCallFn = func(name string, args, response interface{}) error {
		c.Assert(name, gc.Equals, "CloudSpec")
		c.Assert(args, jc.DeepEquals, params.Entities{Entities: []params.Entity{
			{Tag: coretesting.ModelTag.String()},
		}})
		*(response.(*params.CloudSpecResults)) = params.CloudSpecResults{
			Results: []params.CloudSpecResult{{
				Result: &params.CloudSpec{
					Type:             "type",
					Name:             "name",
					Region:           "region",
					Endpoint:         "endpoint",
					IdentityEndpoint: "identity-endpoint",
					StorageEndpoint:  "storage-endpoint",
					Credential: &params.CloudCredential{
						AuthType:   "auth-type",
						Attributes: map[string]string{"k": "v"},
					},
					CACertificates: []string{coretesting.CACert},
					SkipTLSVerify:  true,
				},
			}},
		}
		return nil
	}
	api := cloudspec.NewCloudSpecAPI(&facadeCaller, coretesting.ModelTag)
	cloudSpec, err := api.CloudSpec()
	c.Assert(err, jc.ErrorIsNil)

	credential := cloud.NewCredential(
		"auth-type",
		map[string]string{"k": "v"},
	)
	c.Assert(cloudSpec, jc.DeepEquals, environscloudspec.CloudSpec{
		Type:             "type",
		Name:             "name",
		Region:           "region",
		Endpoint:         "endpoint",
		IdentityEndpoint: "identity-endpoint",
		StorageEndpoint:  "storage-endpoint",
		Credential:       &credential,
		CACertificates:   []string{coretesting.CACert},
		SkipTLSVerify:    true,
	})
}

func (s *CloudSpecSuite) TestWatchCloudSpecChanges(c *gc.C) {
	called := false
	facadeCaller := apitesting.StubFacadeCaller{Stub: &testing.Stub{}}
	facadeCaller.ReturnRawAPICaller = apitesting.BestVersionCaller{
		APICallerFunc: apitesting.APICallerFunc(
			func(objType string, version int, id, request string, a, response interface{}) error {
				c.Assert(request, gc.Equals, "Next")
				return nil
			}),
		BestVersion: 1}
	facadeCaller.FacadeCallFn = func(name string, args, response interface{}) error {
		c.Assert(name, gc.Equals, "WatchCloudSpecsChanges")
		c.Assert(args, jc.DeepEquals, params.Entities{Entities: []params.Entity{
			{Tag: coretesting.ModelTag.String()},
		}})
		*(response.(*params.NotifyWatchResults)) = params.NotifyWatchResults{
			Results: []params.NotifyWatchResult{{
				NotifyWatcherId: "1",
			}},
		}
		called = true
		return nil
	}
	api := cloudspec.NewCloudSpecAPI(&facadeCaller, coretesting.ModelTag)
	w, err := api.WatchCloudSpecChanges()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(w, gc.NotNil)
	c.Assert(called, jc.IsTrue)
}

func (s *CloudSpecSuite) TestCloudSpecOverallError(c *gc.C) {
	expect := errors.New("bewm")
	facadeCaller := apitesting.StubFacadeCaller{Stub: &testing.Stub{}}
	facadeCaller.FacadeCallFn = func(name string, args, response interface{}) error {
		return expect
	}
	api := cloudspec.NewCloudSpecAPI(&facadeCaller, coretesting.ModelTag)
	_, err := api.CloudSpec()
	c.Assert(err, gc.Equals, expect)
}

func (s *CloudSpecSuite) TestCloudSpecResultCountMismatch(c *gc.C) {
	facadeCaller := apitesting.StubFacadeCaller{Stub: &testing.Stub{}}
	facadeCaller.FacadeCallFn = func(name string, args, response interface{}) error {
		return nil
	}
	api := cloudspec.NewCloudSpecAPI(&facadeCaller, coretesting.ModelTag)
	_, err := api.CloudSpec()
	c.Assert(err, gc.ErrorMatches, "expected 1 result, got 0")
}

func (s *CloudSpecSuite) TestCloudSpecResultError(c *gc.C) {
	facadeCaller := apitesting.StubFacadeCaller{Stub: &testing.Stub{}}
	facadeCaller.FacadeCallFn = func(name string, args, response interface{}) error {
		*(response.(*params.CloudSpecResults)) = params.CloudSpecResults{
			Results: []params.CloudSpecResult{{
				Error: &params.Error{
					Code:    params.CodeUnauthorized,
					Message: "dang",
				},
			}},
		}
		return nil
	}
	api := cloudspec.NewCloudSpecAPI(&facadeCaller, coretesting.ModelTag)
	_, err := api.CloudSpec()
	c.Assert(err, jc.Satisfies, params.IsCodeUnauthorized)
	c.Assert(err, gc.ErrorMatches, "API request failed: dang")
}

func (s *CloudSpecSuite) TestCloudSpecInvalidCloudSpec(c *gc.C) {
	facadeCaller := apitesting.StubFacadeCaller{Stub: &testing.Stub{}}
	facadeCaller.FacadeCallFn = func(name string, args, response interface{}) error {
		*(response.(*params.CloudSpecResults)) = params.CloudSpecResults{Results: []params.CloudSpecResult{{
			Result: &params.CloudSpec{
				Type: "",
			},
		}}}
		return nil
	}
	api := cloudspec.NewCloudSpecAPI(&facadeCaller, coretesting.ModelTag)
	_, err := api.CloudSpec()
	c.Assert(err, gc.ErrorMatches, "validating CloudSpec: empty Type not valid")
}
