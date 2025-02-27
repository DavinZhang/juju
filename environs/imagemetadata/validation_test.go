// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package imagemetadata_test

import (
	"path"
	"path/filepath"

	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v2"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/environs/filestorage"
	"github.com/DavinZhang/juju/environs/imagemetadata"
	"github.com/DavinZhang/juju/environs/simplestreams"
	sstesting "github.com/DavinZhang/juju/environs/simplestreams/testing"
	"github.com/DavinZhang/juju/testing"
)

type ValidateSuite struct {
	testing.BaseSuite
	metadataDir string
}

var _ = gc.Suite(&ValidateSuite{})

func (s *ValidateSuite) makeLocalMetadata(c *gc.C, ss *simplestreams.Simplestreams, id, region, series, endpoint, stream string) {
	metadata := []*imagemetadata.ImageMetadata{
		{
			Id:     id,
			Arch:   "amd64",
			Stream: stream,
		},
	}
	cloudSpec := simplestreams.CloudSpec{
		Region:   region,
		Endpoint: endpoint,
	}
	targetStorage, err := filestorage.NewFileStorageWriter(s.metadataDir)
	c.Assert(err, jc.ErrorIsNil)
	err = imagemetadata.MergeAndWriteMetadata(ss, series, metadata, &cloudSpec, targetStorage)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *ValidateSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)
	s.metadataDir = c.MkDir()
}

func (s *ValidateSuite) assertMatch(c *gc.C, ss *simplestreams.Simplestreams, stream string) {
	s.makeLocalMetadata(c, ss, "1234", "region-2", "raring", "some-auth-url", stream)
	metadataPath := filepath.Join(s.metadataDir, "images")
	params := &simplestreams.MetadataLookupParams{
		Region:        "region-2",
		Release:       "raring",
		Architectures: []string{"amd64"},
		Endpoint:      "some-auth-url",
		Stream:        stream,
		Sources: []simplestreams.DataSource{
			sstesting.VerifyDefaultCloudDataSource("test", utils.MakeFileURL(metadataPath))},
	}
	imageIds, resolveInfo, err := imagemetadata.ValidateImageMetadata(ss, params)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(imageIds, gc.DeepEquals, []string{"1234"})
	c.Check(resolveInfo, gc.DeepEquals, &simplestreams.ResolveInfo{
		Source:    "test",
		Signed:    false,
		IndexURL:  utils.MakeFileURL(path.Join(metadataPath, "streams/v1/index.json")),
		MirrorURL: "",
	})
}

func (s *ValidateSuite) TestMatch(c *gc.C) {
	ss := simplestreams.NewSimpleStreams(sstesting.TestDataSourceFactory())
	s.assertMatch(c, ss, "")
	s.assertMatch(c, ss, imagemetadata.ReleasedStream)
	s.assertMatch(c, ss, "daily")
}

func (s *ValidateSuite) assertNoMatch(c *gc.C, ss *simplestreams.Simplestreams, stream string) {
	s.makeLocalMetadata(c, ss, "1234", "region-2", "raring", "some-auth-url", stream)
	params := &simplestreams.MetadataLookupParams{
		Region:        "region-2",
		Release:       "precise",
		Architectures: []string{"amd64"},
		Endpoint:      "some-auth-url",
		Stream:        stream,
		Sources: []simplestreams.DataSource{
			sstesting.VerifyDefaultCloudDataSource("test", "file://"+s.metadataDir)},
	}
	_, _, err := imagemetadata.ValidateImageMetadata(ss, params)
	c.Assert(err, gc.Not(gc.IsNil))
}

func (s *ValidateSuite) TestNoMatch(c *gc.C) {
	ss := simplestreams.NewSimpleStreams(sstesting.TestDataSourceFactory())
	s.assertNoMatch(c, ss, "")
	s.assertNoMatch(c, ss, imagemetadata.ReleasedStream)
	s.assertNoMatch(c, ss, "daily")
}
