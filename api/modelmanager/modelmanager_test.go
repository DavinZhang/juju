// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package modelmanager_test

import (
	"regexp"
	"time"

	"github.com/juju/errors"
	"github.com/juju/names/v4"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/api/base"
	basetesting "github.com/DavinZhang/juju/api/base/testing"
	"github.com/DavinZhang/juju/api/modelmanager"
	apiservererrors "github.com/DavinZhang/juju/apiserver/errors"
	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/core/life"
	"github.com/DavinZhang/juju/core/model"
	"github.com/DavinZhang/juju/core/status"
	"github.com/DavinZhang/juju/environs/config"
	coretesting "github.com/DavinZhang/juju/testing"
)

type modelmanagerSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&modelmanagerSuite{})

func (s *modelmanagerSuite) TestCreateModelBadUser(c *gc.C) {
	client := modelmanager.NewClient(basetesting.BestVersionCaller{})
	_, err := client.CreateModel("mymodel", "not a user", "", "", names.CloudCredentialTag{}, nil)
	c.Assert(err, gc.ErrorMatches, `invalid owner name "not a user"`)
}

func (s *modelmanagerSuite) TestCreateModelBadCloud(c *gc.C) {
	client := modelmanager.NewClient(basetesting.BestVersionCaller{})
	_, err := client.CreateModel("mymodel", "bob", "123!", "", names.CloudCredentialTag{}, nil)
	c.Assert(err, gc.ErrorMatches, `invalid cloud name "123!"`)
}

func (s *modelmanagerSuite) TestCreateModel(c *gc.C) {
	apiCaller := basetesting.APICallerFunc(func(objType string, version int, id, request string, arg, result interface{}) error {
		c.Check(objType, gc.Equals, "ModelManager")
		c.Check(id, gc.Equals, "")
		c.Check(request, gc.Equals, "CreateModel")
		c.Check(arg, jc.DeepEquals, params.ModelCreateArgs{
			Name:        "new-model",
			OwnerTag:    "user-bob",
			Config:      map[string]interface{}{"abc": 123},
			CloudTag:    "cloud-nimbus",
			CloudRegion: "catbus",
		})
		c.Check(result, gc.FitsTypeOf, &params.ModelInfo{})

		out := result.(*params.ModelInfo)
		out.Name = "dowhatimean"
		out.Type = "iaas"
		out.UUID = "youyoueyedee"
		out.ControllerUUID = "youyoueyedeetoo"
		out.ProviderType = "C-123"
		out.DefaultSeries = "M*A*S*H"
		out.CloudTag = "cloud-nimbus"
		out.CloudRegion = "catbus"
		out.OwnerTag = "user-fnord"
		out.Life = "alive"
		return nil
	})

	client := modelmanager.NewClient(apiCaller)
	newModel, err := client.CreateModel(
		"new-model",
		"bob",
		"nimbus",
		"catbus",
		names.CloudCredentialTag{},
		map[string]interface{}{"abc": 123},
	)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(newModel, jc.DeepEquals, base.ModelInfo{
		Name:           "dowhatimean",
		Type:           model.IAAS,
		UUID:           "youyoueyedee",
		ControllerUUID: "youyoueyedeetoo",
		ProviderType:   "C-123",
		DefaultSeries:  "M*A*S*H",
		Cloud:          "nimbus",
		CloudRegion:    "catbus",
		Owner:          "fnord",
		Life:           "alive",
		Status: base.Status{
			Data: make(map[string]interface{}),
		},
		Users:    []base.UserInfo{},
		Machines: []base.Machine{},
	})
}

func (s *modelmanagerSuite) TestListModelsBadUser(c *gc.C) {
	client := modelmanager.NewClient(basetesting.BestVersionCaller{})
	_, err := client.ListModels("not a user")
	c.Assert(err, gc.ErrorMatches, `invalid user name "not a user"`)
}

