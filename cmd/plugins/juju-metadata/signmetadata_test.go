// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/cmd/v3"
	"github.com/juju/cmd/v3/cmdtesting"
	"github.com/juju/loggo"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/environs/simplestreams"
	sstesting "github.com/DavinZhang/juju/environs/simplestreams/testing"
	coretesting "github.com/DavinZhang/juju/testing"
)

type SignMetadataSuite struct {
	coretesting.BaseSuite
}

var _ = gc.Suite(&SignMetadataSuite{})

func (s *SignMetadataSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)
	loggo.GetLogger("").SetLogLevel(loggo.INFO)
}

var expectedLoggingOutput = `signing 2 file\(s\) in .*subdir1.*
signing file .*file1\.json.*
signing file .*file2\.json.*
signing 1 file\(s\) in .*subdir2.*
signing file .*file3\.json.*
`

func makeFileNames(topLevel string) []string {
	return []string{
		filepath.Join(topLevel, "subdir1", "file1.json"),
		filepath.Join(topLevel, "subdir1", "file2.json"),
		filepath.Join(topLevel, "subdir1", "subdir2", "file3.json"),
	}
}

func setupJsonFiles(c *gc.C, topLevel string) {
	err := os.MkdirAll(filepath.Join(topLevel, "subdir1", "subdir2"), 0700)
	c.Assert(err, jc.ErrorIsNil)
	content := []byte("hello world")
	filenames := makeFileNames(topLevel)
	for _, filename := range filenames {
		err = ioutil.WriteFile(filename, content, 0644)
		c.Assert(err, jc.ErrorIsNil)
	}
}

func assertSignedFile(c *gc.C, filename string) {
	r, err := os.Open(filename)
	c.Assert(err, jc.ErrorIsNil)
	defer r.Close()
	data, err := simplestreams.DecodeCheckSignature(r, sstesting.SignedMetadataPublicKey)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(data), gc.Equals, "hello world\n")
}

func assertSignedFiles(c *gc.C, topLevel string) {
	filenames := makeFileNames(topLevel)
	for _, filename := range filenames {
		filename = strings.Replace(filename, ".json", ".sjson", -1)
		assertSignedFile(c, filename)
	}
}

func (s *SignMetadataSuite) TestSignMetadata(c *gc.C) {
	topLevel := c.MkDir()
	keyfile := filepath.Join(topLevel, "privatekey.asc")
	err := ioutil.WriteFile(keyfile, []byte(sstesting.SignedMetadataPrivateKey), 0644)
	c.Assert(err, jc.ErrorIsNil)
	setupJsonFiles(c, topLevel)

	ctx := cmdtesting.Context(c)
	code := cmd.Main(
		newSignMetadataCommand(), ctx, []string{"-d", topLevel, "-k", keyfile, "-p", sstesting.PrivateKeyPassphrase})
	c.Assert(code, gc.Equals, 0)
	output := ctx.Stdout.(*bytes.Buffer).String()
	c.Assert(output, gc.Matches, expectedLoggingOutput)
	assertSignedFiles(c, topLevel)
}

func runSignMetadata(c *gc.C, args ...string) error {
	_, err := cmdtesting.RunCommand(c, newSignMetadataCommand(), args...)
	return err
}

func (s *SignMetadataSuite) TestSignMetadataErrors(c *gc.C) {
	err := runSignMetadata(c, "")
	c.Assert(err, gc.ErrorMatches, `directory must be specified`)
	err = runSignMetadata(c, "-d", "foo")
	c.Assert(err, gc.ErrorMatches, `keyfile must be specified`)
}
