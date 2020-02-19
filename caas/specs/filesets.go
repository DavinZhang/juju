// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package specs

import (
	"fmt"

	"github.com/juju/errors"
	"k8s.io/apimachinery/pkg/api/resource"
)

// FileSet defines a set of files to mount
// into the container.
type FileSet struct {
	Name         string `json:"name" yaml:"name"`
	MountPath    string `json:"mountPath" yaml:"mountPath"`
	VolumeSource `json:",inline" yaml:",inline"`
}

// Validate validates FileSet.
func (fs *FileSet) Validate() error {
	if fs.Name == "" {
		return errors.New("file set name is missing")
	}
	if fs.MountPath == "" {
		return errors.Errorf("mount path is missing for file set %q", fs.Name)
	}
	if err := fs.VolumeSource.Validate(fs.Name); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// VolumeSource represents the source of a volume to mount.
type VolumeSource struct {
	Files     map[string]string `json:"files" yaml:"files"`
	HostPath  *HostPathVol      `json:"hostPath" yaml:"hostPath"`
	EmptyDir  *EmptyDirVol      `json:"emptyDir" yaml:"emptyDir"`
	ConfigMap *ResourceRefVol   `json:"configMap" yaml:"configMap"`
	Secret    *ResourceRefVol   `json:"secret" yaml:"secret"`
}

type validator interface {
	Validate(string) error
}

// Validate validates VolumeSource.
func (vs VolumeSource) Validate(name string) error {
	nonNilSource := 0
	if vs.Files != nil {
		nonNilSource++
	}
	if vs.HostPath != nil {
		nonNilSource++
		if err := vs.HostPath.Validate(name); err != nil {
			return errors.Trace(err)
		}
	}
	if vs.EmptyDir != nil {
		nonNilSource++
		if err := vs.EmptyDir.Validate(name); err != nil {
			return errors.Trace(err)
		}
	}
	if vs.Secret != nil {
		nonNilSource++
		if err := vs.Secret.Validate(name); err != nil {
			return errors.Trace(err)
		}
	}
	if vs.ConfigMap != nil {
		nonNilSource++
		if err := vs.ConfigMap.Validate(name); err != nil {
			return errors.Trace(err)
		}
	}
	if nonNilSource == 0 {
		return errors.NewNotValid(nil, fmt.Sprintf("file set %q requires volume source", name))
	}
	if nonNilSource > 1 {
		return errors.NewNotValid(nil, fmt.Sprintf("file set %q can only have one volume source", name))
	}
	return nil
}

// HostPathVol represents a host path mapped into a pod.
type HostPathVol struct {
	Path string `json:"path" yaml:"path"`
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

// Validate validates HostPathVol.
func (hpv *HostPathVol) Validate(name string) error {
	if hpv.Path == "" {
		return errors.Errorf("Path is missing for %q", name)
	}
	return nil
}

// EmptyDirVol represents an empty directory for a pod.
type EmptyDirVol struct {
	Medium    string             `json:"medium,omitempty" yaml:"medium,omitempty"`
	SizeLimit *resource.Quantity `json:"sizeLimit,omitempty" yaml:"sizeLimit,omitempty"`
}

// Validate validates EmptyDirVol.
func (edv *EmptyDirVol) Validate(name string) error {
	return nil
}

// ResourceRefVol reprents a configmap or secret source could be referenced by a volume.
type ResourceRefVol struct {
	Name        string      `json:"name" yaml:"name"`
	Items       []KeyToPath `json:"items,omitempty" yaml:"items,omitempty"`
	DefaultMode *int32      `json:"defaultMode,omitempty" yaml:"defaultMode,omitempty"`
	Optional    *bool       `json:"optional,omitempty" yaml:"optional,omitempty"`
}

// Validate validates ResourceRefVol.
func (rrv *ResourceRefVol) Validate(name string) error {
	if rrv.Name == "" {
		return errors.Errorf("Name is missing for %q", name)
	}
	for _, item := range rrv.Items {
		if err := item.Validate(name); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

// KeyToPath maps a string key to a path within a volume.
type KeyToPath struct {
	Key  string `json:"key" yaml:"key"`
	Path string `json:"path" yaml:"path"`
	Mode *int32 `json:"mode,omitempty" yaml:"mode,omitempty"`
}

// Validate validates KeyToPath.
func (ktp *KeyToPath) Validate(name string) error {
	if ktp.Key == "" {
		return errors.Errorf("Key is missing for %q", name)
	}
	if ktp.Path == "" {
		return errors.Errorf("Path is missing for %q", name)
	}
	return nil
}