func (s *modelmanagerSuite) TestListModels(c *gc.C) {
	lastConnection := time.Now()
	apiCaller := basetesting.APICallerFunc(
		func(objType string,
			version int,
			id, req string,
			args, resp interface{},
		) error {
			c.Check(objType, gc.Equals, "ModelManager")
			c.Check(id, gc.Equals, "")
			c.Check(req, gc.Equals, "ListModels")
			c.Check(args, jc.DeepEquals, params.Entity{"user-user@remote"})
			results := resp.(*params.UserModelList)
			results.UserModels = []params.UserModel{{
				Model: params.Model{
					Name:     "yo",
					UUID:     "wei",
					Type:     "caas",
					OwnerTag: "user-user@remote",
				},
				LastConnection: &lastConnection,
			}, {
				Model: params.Model{
					Name:     "sup",
					UUID:     "hazzagarn",
					Type:     "iaas",
					OwnerTag: "user-phyllis@thrace",
				},
			}}
			return nil
		},
	)

	client := modelmanager.NewClient(apiCaller)
	models, err := client.ListModels("user@remote")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(models, jc.DeepEquals, []base.UserModel{{
		Name:           "yo",
		UUID:           "wei",
		Type:           model.CAAS,
		Owner:          "user@remote",
		LastConnection: &lastConnection,
	}, {
		Name:  "sup",
		UUID:  "hazzagarn",
		Type:  model.IAAS,
		Owner: "phyllis@thrace",
	}})
}

func (s *modelmanagerSuite) TestDestroyModel(c *gc.C) {
	true_ := true
	false_ := false
	defaultMin := 1 * time.Minute
	for _, v := range []int{4, 7} {
		s.testDestroyModel(c, v, nil, nil, nil, time.Minute)
		s.testDestroyModel(c, v, nil, &true_, nil, time.Minute)
		s.testDestroyModel(c, v, nil, &true_, &defaultMin, time.Minute)
		s.testDestroyModel(c, v, nil, &false_, nil, time.Minute)
		s.testDestroyModel(c, v, &true_, nil, nil, time.Minute)
		s.testDestroyModel(c, v, &true_, &false_, nil, time.Minute)
		s.testDestroyModel(c, v, &true_, &true_, &defaultMin, time.Minute)
		s.testDestroyModel(c, v, &false_, nil, nil, time.Minute)
		s.testDestroyModel(c, v, &false_, &false_, nil, time.Minute)
		s.testDestroyModel(c, v, &false_, &true_, &defaultMin, time.Minute)
	}
}

func (s *modelmanagerSuite) testDestroyModel(c *gc.C, v int, destroyStorage, force *bool, maxWait *time.Duration, timeout time.Duration) {
	var called bool
	apiCaller := basetesting.BestVersionCaller{
		BestVersion: v,
		APICallerFunc: basetesting.APICallerFunc(
			func(objType string,
				version int,
				id, req string,
				args, resp interface{},
			) error {
				c.Check(objType, gc.Equals, "ModelManager")
				c.Check(id, gc.Equals, "")
				c.Check(req, gc.Equals, "DestroyModels")
				if v == 4 {
					c.Check(args, jc.DeepEquals, params.DestroyModelsParams{
						Models: []params.DestroyModelParams{{
							ModelTag:       coretesting.ModelTag.String(),
							DestroyStorage: destroyStorage,
						}},
					})
				} else {
					c.Check(args, jc.DeepEquals, params.DestroyModelsParams{
						Models: []params.DestroyModelParams{{
							ModelTag:       coretesting.ModelTag.String(),
							DestroyStorage: destroyStorage,
							Force:          force,
							MaxWait:        maxWait,
							Timeout:        &timeout,
						}},
					})
				}
				results := resp.(*params.ErrorResults)
				*results = params.ErrorResults{
					Results: []params.ErrorResult{{}},
				}
				called = true
				return nil
			},
		),
	}
	client := modelmanager.NewClient(apiCaller)
	err := client.DestroyModel(coretesting.ModelTag, destroyStorage, force, maxWait, timeout)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(called, jc.IsTrue)
}

func (s *modelmanagerSuite) TestDestroyModelV3(c *gc.C) {
	var called bool
	apiCaller := basetesting.APICallerFunc(
		func(objType string,
			version int,
			id, req string,
			args, resp interface{},
		) error {
			c.Check(objType, gc.Equals, "ModelManager")
			c.Check(id, gc.Equals, "")
			c.Check(req, gc.Equals, "DestroyModels")
			c.Check(args, jc.DeepEquals, params.Entities{
				Entities: []params.Entity{{coretesting.ModelTag.String()}},
			})
			results := resp.(*params.ErrorResults)
			*results = params.ErrorResults{
				Results: []params.ErrorResult{{}},
			}
			called = true
			return nil
		},
	)
	client := modelmanager.NewClient(apiCaller)
	destroyStorage := true
	err := client.DestroyModel(coretesting.ModelTag, &destroyStorage, nil, nil, time.Minute)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(called, jc.IsTrue)
}

