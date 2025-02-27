// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	stdtesting "testing"

	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v2"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/environs/filestorage"
	"github.com/DavinZhang/juju/environs/simplestreams"
	"github.com/DavinZhang/juju/environs/storage"
	"github.com/DavinZhang/juju/testing"
)

func TestPackage(t *stdtesting.T) {
	gc.TestingT(t)
}

var _ = gc.Suite(&datasourceSuite{})

type datasourceSuite struct {
	testing.FakeJujuXDGDataHomeSuite
	stor    storage.Storage
	baseURL string
}

func (s *datasourceSuite) SetUpTest(c *gc.C) {
	s.FakeJujuXDGDataHomeSuite.SetUpTest(c)

	storageDir := c.MkDir()
	stor, err := filestorage.NewFileStorageWriter(storageDir)
	c.Assert(err, jc.ErrorIsNil)
	s.stor = stor
	s.baseURL, err = s.stor.URL("")
	c.Assert(err, jc.ErrorIsNil)
}

func (s *datasourceSuite) TestFetch(c *gc.C) {
	sampleData := "hello world"
	s.stor.Put("foo/bar/data.txt", bytes.NewReader([]byte(sampleData)), int64(len(sampleData)))
	ds := storage.NewStorageSimpleStreamsDataSource("test datasource", s.stor, "", simplestreams.DEFAULT_CLOUD_DATA, false)
	rc, url, err := ds.Fetch("foo/bar/data.txt")
	c.Assert(err, jc.ErrorIsNil)
	defer rc.Close()
	c.Assert(url, gc.Equals, s.baseURL+"/foo/bar/data.txt")
	data, err := ioutil.ReadAll(rc)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(data, gc.DeepEquals, []byte(sampleData))
}

func (s *datasourceSuite) TestFetchWithBasePath(c *gc.C) {
	sampleData := "hello world"
	s.stor.Put("base/foo/bar/data.txt", bytes.NewReader([]byte(sampleData)), int64(len(sampleData)))
	ds := storage.NewStorageSimpleStreamsDataSource("test datasource", s.stor, "base", simplestreams.DEFAULT_CLOUD_DATA, false)
	rc, url, err := ds.Fetch("foo/bar/data.txt")
	c.Assert(err, jc.ErrorIsNil)
	defer rc.Close()
	c.Assert(url, gc.Equals, s.baseURL+"/base/foo/bar/data.txt")
	data, err := ioutil.ReadAll(rc)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(data, gc.DeepEquals, []byte(sampleData))
}

func (s *datasourceSuite) TestFetchWithRetry(c *gc.C) {
	stor := &fakeStorage{shouldRetry: true}
	ds := storage.NewStorageSimpleStreamsDataSource("test datasource", stor, "base", simplestreams.DEFAULT_CLOUD_DATA, false)
	ds.SetAllowRetry(true)
	_, _, err := ds.Fetch("foo/bar/data.txt")
	c.Assert(err, gc.ErrorMatches, "an error")
	c.Assert(stor.getName, gc.Equals, "base/foo/bar/data.txt")
	c.Assert(stor.invokeCount, gc.Equals, 10)
}

func (s *datasourceSuite) TestFetchWithNoRetry(c *gc.C) {
	// NB shouldRetry below is true indicating the fake storage is capable of
	// retrying, not that it will retry.
	stor := &fakeStorage{shouldRetry: true}
	ds := storage.NewStorageSimpleStreamsDataSource("test datasource", stor, "base", simplestreams.DEFAULT_CLOUD_DATA, false)
	_, _, err := ds.Fetch("foo/bar/data.txt")
	c.Assert(err, gc.ErrorMatches, "an error")
	c.Assert(stor.getName, gc.Equals, "base/foo/bar/data.txt")
	c.Assert(stor.invokeCount, gc.Equals, 1)
}

