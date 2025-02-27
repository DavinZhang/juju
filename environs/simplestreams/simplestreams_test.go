// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package simplestreams_test

import (
	"strings"
	"testing"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/environs/simplestreams"
	sstesting "github.com/DavinZhang/juju/environs/simplestreams/testing"
)

func Test(t *testing.T) {
	registerSimpleStreamsTests()
	gc.Suite(&jsonSuite{})
	gc.TestingT(t)
}

func registerSimpleStreamsTests() {
	gc.Suite(&simplestreamsSuite{
		LocalLiveSimplestreamsSuite: sstesting.LocalLiveSimplestreamsSuite{
			Source:         sstesting.VerifyDefaultCloudDataSource("test", "test:"),
			RequireSigned:  false,
			DataType:       "image-ids",
			StreamsVersion: "v1",
			ValidConstraint: sstesting.NewTestConstraint(simplestreams.LookupParams{
				CloudSpec: simplestreams.CloudSpec{
					Region:   "us-east-1",
					Endpoint: "https://ec2.us-east-1.amazonaws.com",
				},
				Releases: []string{"precise"},
				Arches:   []string{"amd64", "arm"},
			}),
		},
	})
}

type simplestreamsSuite struct {
	sstesting.TestDataSuite
	sstesting.LocalLiveSimplestreamsSuite
}

func (s *simplestreamsSuite) SetUpSuite(c *gc.C) {
	s.LocalLiveSimplestreamsSuite.SetUpSuite(c)
	s.TestDataSuite.SetUpSuite(c)
}

func (s *simplestreamsSuite) TearDownSuite(c *gc.C) {
	s.TestDataSuite.TearDownSuite(c)
	s.LocalLiveSimplestreamsSuite.TearDownSuite(c)
}

func (s *simplestreamsSuite) TestGetProductsPath(c *gc.C) {
	indexRef, err := s.GetIndexRef(sstesting.Index_v1)
	c.Assert(err, jc.ErrorIsNil)
	path, err := indexRef.GetProductsPath(s.ValidConstraint)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(path, gc.Equals, "streams/v1/image_metadata.json")
}

func (*simplestreamsSuite) TestExtractCatalogsForProductsAcceptsNil(c *gc.C) {
	empty := simplestreams.CloudMetadata{}
	c.Check(simplestreams.ExtractCatalogsForProducts(empty, nil), gc.HasLen, 0)
}

func (*simplestreamsSuite) TestExtractCatalogsForProductsReturnsMatch(c *gc.C) {
	metadata := simplestreams.CloudMetadata{
		Products: map[string]simplestreams.MetadataCatalog{
			"foo": {},
		},
	}
	c.Check(
		simplestreams.ExtractCatalogsForProducts(metadata, []string{"foo"}),
		gc.DeepEquals,
		[]simplestreams.MetadataCatalog{metadata.Products["foo"]})
}

func (*simplestreamsSuite) TestExtractCatalogsForProductsIgnoresNonMatches(c *gc.C) {
	metadata := simplestreams.CloudMetadata{
		Products: map[string]simplestreams.MetadataCatalog{
			"one-product": {},
		},
	}
	absentProducts := []string{"another-product"}
	c.Check(simplestreams.ExtractCatalogsForProducts(metadata, absentProducts), gc.HasLen, 0)
}

func (*simplestreamsSuite) TestExtractCatalogsForProductsPreservesOrder(c *gc.C) {
	products := map[string]simplestreams.MetadataCatalog{
		"1": {},
		"2": {},
		"3": {},
		"4": {},
	}

	metadata := simplestreams.CloudMetadata{Products: products}

	c.Check(
		simplestreams.ExtractCatalogsForProducts(metadata, []string{"1", "3", "4", "2"}),
		gc.DeepEquals,
		[]simplestreams.MetadataCatalog{
			products["1"],
			products["3"],
			products["4"],
			products["2"],
		})
}

func (*simplestreamsSuite) TestExtractIndexesAcceptsEmpty(c *gc.C) {
	ind := simplestreams.Indices{}
	c.Check(simplestreams.ExtractIndexes(ind, nil), gc.HasLen, 0)
}

func (*simplestreamsSuite) TestExtractIndexesReturnsIndex(c *gc.C) {
	metadata := simplestreams.IndexMetadata{}
	ind := simplestreams.Indices{Indexes: map[string]*simplestreams.IndexMetadata{"foo": &metadata}}
	c.Check(simplestreams.ExtractIndexes(ind, nil), gc.DeepEquals, simplestreams.IndexMetadataSlice{&metadata})
}