func (s *modelmanagerSuite) TestDestroyModelV3DestroyStorageNotTrue(c *gc.C) {
	client := modelmanager.NewClient(basetesting.BestVersionCaller{})
	for _, destroyStorage := range []*bool{nil, new(bool)} {
		err := client.DestroyModel(coretesting.ModelTag, destroyStorage, nil, nil, time.Minute)
		c.Assert(err, gc.ErrorMatches, "this Juju controller requires destroyStorage to be true")
	}
}

func (s *modelmanagerSuite) TestModelDefaults(c *gc.C) {
	apiCaller := basetesting.BestVersionCaller{
		APICallerFunc: basetesting.APICallerFunc(
			func(objType string,
				version int,
				id, request string,
				a, result interface{},
			) error {
				c.Check(objType, gc.Equals, "ModelManager")
				c.Check(id, gc.Equals, "")
				c.Check(request, gc.Equals, "ModelDefaultsForClouds")
				c.Check(a, jc.DeepEquals, params.Entities{
					Entities: []params.Entity{{Tag: names.NewCloudTag("aws").String()}},
				})
				c.Assert(result, gc.FitsTypeOf, &params.ModelDefaultsResults{})
				results := result.(*params.ModelDefaultsResults)
				results.Results = []params.ModelDefaultsResult{{Config: map[string]params.ModelDefaults{
					"foo": {"bar", "model", []params.RegionDefaults{{
						"dummy-region",
						"dummy-value"}}},
				}}}
				return nil
			},
		), BestVersion: 6}
	client := modelmanager.NewClient(apiCaller)
	result, err := client.ModelDefaults("aws")
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(result, jc.DeepEquals, config.ModelDefaultAttributes{
		"foo": {"bar", "model", []config.RegionDefaultValue{{
			"dummy-region",
			"dummy-value"}}},
	})
}

func (s *modelmanagerSuite) TestModelDefaultsOldVersionFails(c *gc.C) {
	apiCaller := basetesting.BestVersionCaller{
		APICallerFunc: basetesting.APICallerFunc(
			func(objType string,
				version int,
				id, request string,
				a, result interface{},
			) error {
				c.Fail()
				return nil
			},
		), BestVersion: 5}
	client := modelmanager.NewClient(apiCaller)
	_, err := client.ModelDefaults("aws")
	c.Assert(err, gc.ErrorMatches, "model defaults for cloud aws not supported for this version of Juju")
}

func (s *modelmanagerSuite) TestModelDefaultsOldVersion(c *gc.C) {
	apiCaller := basetesting.BestVersionCaller{
		APICallerFunc: basetesting.APICallerFunc(
			func(objType string,
				version int,
				id, request string,
				a, result interface{},
			) error {
				c.Check(objType, gc.Equals, "ModelManager")
				c.Check(id, gc.Equals, "")
				c.Check(request, gc.Equals, "ModelDefaults")
				c.Check(a, gc.IsNil)
				c.Assert(result, gc.FitsTypeOf, &params.ModelDefaultsResult{})
				results := result.(*params.ModelDefaultsResult)
				results.Config = map[string]params.ModelDefaults{
					"foo": {"bar", "model", []params.RegionDefaults{{
						"dummy-region",
						"dummy-value"}},
					}}
				return nil
			},
		), BestVersion: 5,
	}
	client := modelmanager.NewClient(apiCaller)
	result, err := client.ModelDefaults("")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, jc.DeepEquals, config.ModelDefaultAttributes{
		"foo": {"bar", "model", []config.RegionDefaultValue{{
			"dummy-region",
			"dummy-value"}}},
	})
}

