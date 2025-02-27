// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charmhub

import (
	"context"
	"encoding/json"
	http "net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/charmhub/path"
	"github.com/DavinZhang/juju/charmhub/transport"
)

type ResourcesSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&ResourcesSuite{})

func (s *ResourcesSuite) TestListResourceRevisions(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	baseURL := MustParseURL(c, "http://api.foo.bar")

	path := path.MakePath(baseURL)
	name := "meshuggah"
	resource := "image"

	restClient := NewMockRESTClient(ctrl)
	s.expectGet(c, restClient, path, name, resource)

	client := NewResourcesClient(path, restClient, &FakeLogger{})
	response, err := client.ListResourceRevisions(context.TODO(), name, resource)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(response, gc.HasLen, 3)
}

func (s *ResourcesSuite) TestListResourceRevisionsFailure(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	baseURL := MustParseURL(c, "http://api.foo.bar")

	path := path.MakePath(baseURL)
	name := "meshuggah"
	resource := "image"

	restClient := NewMockRESTClient(ctrl)
	s.expectGetFailure(restClient)

	client := NewResourcesClient(path, restClient, &FakeLogger{})
	_, err := client.ListResourceRevisions(context.TODO(), name, resource)
	c.Assert(err, gc.Not(jc.ErrorIsNil))
}

func (s *ResourcesSuite) expectGet(c *gc.C, client *MockRESTClient, p path.Path, charm, resource string) {
	namedPath, err := p.Join(charm, resource, "revisions")
	c.Assert(err, jc.ErrorIsNil)

	client.EXPECT().Get(gomock.Any(), namedPath, gomock.Any()).Do(func(_ context.Context, _ path.Path, response *transport.ResourcesResponse) {
		response.Revisions = make([]transport.ResourceRevision, 3)
	}).Return(RESTResponse{}, nil)
}

func (s *ResourcesSuite) expectGetFailure(client *MockRESTClient) {
	client.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(RESTResponse{StatusCode: http.StatusInternalServerError}, errors.Errorf("boom"))
}

func (s *ResourcesSuite) TestListResourceRevisionsRequestPayload(c *gc.C) {
	resourcesResponse := transport.ResourcesResponse{Revisions: []transport.ResourceRevision{
		{Name: "image", Revision: 3, Type: "image"},
		{Name: "image", Revision: 2, Type: "image"},
		{Name: "image", Revision: 1, Type: "image"},
	}}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(resourcesResponse)
		c.Assert(err, jc.ErrorIsNil)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	config := Config{
		URL: server.URL,
	}
	basePath, err := config.BasePath()
	c.Assert(err, jc.ErrorIsNil)

	resourcesPath, err := basePath.Join("resources")
	c.Assert(err, jc.ErrorIsNil)

	apiRequester := NewAPIRequester(DefaultHTTPTransport(&FakeLogger{}), &FakeLogger{})
	restClient := NewHTTPRESTClient(apiRequester, nil)

	client := NewResourcesClient(resourcesPath, restClient, &FakeLogger{})
	response, err := client.ListResourceRevisions(context.TODO(), "wordpress", "image")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(response, gc.DeepEquals, resourcesResponse.Revisions)
}
