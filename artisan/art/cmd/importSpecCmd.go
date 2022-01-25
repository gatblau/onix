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

// ImportSpecCmd Import the contents from a tarball to create an artisan package in the local registry
type ImportSpecCmd struct {
	cmd       *cobra.Command
	creds     string
	localPath string
}

func NewImportSpecCmd() *ImportSpecCmd {
	c := &ImportSpecCmd{
		cmd: &cobra.Command{
			Use:   "spec [FLAGS] URI",
			Short: "import one or more tarball files into the local registry using a specification",
			Long: `Usage: art import spec [FLAGS] URI

Import one or more tarball files into the local registry using a specification (spec.yaml file).
A specification defines the list of all artefacts required to run an application (i.e. packages and images).
If packages contains OS dependencies, then the specification contains all resources is needed to deploy an app without network access.
The specification is a YAML file called spec.yaml that is located in the same place as the tarball files to import.
The spec.yaml specifies a list of packages as two separate key-value maps, namely packages and images as follows:

spec.yaml
---
# the version for the specification
version: 1.0

# key-value pair of artisan packages 
packages:
  TEST_TESTPK: localhost:8082/test/testpk:v3

# key-value pair of container images
images:
  QUAY_MINIO_LATEST: quay.io/minio/minio:latest
  POSTGRES: postgres:13
...

Examples:
   # import a specification from a folder
   art import spec ./test
   
   # import a specification from an S3 endpoint
   art import spec s3s://my-s3-service.com/my-app/v1.0 -c S3_ID:S3_SECRET
 
   note: in the example above, the URI is the path to the folder where the spec.yaml file is located
`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.creds, "creds", "c", "", "the credentials used to retrieve the specification from an endpoint")
	c.cmd.Flags().StringVarP(&c.localPath, "path", "p", "", "if specified, download the spec tarball files to the path")
	return c
}

func (c *ImportSpecCmd) Run(cmd *cobra.Command, args []string) {
	// check a package name has been provided
	if args != nil && len(args) < 1 {
		core.RaiseErr("the URI of the specification is required")
	}
	// import the tar archive(s)
	err := export.ImportSpec(args[0], c.creds, c.localPath)
	core.CheckErr(err, "cannot import spec")
}