func (s *modelmanagerSuite) TestSetModelDefaults(c *gc.C) {
	called := false
	apiCaller := basetesting.APICallerFunc(
		func(objType string,
			version int,
			id, request string,
			a, result interface{},
		) error {
			c.Check(objType, gc.Equals, "ModelManager")
			c.Check(id, gc.Equals, "")
			c.Check(request, gc.Equals, "SetModelDefaults")
			c.Check(a, jc.DeepEquals, params.SetModelDefaults{
				Config: []params.ModelDefaultValues{{
					CloudTag:    "cloud-mycloud",
					CloudRegion: "region",
					Config: map[string]interface{}{
						"some-name":  "value",
						"other-name": true,
					},
				}}})
			c.Assert(result, gc.FitsTypeOf, &params.ErrorResults{})
			*(result.(*params.ErrorResults)) = params.ErrorResults{
				Results: []params.ErrorResult{{Error: nil}},
			}
			called = true
			return nil
		},
	)
	client := modelmanager.NewClient(apiCaller)
	err := client.SetModelDefaults("mycloud", "region", map[string]interface{}{
		"some-name":  "value",
		"other-name": true,
	})
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(called, jc.IsTrue)
}

func (s *modelmanagerSuite) TestUnsetModelDefaults(c *gc.C) {
	called := false
	apiCaller := basetesting.APICallerFunc(
		func(objType string,
			version int,
			id, request string,
			a, result interface{},
		) error {
			c.Check(objType, gc.Equals, "ModelManager")
			c.Check(id, gc.Equals, "")
			c.Check(request, gc.Equals, "UnsetModelDefaults")
			c.Check(a, jc.DeepEquals, params.UnsetModelDefaults{
				Keys: []params.ModelUnsetKeys{{
					CloudTag:    "cloud-mycloud",
					CloudRegion: "region",
					Keys:        []string{"foo", "bar"},
				}}})
			c.Assert(result, gc.FitsTypeOf, &params.ErrorResults{})
			*(result.(*params.ErrorResults)) = params.ErrorResults{
				Results: []params.ErrorResult{{Error: nil}},
			}
			called = true
			return nil
		},
	)
	client := modelmanager.NewClient(apiCaller)
	err := client.UnsetModelDefaults("mycloud", "region", "foo", "bar")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(called, jc.IsTrue)
}

func (s *modelmanagerSuite) TestModelStatus(c *gc.C) {
	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 4,
		APICallerFunc: func(objType string, version int, id, request string, arg, result interface{}) error {
			c.Check(objType, gc.Equals, "ModelManager")
			c.Check(id, gc.Equals, "")
			c.Check(request, gc.Equals, "ModelStatus")
			c.Check(arg, jc.DeepEquals, params.Entities{
				[]params.Entity{
					{Tag: coretesting.ModelTag.String()},
					{Tag: coretesting.ModelTag.String()},
				},
			})
			c.Check(result, gc.FitsTypeOf, &params.ModelStatusResults{})

			out := result.(*params.ModelStatusResults)
			out.Results = []params.ModelStatus{
				{
					ModelTag:           coretesting.ModelTag.String(),
					OwnerTag:           "user-glenda",
					ApplicationCount:   3,
					HostedMachineCount: 2,
					Life:               "alive",
					Machines: []params.ModelMachineInfo{{
						Id:         "0",
						InstanceId: "inst-ance",
						Status:     "pending",
					}},
				},
				{Error: apiservererrors.ServerError(errors.New("model error"))},
			}
			return nil
		},
	}

	client := modelmanager.NewClient(apiCaller)
	results, err := client.ModelStatus(coretesting.ModelTag, coretesting.ModelTag)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results[0], jc.DeepEquals, base.ModelStatus{
		UUID:               coretesting.ModelTag.Id(),
		TotalMachineCount:  1,
		HostedMachineCount: 2,
		ApplicationCount:   3,
		Owner:              "glenda",
		Life:               life.Alive,
		Machines:           []base.Machine{{Id: "0", InstanceId: "inst-ance", Status: "pending"}},
	})
	c.Assert(results[1].Error, gc.ErrorMatches, "model error")
}

func (s *modelmanagerSuite) TestModelStatusEmpty(c *gc.C) {
	apiCaller := basetesting.APICallerFunc(func(objType string, version int, id, request string, arg, result interface{}) error {
		c.Check(objType, gc.Equals, "ModelManager")
		c.Check(id, gc.Equals, "")
		c.Check(request, gc.Equals, "ModelStatus")
		c.Check(result, gc.FitsTypeOf, &params.ModelStatusResults{})

		return nil
	})

	client := modelmanager.NewClient(apiCaller)
	results, err := client.ModelStatus()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, jc.DeepEquals, []base.ModelStatus{})
}

