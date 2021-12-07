// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"fmt"

	"github.com/juju/cmd/v3"
	"github.com/juju/errors"
	"github.com/juju/gnuflag"
	"github.com/juju/loggo"
	"github.com/juju/utils/v2"

	jujucmd "github.com/DavinZhang/juju/cmd"
	"github.com/DavinZhang/juju/environs/filestorage"
	"github.com/DavinZhang/juju/environs/simplestreams"
	"github.com/DavinZhang/juju/environs/storage"
	envtools "github.com/DavinZhang/juju/environs/tools"
	"github.com/DavinZhang/juju/juju/keys"
	"github.com/DavinZhang/juju/juju/osenv"
	coretools "github.com/DavinZhang/juju/tools"
)

func newGenerateAgentsCommand() cmd.Command {
	return &generateAgentsCommand{}
}

// generateAgentsCommand is used to generate simplestreams metadata for juju
// agents.
type generateAgentsCommand struct {
	cmd.CommandBase
	fetch           bool
	metadataDir     string
	stream          string
	clean           bool
	public          bool
	preventFallback bool
}

const generateAgentsDoc = `
generate-agents creates the simplestreams metadata for agents.

This command works by scanning a directory for agent binary tarballs from which 
to generate simplestreams agent metadata. The working directory is specified
using the -d argument (defaults to $JUJU_DATA or if not defined
$XDG_DATA_HOME/juju or if that is not defined ~/.local/share/acpoc). The working 
directory path must contain a "tools" subdirectory. The "tools" directory is 
expected to contain a streams release subdirectory (released, devel, proposed) 
containing agent binary tarballs. The output of the streams metadata will be
written next to the release subdirectory.

The stream for which metadata is generated is specified using the --stream 
parameter (default is "released"). Metadata can be generated for any supported 
stream - released, proposed, testing, devel.

Agent tarballs can are located in either a sub directory called "releases" 
(legacy), or a directory named after the stream. The tarballs are expected to be
in the format of "juju-<version>-<series>-<arch>.tgz" and will not be picked
up for writing the streams. If in the event there is no tarballs, then the
command will fallback to using the public cloud metadata, unless the
prevent-fallback flag is used.

By default, if no --stream argument is provided, metadata for agent binaries in 
the "released" stream is generated by scanning for agent binary tarballs in the 
"releases" directory. By specifying a stream explicitly, agent binary tarballs 
are expected to be located in a directory named after the stream.

Newly generated metadata will be merged with any existing metadata that is already there.
To first remove metadata for the specified stream before generating new metadata,
use the --clean option.

Examples:

# generate metadata for "released":
juju metadata generate-agents -d <workingdir>

# generate metadata for "released":
juju metadata generate-agents -d <workingdir> --stream released

# generate metadata for "proposed":
juju metadata generate-agents -d <workingdir> --stream proposed

# generate metadata for "proposed", first removing existing "proposed" metadata:
juju metadata generate-agents -d <workingdir> --stream proposed --clean
`

func (c *generateAgentsCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:    "generate-agents",
		Purpose: "generate simplestreams agent metadata",
		Doc:     generateAgentsDoc,
		Aliases: []string{"generate-tools"},
	})
}

func (c *generateAgentsCommand) SetFlags(f *gnuflag.FlagSet) {
	f.StringVar(&c.metadataDir, "d", "", "local directory in which to store metadata")
	f.StringVar(&c.stream, "stream", envtools.ReleasedStream,
		"simplestreams stream for which to generate the metadata")
	f.BoolVar(&c.clean, "clean", false,
		"remove any existing metadata for the specified stream before generating new metadata")
	f.BoolVar(&c.public, "public", false,
		"agent binaries are for a public cloud, so generate mirror information")
	f.BoolVar(&c.preventFallback, "prevent-fallback", false,
		"prevent falling back to the public cloud if no binaries are found")
}

func (c *generateAgentsCommand) Run(context *cmd.Context) error {
	writer := loggo.NewMinimumLevelWriter(
		cmd.NewCommandLogWriter("juju.environs.tools", context.Stdout, context.Stderr),
		loggo.INFO)
	_ = loggo.RegisterWriter("toolsmetadata", writer)
	defer func() { _, _ = loggo.RemoveWriter("toolsmetadata") }()
	if c.metadataDir == "" {
		c.metadataDir = osenv.JujuXDGDataHomeDir()
	} else {
		c.metadataDir = context.AbsPath(c.metadataDir)
	}

	sourceStorage, err := filestorage.NewFileStorageReader(c.metadataDir)
	if err != nil {
		return errors.Trace(err)
	}

	ss := simplestreams.NewSimpleStreams(simplestreams.DefaultDataSourceFactory())

	fmt.Fprintf(context.Stdout, "Finding agent binaries in %s for stream %s.\n", c.metadataDir, c.stream)
	toolsList, err := envtools.ReadList(sourceStorage, c.stream, -1, -1)
	if err == envtools.ErrNoTools {
		if c.preventFallback {
			return errors.Trace(err)
		}

		var source string
		source, err = envtools.ToolsURL(envtools.DefaultBaseURL)
		if err != nil {
			return errors.Trace(err)
		}
		toolsList, err = envtools.FindToolsForCloud(ss, makeDataSources(ss, source), simplestreams.CloudSpec{}, []string{c.stream}, -1, -1, coretools.Filter{})
	}
	if err != nil {
		return errors.Trace(err)
	}

	targetStorage, err := filestorage.NewFileStorageWriter(c.metadataDir)
	if err != nil {
		return errors.Trace(err)
	}
	writeMirrors := envtools.DoNotWriteMirrors
	if c.public {
		writeMirrors = envtools.WriteMirrors
	}
	return errors.Trace(mergeAndWriteMetadata(ss, targetStorage, c.stream, c.stream, c.clean, toolsList, writeMirrors))
}

func makeDataSources(ss *simplestreams.Simplestreams, urls ...string) []simplestreams.DataSource {
	dataSources := make([]simplestreams.DataSource, len(urls))
	for i, url := range urls {
		dataSources[i] = ss.NewDataSource(
			simplestreams.Config{
				Description:          "local source",
				BaseURL:              url,
				PublicSigningKey:     keys.JujuPublicKey,
				HostnameVerification: utils.VerifySSLHostnames,
				Priority:             simplestreams.CUSTOM_CLOUD_DATA,
			},
		)
	}
	return dataSources
}

// This is essentially the same as tools.MergeAndWriteMetadata, but also
// resolves metadata for existing agents by fetching them and computing
// size/sha256 locally.
func mergeAndWriteMetadata(ss envtools.SimplestreamsFetcher,
	stor storage.Storage, toolsDir, stream string, clean bool, toolsList coretools.List, writeMirrors envtools.ShouldWriteMirrors,
) error {
	existing, err := envtools.ReadAllMetadata(ss, stor)
	if err != nil {
		return err
	}
	if clean {
		delete(existing, stream)
	}
	metadata := envtools.MetadataFromTools(toolsList, toolsDir)
	var mergedMetadata []*envtools.ToolsMetadata
	if mergedMetadata, err = envtools.MergeMetadata(metadata, existing[stream]); err != nil {
		return err
	}
	if err = envtools.ResolveMetadata(stor, toolsDir, mergedMetadata); err != nil {
		return err
	}
	existing[stream] = mergedMetadata
	return envtools.WriteMetadata(stor, existing, []string{stream}, writeMirrors)
}
