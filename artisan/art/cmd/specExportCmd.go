/*
  Onix Config Manager - Artisan
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package cmd

import (
	"github.com/gatblau/onix/artisan/core"
	. "github.com/gatblau/onix/artisan/release"
	"github.com/spf13/cobra"
	"strings"
)

// SpecExportCmd save one or more packages or images tar archives using a specification of artefacts to export in a yaml file
type SpecExportCmd struct {
	Cmd         *cobra.Command
	srcCreds    string
	targetCreds string
	output      string
	filter      string
}

func NewSpecExportCmd() *SpecExportCmd {
	c := &SpecExportCmd{
		Cmd: &cobra.Command{
			Use:   "export [FLAGS] SPEC-FILE-PATH",
			Short: "export an application release specification",
			Long: `Usage: art spec export [FLAGS] SPEC-FILE-PATH

Exports an application release specification using the information in a spec.yaml file.

The specification YAML file should look like the one in the sample below:

spec.yaml
---
# the application release version to which this specification applies
version: "1.0"

# a list of packages to export
packages:
	- package-key-1: "my-package-1:tag1"
	- package-key-2: "my-package-2:tag2"
	- package-key-3: "my-package-3:tag3"

# a list of container images to export
images:
	- image-key-1: "my-image-1:tag1"
	- image-key-1: "my-image-2:tag2"
	- image-key-1: "my-image-2:tag2"
...

Examples:
   # exports the artefacts defined in the spec.yaml file in the current folder to tar archives in the test folder 
   art spec export spec -o ./test

   # exports defined in the spec.yaml file in the ./v1 folder to tar archives in an authenticated and TLS enabled s3 bucket
   art spec export -u REG_USER:REG_PWD -o s3s://endpoint/bucket -c S3_ID:S3_SECRET ./v1
`,
		},
	}
	c.Cmd.Run = c.Run
	c.Cmd.Flags().StringVarP(&c.output, "output", "o", "", "the URI where the tar archive(s) will be saved; URI can be file system (absolute or relative path) or s3 bucket (s3:// or s3s:// using TLS)")
	c.Cmd.Flags().StringVarP(&c.srcCreds, "user", "u", "", "the credentials used to pull packages from an authenticated artisan registry, if the packages are not already in the local registry")
	c.Cmd.Flags().StringVarP(&c.targetCreds, "creds", "c", "", "the credentials to write packages to a destination, if such destination implements authentication (e.g. s3)")
	c.Cmd.Flags().StringVarP(&c.filter, "filter", "f", "", "a regular expression used to select spec artefacts to be exported; any artefacts not matched by the filter are skipped (e.g. -f \"^quay.*$\")")
	c.Cmd.MarkFlagRequired("output")
	return c
}

func (c *SpecExportCmd) Run(cmd *cobra.Command, args []string) {
	var path string
	// if no spec path is provided assume current folder
	if len(args) == 0 || len(args) > 1 {
		path = "."
	} else {
		path = args[0]
	}
	// checks the path points to a local folder
	if strings.Contains(path, "://") {
		core.RaiseErr("SPEC-FILE-PATH should point to a local folder, instead it was %s", path)
	}
	// load the spec file
	spec, err := NewSpec(path, "")
	core.CheckErr(err, "cannot load spec.yaml")
	core.CheckErr(ExportSpec(
		ExportOptions{
			Specification: spec,
			TargetUri:     c.output,
			SourceCreds:   c.srcCreds,
			TargetCreds:   c.targetCreds,
			Filter:        c.filter,
		}), "cannot export spec")
}