func (s *modelmanagerSuite) TestModelStatusError(c *gc.C) {
	apiCaller := basetesting.APICallerFunc(
		func(objType string, version int, id, request string, args, result interface{}) error {
			return errors.New("model error")
		})
	client := modelmanager.NewClient(apiCaller)
	out, err := client.ModelStatus(coretesting.ModelTag, coretesting.ModelTag)
	c.Assert(err, gc.ErrorMatches, "model error")
	c.Assert(out, gc.IsNil)
}

func createModelSummary() *params.ModelSummary {
	return &params.ModelSummary{
		Name:               "name",
		UUID:               "uuid",
		Type:               "iaas",
		ControllerUUID:     "controllerUUID",
		ProviderType:       "aws",
		DefaultSeries:      "xenial",
		CloudTag:           "cloud-aws",
		CloudRegion:        "us-east-1",
		CloudCredentialTag: "cloudcred-foo_bob_one",
		OwnerTag:           "user-admin",
		Life:               life.Alive,
		Status:             params.EntityStatus{Status: status.Status("active")},
		UserAccess:         params.ModelAdminAccess,
		Counts:             []params.ModelEntityCount{},
	}
}

func (s *modelmanagerSuite) TestListModelSummaries(c *gc.C) {
	userTag := names.NewUserTag("commander")
	testModelInfo := createModelSummary()

	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 4,
		APICallerFunc: func(objType string, version int, id, request string, arg, result interface{}) error {
			c.Check(objType, gc.Equals, "ModelManager")
			c.Check(id, gc.Equals, "")
			c.Check(request, gc.Equals, "ListModelSummaries")
			c.Check(arg, gc.Equals, params.ModelSummariesRequest{
				UserTag: userTag.String(),
				All:     true,
			})
			c.Check(result, gc.FitsTypeOf, &params.ModelSummaryResults{})

			out := result.(*params.ModelSummaryResults)
			out.Results = []params.ModelSummaryResult{
				{Result: testModelInfo},
				{Error: apiservererrors.ServerError(errors.New("model error"))},
			}
			return nil
		},
	}

	client := modelmanager.NewClient(apiCaller)
	results, err := client.ListModelSummaries(userTag.Id(), true)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(results, gc.HasLen, 2)
	c.Assert(results[0], jc.DeepEquals, base.UserModelSummary{Name: testModelInfo.Name,
		UUID:            testModelInfo.UUID,
		Type:            model.IAAS,
		ControllerUUID:  testModelInfo.ControllerUUID,
		ProviderType:    testModelInfo.ProviderType,
		DefaultSeries:   testModelInfo.DefaultSeries,
		Cloud:           "aws",
		CloudRegion:     "us-east-1",
		CloudCredential: "foo/bob/one",
		Owner:           "admin",
		Life:            "alive",
		Status: base.Status{
			Status: status.Active,
			Data:   map[string]interface{}{},
		},
		ModelUserAccess: "admin",
		Counts:          []base.EntityCount{},
	})
	c.Assert(errors.Cause(results[1].Error), gc.ErrorMatches, "model error")
}

func (s *modelmanagerSuite) TestListModelSummariesParsingErrors(c *gc.C) {
	badOwnerInfo := createModelSummary()
	badOwnerInfo.OwnerTag = "owner-user"

	badCloudInfo := createModelSummary()
	badCloudInfo.CloudTag = "not-cloud"

	badCredentialsInfo := createModelSummary()
	badCredentialsInfo.CloudCredentialTag = "not-credential"

	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 4,
		APICallerFunc: func(objType string, version int, id, request string, arg, result interface{}) error {
			out := result.(*params.ModelSummaryResults)
			out.Results = []params.ModelSummaryResult{
				{Result: badOwnerInfo},
				{Result: badCloudInfo},
				{Result: badCredentialsInfo},
			}
			return nil
		},
	}

	client := modelmanager.NewClient(apiCaller)
	results, err := client.ListModelSummaries("commander", true)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(results, gc.HasLen, 3)
	c.Assert(results[0].Error, gc.ErrorMatches, `while parsing model owner tag: "owner-user" is not a valid tag`)
	c.Assert(results[1].Error, gc.ErrorMatches, `while parsing model cloud tag: "not-cloud" is not a valid tag`)
	c.Assert(results[2].Error, gc.ErrorMatches, `while parsing model cloud credential tag: "not-credential" is not a valid tag`)
}