func (*simplestreamsSuite) TestExtractIndexesReturnsAllIndexes(c *gc.C) {
	ind := simplestreams.Indices{
		Indexes: map[string]*simplestreams.IndexMetadata{
			"foo": {},
			"bar": {},
		},
	}

	array := simplestreams.ExtractIndexes(ind, nil)

	c.Assert(array, gc.HasLen, len(ind.Indexes))
	c.Check(array[0], gc.NotNil)
	c.Check(array[1], gc.NotNil)
	c.Check(array[0], gc.Not(gc.Equals), array[1])
	c.Check(
		(array[0] == ind.Indexes["foo"]),
		gc.Not(gc.Equals),
		(array[1] == ind.Indexes["foo"]))
	c.Check(
		(array[0] == ind.Indexes["bar"]),
		gc.Not(gc.Equals),
		(array[1] == ind.Indexes["bar"]))
}

func (*simplestreamsSuite) TestExtractIndexesReturnsSpecifiedIndexes(c *gc.C) {
	ind := simplestreams.Indices{
		Indexes: map[string]*simplestreams.IndexMetadata{
			"foo":    {},
			"bar":    {},
			"foobar": {},
		},
	}

	array := simplestreams.ExtractIndexes(ind, []string{"foobar"})
	c.Assert(array, gc.HasLen, 1)
	c.Assert(array[0], gc.Equals, ind.Indexes["foobar"])
}

func (*simplestreamsSuite) TestHasCloudAcceptsNil(c *gc.C) {
	metadata := simplestreams.IndexMetadata{Clouds: nil}
	c.Check(simplestreams.HasCloud(metadata, simplestreams.CloudSpec{}), jc.IsTrue)
}

func (*simplestreamsSuite) TestHasCloudFindsMatch(c *gc.C) {
	metadata := simplestreams.IndexMetadata{
		Clouds: []simplestreams.CloudSpec{
			{Region: "r1", Endpoint: "http://e1"},
			{Region: "r2", Endpoint: "http://e2"},
		},
	}
	c.Check(simplestreams.HasCloud(metadata, metadata.Clouds[1]), jc.IsTrue)
}

func (*simplestreamsSuite) TestHasCloudFindsMatchWithTrailingSlash(c *gc.C) {
	metadata := simplestreams.IndexMetadata{
		Clouds: []simplestreams.CloudSpec{
			{Region: "r1", Endpoint: "http://e1/"},
			{Region: "r2", Endpoint: "http://e2"},
		},
	}
	spec := simplestreams.CloudSpec{Region: "r1", Endpoint: "http://e1"}
	c.Check(simplestreams.HasCloud(metadata, spec), jc.IsTrue)
	spec = simplestreams.CloudSpec{Region: "r1", Endpoint: "http://e1/"}
	c.Check(simplestreams.HasCloud(metadata, spec), jc.IsTrue)
	spec = simplestreams.CloudSpec{Region: "r2", Endpoint: "http://e2/"}
	c.Check(simplestreams.HasCloud(metadata, spec), jc.IsTrue)
}

func (*simplestreamsSuite) TestHasCloudReturnsFalseIfCloudsDoNotMatch(c *gc.C) {
	metadata := simplestreams.IndexMetadata{
		Clouds: []simplestreams.CloudSpec{
			{Region: "r1", Endpoint: "http://e1"},
			{Region: "r2", Endpoint: "http://e2"},
		},
	}
	otherCloud := simplestreams.CloudSpec{Region: "r9", Endpoint: "http://e9"}
	c.Check(simplestreams.HasCloud(metadata, otherCloud), jc.IsFalse)
}

func (*simplestreamsSuite) TestHasCloudRequiresIdenticalRegion(c *gc.C) {
	metadata := simplestreams.IndexMetadata{
		Clouds: []simplestreams.CloudSpec{
			{Region: "around", Endpoint: "http://nearby"},
		},
	}
	similarCloud := metadata.Clouds[0]
	similarCloud.Region = "elsewhere"
	c.Assert(similarCloud, gc.Not(gc.Equals), metadata.Clouds[0])

	c.Check(simplestreams.HasCloud(metadata, similarCloud), jc.IsFalse)
}

func (*simplestreamsSuite) TestHasCloudRequiresIdenticalEndpoint(c *gc.C) {
	metadata := simplestreams.IndexMetadata{
		Clouds: []simplestreams.CloudSpec{
			{Region: "around", Endpoint: "http://nearby"},
		},
	}
	similarCloud := metadata.Clouds[0]
	similarCloud.Endpoint = "http://far"
	c.Assert(similarCloud, gc.Not(gc.Equals), metadata.Clouds[0])

	c.Check(simplestreams.HasCloud(metadata, similarCloud), jc.IsFalse)
}

