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
	"github.com/gatblau/onix/artisan/export"
	"github.com/spf13/cobra"
)

// ExportSpecCmd save one or more packages or images tar archives using a specification of artefacts to export in a yaml file
type ExportSpecCmd struct {
	cmd         *cobra.Command
	srcCreds    string
	targetCreds string
	output      string
}

func NewSaveSpecCmd() *ExportSpecCmd {
	c := &ExportSpecCmd{
		cmd: &cobra.Command{
			Use:   "spec [FLAGS] SPEC-FILE",
			Short: "export one or more packages and / or container images to tar archives defined in a spec.yaml file",
			Long: `Usage: art export spec [FLAGS] SPEC-FILE

Exports one or more packages and / or container images tar archives using a specification of the artefacts to export in a yaml file.
The yaml specification file is as follows:

spec.yaml
---
packages:
	- package-key-1: my-package-1:tag1
	- package-key-2: my-package-2:tag2
	- package-key-3: my-package-3:tag3
images:
	- image-key-1: my-image-1:tag1
	- image-key-1: my-image-2:tag2
	- image-key-1: my-image-2:tag2
...

Note: artefacts require a key-value pair to define them.
The key is used to derived the name of the tar archive produced.

Examples:
   # exports the artefacts defined in the spec.yaml file in the current folder to tar archives in the target folder 
   art export spec -o ./test . 

   # exports defined in the spec.yaml file in the ./v1 folder to tar archives in an authenticated and TLS enabled s3 bucket
   art export spec -u REG_USER:REG_PWD -o s3s://endpoint/bucket -c S3_ID:S3_SECRET ./v1
`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.output, "output", "o", "", "the URI where the tar archive(s) will be saved; URI can be file system (absolute or relative path) or s3 bucket (s3:// or s3s:// using TLS)")
	c.cmd.Flags().StringVarP(&c.srcCreds, "user", "u", "", "the credentials used to pull packages from an authenticated artisan registry, if the packages are not already in the local registry")
	c.cmd.Flags().StringVarP(&c.targetCreds, "creds", "c", "", "the credentials to write packages to a destination, if such destination implements authentication (e.g. s3)")
	return c
}

func (c *ExportSpecCmd) Run(cmd *cobra.Command, args []string) {
	var path string
	// if no spec path is provided assume current folder
	if len(args) == 0 || len(args) > 1 {
		path = "."
	}
	spec, err := export.NewSpec(path)
	core.CheckErr(err, "cannot load spec.yaml")
	core.CheckErr(spec.Export(c.output, c.srcCreds, c.targetCreds), "cannot export spec")
}