func (s *modelmanagerSuite) TestListModelSummariesInvalidUserIn(c *gc.C) {
	apiCaller := basetesting.APICallerFunc(
		func(objType string, version int, id, request string, args, result interface{}) error {
			return nil
		})
	client := modelmanager.NewClient(apiCaller)
	out, err := client.ListModelSummaries("++)captain", false)
	c.Assert(err, gc.ErrorMatches, regexp.QuoteMeta(`invalid user name "++)captain"`))
	c.Assert(out, gc.IsNil)
}

func (s *modelmanagerSuite) TestListModelSummariesServerError(c *gc.C) {
	apiCaller := basetesting.APICallerFunc(
		func(objType string, version int, id, request string, args, result interface{}) error {
			return errors.New("captain, error")
		})
	client := modelmanager.NewClient(apiCaller)
	out, err := client.ListModelSummaries("captain", false)
	c.Assert(err, gc.ErrorMatches, "captain, error")
	c.Assert(out, gc.IsNil)
}

func (s *modelmanagerSuite) TestChangeModelCredential(c *gc.C) {
	credentialTag := names.NewCloudCredentialTag("foo/bob/bar")
	called := false
	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 5,
		APICallerFunc: func(objType string, version int, id, request string, arg, result interface{}) error {
			c.Check(objType, gc.Equals, "ModelManager")
			c.Check(id, gc.Equals, "")
			c.Check(request, gc.Equals, "ChangeModelCredential")
			c.Check(arg, jc.DeepEquals, params.ChangeModelCredentialsParams{
				[]params.ChangeModelCredentialParams{
					{ModelTag: coretesting.ModelTag.String(), CloudCredentialTag: credentialTag.String()},
				},
			})
			c.Check(result, gc.FitsTypeOf, &params.ErrorResults{})
			called = true
			out := result.(*params.ErrorResults)
			out.Results = []params.ErrorResult{{}}
			return nil
		},
	}

	client := modelmanager.NewClient(apiCaller)
	err := client.ChangeModelCredential(coretesting.ModelTag, credentialTag)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(called, jc.IsTrue)
}

func (s *modelmanagerSuite) TestChangeModelCredentialManyResults(c *gc.C) {
	credentialTag := names.NewCloudCredentialTag("foo/bob/bar")
	called := false
	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 5,
		APICallerFunc: func(objType string, version int, id, request string, arg, result interface{}) error {
			called = true
			out := result.(*params.ErrorResults)
			out.Results = []params.ErrorResult{{}, {}}
			return nil
		},
	}

	client := modelmanager.NewClient(apiCaller)
	err := client.ChangeModelCredential(coretesting.ModelTag, credentialTag)
	c.Assert(err, gc.ErrorMatches, `expected 1 result, got 2`)
	c.Assert(called, jc.IsTrue)
}

func (s *modelmanagerSuite) TestChangeModelCredentialCallFailed(c *gc.C) {
	credentialTag := names.NewCloudCredentialTag("foo/bob/bar")
	called := false
	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 5,
		APICallerFunc: func(objType string, version int, id, request string, arg, result interface{}) error {
			called = true
			return errors.New("failed call")
		},
	}

	client := modelmanager.NewClient(apiCaller)
	err := client.ChangeModelCredential(coretesting.ModelTag, credentialTag)
	c.Assert(err, gc.ErrorMatches, `failed call`)
	c.Assert(called, jc.IsTrue)
}

func (s *modelmanagerSuite) TestChangeModelCredentialUpdateFailed(c *gc.C) {
	credentialTag := names.NewCloudCredentialTag("foo/bob/bar")
	called := false
	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 5,
		APICallerFunc: func(objType string, version int, id, request string, arg, result interface{}) error {
			called = true
			out := result.(*params.ErrorResults)
			out.Results = []params.ErrorResult{{Error: apiservererrors.ServerError(errors.New("update error"))}}
			return nil
		},
	}

	client := modelmanager.NewClient(apiCaller)
	err := client.ChangeModelCredential(coretesting.ModelTag, credentialTag)
	c.Assert(err, gc.ErrorMatches, `update error`)
	c.Assert(called, jc.IsTrue)
}