func (*simplestreamsSuite) TestHasProductAcceptsNils(c *gc.C) {
	metadata := simplestreams.IndexMetadata{}
	c.Check(simplestreams.HasProduct(metadata, nil), jc.IsFalse)
}

func (*simplestreamsSuite) TestHasProductFindsMatchingProduct(c *gc.C) {
	metadata := simplestreams.IndexMetadata{ProductIds: []string{"x", "y", "z"}}
	c.Check(
		simplestreams.HasProduct(metadata, []string{"a", "b", metadata.ProductIds[1]}),
		gc.Equals,
		true)
}

func (*simplestreamsSuite) TestHasProductReturnsFalseIfProductsDoNotMatch(c *gc.C) {
	metadata := simplestreams.IndexMetadata{ProductIds: []string{"x", "y", "z"}}
	c.Check(simplestreams.HasProduct(metadata, []string{"a", "b", "c"}), jc.IsFalse)
}

func (*simplestreamsSuite) TestFilterReturnsNothingForEmptyArray(c *gc.C) {
	empty := simplestreams.IndexMetadataSlice{}
	c.Check(
		simplestreams.Filter(empty, func(*simplestreams.IndexMetadata) bool { return true }),
		gc.HasLen,
		0)
}

func (*simplestreamsSuite) TestFilterRemovesNonMatches(c *gc.C) {
	array := simplestreams.IndexMetadataSlice{&simplestreams.IndexMetadata{}}
	c.Check(
		simplestreams.Filter(array, func(*simplestreams.IndexMetadata) bool { return false }),
		gc.HasLen,
		0)
}

func (*simplestreamsSuite) TestFilterIncludesMatches(c *gc.C) {
	metadata := simplestreams.IndexMetadata{}
	array := simplestreams.IndexMetadataSlice{&metadata}
	c.Check(
		simplestreams.Filter(array, func(*simplestreams.IndexMetadata) bool { return true }),
		gc.DeepEquals,
		simplestreams.IndexMetadataSlice{&metadata})
}

func (*simplestreamsSuite) TestFilterLeavesOriginalUnchanged(c *gc.C) {
	item1 := simplestreams.IndexMetadata{CloudName: "aws"}
	item2 := simplestreams.IndexMetadata{CloudName: "openstack"}
	array := simplestreams.IndexMetadataSlice{&item1, &item2}

	result := simplestreams.Filter(array, func(metadata *simplestreams.IndexMetadata) bool {
		return metadata.CloudName == "aws"
	})
	// This exercises both the "leave out" and the "include" code paths.
	c.Assert(result, gc.HasLen, 1)

	// The original, however, has not changed.
	c.Assert(array, gc.HasLen, 2)
	c.Check(array, gc.DeepEquals, simplestreams.IndexMetadataSlice{&item1, &item2})
}

func (*simplestreamsSuite) TestFilterPreservesOrder(c *gc.C) {
	array := simplestreams.IndexMetadataSlice{
		&simplestreams.IndexMetadata{CloudName: "aws"},
		&simplestreams.IndexMetadata{CloudName: "maas"},
		&simplestreams.IndexMetadata{CloudName: "openstack"},
	}

	c.Check(
		simplestreams.Filter(array, func(metadata *simplestreams.IndexMetadata) bool { return true }),
		gc.DeepEquals,
		array)
}

func (*simplestreamsSuite) TestFilterCombinesMatchesAndNonMatches(c *gc.C) {
	array := simplestreams.IndexMetadataSlice{
		&simplestreams.IndexMetadata{Format: "1.0"},
		&simplestreams.IndexMetadata{Format: "1.1"},
		&simplestreams.IndexMetadata{Format: "2.0"},
		&simplestreams.IndexMetadata{Format: "2.1"},
	}

	dotOFormats := simplestreams.Filter(array, func(metadata *simplestreams.IndexMetadata) bool {
		return strings.HasSuffix(metadata.Format, ".0")
	})

	c.Check(dotOFormats, gc.DeepEquals, simplestreams.IndexMetadataSlice{array[0], array[2]})
}

// countingSource is used to check that a DataSource has been queried.
type countingSource struct {
	simplestreams.DataSource
	count int
}

func (s *countingSource) URL(path string) (string, error) {
	s.count++
	return s.DataSource.URL(path)
}