func (s *datasourceSuite) TestURL(c *gc.C) {
	sampleData := "hello world"
	s.stor.Put("bar/data.txt", bytes.NewReader([]byte(sampleData)), int64(len(sampleData)))
	ds := storage.NewStorageSimpleStreamsDataSource("test datasource", s.stor, "", simplestreams.DEFAULT_CLOUD_DATA, false)
	url, err := ds.URL("bar")
	c.Assert(err, jc.ErrorIsNil)
	expectedURL, _ := s.stor.URL("bar")
	c.Assert(url, gc.Equals, expectedURL)
}

func (s *datasourceSuite) TestURLWithBasePath(c *gc.C) {
	sampleData := "hello world"
	s.stor.Put("base/bar/data.txt", bytes.NewReader([]byte(sampleData)), int64(len(sampleData)))
	ds := storage.NewStorageSimpleStreamsDataSource("test datasource", s.stor, "base", simplestreams.DEFAULT_CLOUD_DATA, false)
	url, err := ds.URL("bar")
	c.Assert(err, jc.ErrorIsNil)
	expectedURL, _ := s.stor.URL("base/bar")
	c.Assert(url, gc.Equals, expectedURL)
}

var _ = gc.Suite(&storageSuite{})

type storageSuite struct{}

type fakeStorage struct {
	getName     string
	listPrefix  string
	invokeCount int
	shouldRetry bool
}

func (s *fakeStorage) Get(name string) (io.ReadCloser, error) {
	s.getName = name
	s.invokeCount++
	return nil, fmt.Errorf("an error")
}

func (s *fakeStorage) List(prefix string) ([]string, error) {
	s.listPrefix = prefix
	s.invokeCount++
	return nil, fmt.Errorf("an error")
}

func (s *fakeStorage) URL(name string) (string, error) {
	return "", nil
}

func (s *fakeStorage) DefaultConsistencyStrategy() utils.AttemptStrategy {
	// TODO(katco): 2016-08-09: lp:1611427
	return utils.AttemptStrategy{Min: 10}
}

func (s *fakeStorage) ShouldRetry(error) bool {
	return s.shouldRetry
}

func (s *storageSuite) TestGetWithRetry(c *gc.C) {
	stor := &fakeStorage{shouldRetry: true}
	// TODO(katco): 2016-08-09: lp:1611427
	attempt := utils.AttemptStrategy{Min: 5}
	storage.GetWithRetry(stor, "foo", attempt)
	c.Assert(stor.getName, gc.Equals, "foo")
	c.Assert(stor.invokeCount, gc.Equals, 5)
}

func (s *storageSuite) TestGet(c *gc.C) {
	stor := &fakeStorage{shouldRetry: true}
	storage.Get(stor, "foo")
	c.Assert(stor.getName, gc.Equals, "foo")
	c.Assert(stor.invokeCount, gc.Equals, 10)
}

func (s *storageSuite) TestGetNoRetryAllowed(c *gc.C) {
	stor := &fakeStorage{}
	storage.Get(stor, "foo")
	c.Assert(stor.getName, gc.Equals, "foo")
	c.Assert(stor.invokeCount, gc.Equals, 1)
}

func (s *storageSuite) TestListWithRetry(c *gc.C) {
	stor := &fakeStorage{shouldRetry: true}
	// TODO(katco): 2016-08-09: lp:1611427
	attempt := utils.AttemptStrategy{Min: 5}
	storage.ListWithRetry(stor, "foo", attempt)
	c.Assert(stor.listPrefix, gc.Equals, "foo")
	c.Assert(stor.invokeCount, gc.Equals, 5)
}

func (s *storageSuite) TestList(c *gc.C) {
	stor := &fakeStorage{shouldRetry: true}
	storage.List(stor, "foo")
	c.Assert(stor.listPrefix, gc.Equals, "foo")
	c.Assert(stor.invokeCount, gc.Equals, 10)
}

func (s *storageSuite) TestListNoRetryAllowed(c *gc.C) {
	stor := &fakeStorage{}
	storage.List(stor, "foo")
	c.Assert(stor.listPrefix, gc.Equals, "foo")
	c.Assert(stor.invokeCount, gc.Equals, 1)
}