func (s *modelmanagerSuite) TestChangeModelCredentialV4(c *gc.C) {
	credentialTag := names.NewCloudCredentialTag("foo/bob/bar")
	called := false
	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 4,
		APICallerFunc: func(objType string, version int, id, request string, arg, result interface{}) error {
			called = true
			return nil
		},
	}

	client := modelmanager.NewClient(apiCaller)
	err := client.ChangeModelCredential(coretesting.ModelTag, credentialTag)
	c.Assert(err, gc.ErrorMatches, `ChangeModelCredential in version 4 not implemented`)
	c.Assert(called, jc.IsFalse)
}

type dumpModelSuite struct {
	coretesting.BaseSuite
}

var _ = gc.Suite(&dumpModelSuite{})

func (s *dumpModelSuite) TestDumpModelV3(c *gc.C) {
	expected := map[string]interface{}{
		"model-uuid": "some-uuid",
		"other-key":  "special",
	}
	results := params.StringResults{Results: []params.StringResult{{
		Result: "model-uuid: some-uuid\nother-key: special\n",
	}}}
	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 3,
		APICallerFunc: basetesting.APICallerFunc(
			func(objType string, version int, id, request string, args, result interface{}) error {
				c.Check(objType, gc.Equals, "ModelManager")
				c.Check(request, gc.Equals, "DumpModels")
				c.Check(version, gc.Equals, 3)
				c.Assert(args, gc.DeepEquals, params.DumpModelRequest{
					Entities:   []params.Entity{{coretesting.ModelTag.String()}},
					Simplified: true})
				res, ok := result.(*params.StringResults)
				c.Assert(ok, jc.IsTrue)
				*res = results
				return nil
			}),
	}
	client := modelmanager.NewClient(apiCaller)
	out, err := client.DumpModel(coretesting.ModelTag, true)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(out, jc.DeepEquals, expected)
}

func (s *dumpModelSuite) TestDumpModelV2(c *gc.C) {
	expected := map[string]interface{}{
		"model-uuid": "some-uuid",
		"other-key":  "special",
	}
	results := params.MapResults{Results: []params.MapResult{{
		Result: expected,
	}}}
	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 2,
		APICallerFunc: basetesting.APICallerFunc(
			func(objType string, version int, id, request string, args, result interface{}) error {
				c.Check(objType, gc.Equals, "ModelManager")
				c.Check(request, gc.Equals, "DumpModels")
				c.Check(version, gc.Equals, 2)
				c.Assert(args, gc.DeepEquals, params.Entities{[]params.Entity{{coretesting.ModelTag.String()}}})
				res, ok := result.(*params.MapResults)
				c.Assert(ok, jc.IsTrue)
				*res = results
				return nil
			}),
	}
	client := modelmanager.NewClient(apiCaller)
	out, err := client.DumpModel(coretesting.ModelTag, false)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(out, jc.DeepEquals, expected)
}

func (s *dumpModelSuite) TestDumpModelErrorV3(c *gc.C) {
	results := params.StringResults{Results: []params.StringResult{{
		Error: &params.Error{Message: "fake error"},
	}}}
	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 3,
		APICallerFunc: basetesting.APICallerFunc(func(objType string, version int, id, request string, args, result interface{}) error {
			res, ok := result.(*params.StringResults)
			c.Assert(ok, jc.IsTrue)
			*res = results
			return nil
		}),
	}
	client := modelmanager.NewClient(apiCaller)
	out, err := client.DumpModel(coretesting.ModelTag, false)
	c.Assert(err, gc.ErrorMatches, "fake error")
	c.Assert(out, gc.IsNil)
}

func (s *dumpModelSuite) TestDumpModelErrorV2(c *gc.C) {
	results := params.MapResults{Results: []params.MapResult{{
		Error: &params.Error{Message: "fake error"},
	}}}
	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 2,
		APICallerFunc: basetesting.APICallerFunc(
			func(objType string, version int, id, request string, args, result interface{}) error {
				res, ok := result.(*params.MapResults)
				c.Assert(ok, jc.IsTrue)
				*res = results
				return nil
			}),
	}
	client := modelmanager.NewClient(apiCaller)
	out, err := client.DumpModel(coretesting.ModelTag, false)
	c.Assert(err, gc.ErrorMatches, "fake error")
	c.Assert(out, gc.IsNil)
}