func (s *simplestreamsSuite) TestGetMetadataNoMatching(c *gc.C) {
	source := &countingSource{
		DataSource: sstesting.VerifyDefaultCloudDataSource("test", "test:/daily"),
	}
	sources := []simplestreams.DataSource{source, source, source}
	constraint := sstesting.NewTestConstraint(simplestreams.LookupParams{
		CloudSpec: simplestreams.CloudSpec{
			Region:   "us-east-1",
			Endpoint: "https://ec2.us-east-1.amazonaws.com",
		},
		Releases: []string{"precise"},
		Arches:   []string{"not-a-real-arch"}, // never matches
	})
	params := simplestreams.GetMetadataParams{
		StreamsVersion:   s.StreamsVersion,
		LookupConstraint: constraint,
		ValueParams:      simplestreams.ValueParams{DataType: "image-ids"},
	}

	ss := simplestreams.NewSimpleStreams(sstesting.TestDataSourceFactory())
	items, resolveInfo, err := ss.GetMetadata(sources, params)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(items, gc.HasLen, 0)
	c.Assert(resolveInfo, gc.DeepEquals, &simplestreams.ResolveInfo{
		Source:    "test",
		Signed:    false,
		IndexURL:  "test:/daily/streams/v1/index.json",
		MirrorURL: "",
	})

	// There should be 4 calls to each data-source:
	// one for .sjson, one for .json, repeated for legacy vs new index files.
	c.Assert(source.count, gc.Equals, 4*len(sources))
}

func (s *simplestreamsSuite) TestMetadataCatalog(c *gc.C) {
	metadata := s.AssertGetMetadata(c)
	c.Check(len(metadata.Products), gc.Equals, 6)
	c.Check(len(metadata.Aliases), gc.Equals, 1)
	metadataCatalog := metadata.Products["com.ubuntu.cloud:server:12.04:amd64"]
	c.Check(len(metadataCatalog.Items), gc.Equals, 2)
	c.Check(metadataCatalog.Series, gc.Equals, "precise")
	c.Check(metadataCatalog.Version, gc.Equals, "12.04")
	c.Check(metadataCatalog.Arch, gc.Equals, "amd64")
	c.Check(metadataCatalog.RegionName, gc.Equals, "au-east-1")
	c.Check(metadataCatalog.Endpoint, gc.Equals, "https://somewhere")
}

func (s *simplestreamsSuite) TestItemCollection(c *gc.C) {
	ic := s.AssertGetItemCollections(c, "20121218")
	c.Check(ic.RegionName, gc.Equals, "au-east-2")
	c.Check(ic.Endpoint, gc.Equals, "https://somewhere-else")
	c.Assert(len(ic.Items) > 0, jc.IsTrue)
	ti := ic.Items["usww2he"].(*sstesting.TestItem)
	c.Check(ti.Id, gc.Equals, "ami-442ea674")
	c.Check(ti.Storage, gc.Equals, "ebs")
	c.Check(ti.VirtType, gc.Equals, "hvm")
	c.Check(ti.RegionName, gc.Equals, "us-east-1")
	c.Check(ti.Endpoint, gc.Equals, "https://ec2.us-east-1.amazonaws.com")
}

func (s *simplestreamsSuite) TestDenormalisationFromCollection(c *gc.C) {
	ic := s.AssertGetItemCollections(c, "20121218")
	ti := ic.Items["usww1pe"].(*sstesting.TestItem)
	c.Check(ti.RegionName, gc.Equals, ic.RegionName)
	c.Check(ti.Endpoint, gc.Equals, ic.Endpoint)
}

func (s *simplestreamsSuite) TestDenormalisationFromCatalog(c *gc.C) {
	metadata := s.AssertGetMetadata(c)
	metadataCatalog := metadata.Products["com.ubuntu.cloud:server:12.04:amd64"]
	ic := metadataCatalog.Items["20111111"]
	ti := ic.Items["usww3pe"].(*sstesting.TestItem)
	c.Check(ti.RegionName, gc.Equals, metadataCatalog.RegionName)
	c.Check(ti.Endpoint, gc.Equals, metadataCatalog.Endpoint)
}

func (s *simplestreamsSuite) TestDenormalisationFromTopLevel(c *gc.C) {
	metadata := s.AssertGetMetadata(c)
	metadataCatalog := metadata.Products["com.ubuntu.cloud:server:14.04:amd64"]
	ic := metadataCatalog.Items["20140118"]
	ti := ic.Items["nzww1pe"].(*sstesting.TestItem)
	c.Check(ti.RegionName, gc.Equals, metadata.RegionName)
	c.Check(ti.Endpoint, gc.Equals, metadata.Endpoint)
}

