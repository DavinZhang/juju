// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"github.com/DavinZhang/juju/api/imagemetadatamanager"
	"github.com/DavinZhang/juju/cmd/modelcmd"
)

type cloudImageMetadataCommandBase struct {
	modelcmd.ModelCommandBase
	modelcmd.IAASOnlyCommand
}

// NewImageMetadataAPI returns a image metadata api for the root api endpoint
// that the environment command returns.
func (c *cloudImageMetadataCommandBase) NewImageMetadataAPI() (*imagemetadatamanager.Client, error) {
	root, err := c.NewAPIRoot()
	if err != nil {
		return nil, err
	}
	return imagemetadatamanager.NewClient(root), nil
}

// MetadataInfo defines the serialization behaviour of image metadata information.
type MetadataInfo struct {
	Source          string `yaml:"source" json:"source"`
	Series          string `yaml:"series" json:"series"`
	Arch            string `yaml:"arch" json:"arch"`
	Region          string `yaml:"region" json:"region"`
	ImageId         string `yaml:"image-id" json:"image-id"`
	Stream          string `yaml:"stream" json:"stream"`
	VirtType        string `yaml:"virt-type,omitempty" json:"virt-type,omitempty"`
	RootStorageType string `yaml:"storage-type,omitempty" json:"storage-type,omitempty"`
}