func (s *dumpModelSuite) TestDumpModelDB(c *gc.C) {
	expected := map[string]interface{}{
		"models": []map[string]interface{}{{
			"name": "admin",
			"uuid": "some-uuid",
		}},
		"machines": []map[string]interface{}{{
			"id":   "0",
			"life": 0,
		}},
	}
	results := params.MapResults{Results: []params.MapResult{{
		Result: expected,
	}}}
	apiCaller := basetesting.APICallerFunc(
		func(objType string, version int, id, request string, args, result interface{}) error {
			c.Check(objType, gc.Equals, "ModelManager")
			c.Check(request, gc.Equals, "DumpModelsDB")
			in, ok := args.(params.Entities)
			c.Assert(ok, jc.IsTrue)
			c.Assert(in, gc.DeepEquals, params.Entities{[]params.Entity{{coretesting.ModelTag.String()}}})
			res, ok := result.(*params.MapResults)
			c.Assert(ok, jc.IsTrue)
			*res = results
			return nil
		})
	client := modelmanager.NewClient(apiCaller)
	out, err := client.DumpModelDB(coretesting.ModelTag)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(out, jc.DeepEquals, expected)
}

func (s *dumpModelSuite) TestDumpModelDBError(c *gc.C) {
	results := params.MapResults{Results: []params.MapResult{{
		Error: &params.Error{Message: "fake error"},
	}}}
	apiCaller := basetesting.APICallerFunc(
		func(objType string, version int, id, request string, args, result interface{}) error {
			res, ok := result.(*params.MapResults)
			c.Assert(ok, jc.IsTrue)
			*res = results
			return nil
		})
	client := modelmanager.NewClient(apiCaller)
	out, err := client.DumpModelDB(coretesting.ModelTag)
	c.Assert(err, gc.ErrorMatches, "fake error")
	c.Assert(out, gc.IsNil)
}

type validateUpdateModelSuite struct {
	coretesting.BaseSuite
}

var _ = gc.Suite(&validateUpdateModelSuite{})

func (s *validateUpdateModelSuite) TestValidateModelUpgradeWithWongAPIVersion(c *gc.C) {
	called := false
	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 8,
		APICallerFunc: func(objType string, version int, id, request string, arg, result interface{}) error {
			called = true
			return nil
		},
	}

	client := modelmanager.NewClient(apiCaller)
	err := client.ValidateModelUpgrade(coretesting.ModelTag, false)
	c.Assert(err, gc.ErrorMatches, `ValidateModelUpgrade in version 8 not implemented`)
	c.Assert(called, jc.IsFalse)
}

func (s *validateUpdateModelSuite) TestValidateModelUpgradeWithErrors(c *gc.C) {
	results := params.ErrorResults{Results: []params.ErrorResult{{
		Error: &params.Error{Message: "fake error"},
	}}}

	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 9,
		APICallerFunc: func(objType string, version int, id, request string, args, result interface{}) error {
			res, ok := result.(*params.ErrorResults)
			c.Assert(ok, jc.IsTrue)
			*res = results
			return nil
		},
	}

	client := modelmanager.NewClient(apiCaller)
	err := client.ValidateModelUpgrade(coretesting.ModelTag, false)
	c.Assert(err, gc.ErrorMatches, "fake error")
}

func (s *validateUpdateModelSuite) TestValidateModelUpgrade(c *gc.C) {
	results := params.ErrorResults{Results: []params.ErrorResult{{}}}

	apiCaller := basetesting.BestVersionCaller{
		BestVersion: 9,
		APICallerFunc: func(objType string, version int, id, request string, args, result interface{}) error {
			c.Check(objType, gc.Equals, "ModelManager")
			c.Check(request, gc.Equals, "ValidateModelUpgrades")
			in, ok := args.(params.ValidateModelUpgradeParams)
			c.Assert(ok, jc.IsTrue)
			c.Assert(in, gc.DeepEquals, params.ValidateModelUpgradeParams{
				Models: []params.ValidateModelUpgradeParam{{
					ModelTag: coretesting.ModelTag.String(),
				}},
				Force: true,
			})

			res, ok := result.(*params.ErrorResults)
			c.Assert(ok, jc.IsTrue)
			*res = results
			return nil
		},
	}

	client := modelmanager.NewClient(apiCaller)
	err := client.ValidateModelUpgrade(coretesting.ModelTag, true)
	c.Assert(err, jc.ErrorIsNil)
}