func (s *simplestreamsSuite) TestDealiasing(c *gc.C) {
	metadata := s.AssertGetMetadata(c)
	metadataCatalog := metadata.Products["com.ubuntu.cloud:server:12.04:amd64"]
	ic := metadataCatalog.Items["20121218"]
	ti := ic.Items["usww3he"].(*sstesting.TestItem)
	c.Check(ti.RegionName, gc.Equals, "us-west-3")
	c.Check(ti.Endpoint, gc.Equals, "https://ec2.us-west-3.amazonaws.com")
}

type storageVirtTest struct {
	product, coll, item, storage, virt string
}

func (s *simplestreamsSuite) TestStorageVirtFromTopLevel(c *gc.C) {
	s.assertImageMetadata(c,
		storageVirtTest{"com.ubuntu.cloud:server:13.04:amd64", "20160318", "nzww1pe", "ebs", "pv"},
	)
}

func (s *simplestreamsSuite) TestStorageVirtFromCatalog(c *gc.C) {
	s.assertImageMetadata(c,
		storageVirtTest{"com.ubuntu.cloud:server:14.10:amd64", "20160218", "nzww1pe", "ebs", "pv"},
	)
}

func (s *simplestreamsSuite) TestStorageVirtFromCollection(c *gc.C) {
	s.assertImageMetadata(c,
		storageVirtTest{"com.ubuntu.cloud:server:12.10:amd64", "20160118", "nzww1pe", "ebs", "pv"},
	)
}

func (s *simplestreamsSuite) TestStorageVirtFromItem(c *gc.C) {
	s.assertImageMetadata(c,
		storageVirtTest{"com.ubuntu.cloud:server:14.04:amd64", "20140118", "nzww1pe", "ssd", "hvm"},
	)
}

func (s *simplestreamsSuite) assertImageMetadata(c *gc.C, one storageVirtTest) {
	metadata := s.AssertGetMetadata(c)
	metadataCatalog := metadata.Products[one.product]
	ic := metadataCatalog.Items[one.coll]
	ti := ic.Items[one.item].(*sstesting.TestItem)
	c.Check(ti.Storage, gc.Equals, one.storage)
	c.Check(ti.VirtType, gc.Equals, one.virt)
}

var getMirrorTests = []struct {
	region    string
	endpoint  string
	err       string
	mirrorURL string
	path      string
}{{
	// defaults
	mirrorURL: "http://some-mirror/",
	path:      "com.ubuntu.juju:download.json",
}, {
	// default mirror index entry
	region:    "some-region",
	endpoint:  "https://some-endpoint.com",
	mirrorURL: "http://big-mirror/",
	path:      "big:download.json",
}, {
	// endpoint with trailing "/"
	region:    "some-region",
	endpoint:  "https://some-endpoint.com/",
	mirrorURL: "http://big-mirror/",
	path:      "big:download.json",
}}

func (s *simplestreamsSuite) TestGetMirrorMetadata(c *gc.C) {
	for i, t := range getMirrorTests {
		c.Logf("test %d", i)
		if t.region == "" {
			t.region = "us-east-2"
		}
		if t.endpoint == "" {
			t.endpoint = "https://ec2.us-east-2.amazonaws.com"
		}
		cloud := simplestreams.CloudSpec{t.region, t.endpoint}
		params := simplestreams.ValueParams{
			DataType:        "content-download",
			MirrorContentId: "com.ubuntu.juju:released:agents",
		}
		ss := simplestreams.NewSimpleStreams(sstesting.TestDataSourceFactory())
		indexRef, err := ss.GetIndexWithFormat(
			s.Source, s.IndexPath(), sstesting.Index_v1,
			simplestreams.MirrorsPath("v1"), s.RequireSigned, cloud, params)
		if !c.Check(err, jc.ErrorIsNil) {
			continue
		}
		if t.err != "" {
			c.Check(err, gc.ErrorMatches, t.err)
			continue
		}
		if !c.Check(err, jc.ErrorIsNil) {
			continue
		}
		mirrorURL, err := indexRef.Source.URL("")
		if !c.Check(err, jc.ErrorIsNil) {
			continue
		}
		c.Check(mirrorURL, gc.Equals, t.mirrorURL)
		c.Check(indexRef.MirroredProductsPath, gc.Equals, t.path)
	}
}
