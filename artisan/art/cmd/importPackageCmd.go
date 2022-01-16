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
	"github.com/gatblau/onix/artisan/registry"
	"github.com/spf13/cobra"
	"log"
)

// ImportPackageCmd import one or more packages from a tarball file into the local registry
type ImportPackageCmd struct {
	cmd   *cobra.Command
	creds string
}

func NewImportPackageCmd() *ImportPackageCmd {
	c := &ImportPackageCmd{
		cmd: &cobra.Command{
			Use:   "package [FLAGS] TARBALL [TARBALL...]",
			Short: "import packages from one or more tarball files into the local registry",
			Long: `Usage: art import package [FLAGS] TARBALL [TARBALL...]

Import packages from one or more tarball files into the local registry.
The tarball file(s) can be located either in the filesystem or an S3 endpoint.

Examples:
   # import one or more packages in the tarball file from the filesystem
   art import package ./test/archive.tar 
   
   # import one or more packages in the tarball file from an S3 endpoint
   art import package s3s://endpoint/bucket/archive.tar -u S3_ID:S3_SECRET
`,
		},
	}
	c.cmd.Run = c.Run
	c.cmd.Flags().StringVarP(&c.creds, "user", "u", "", "the credentials used to retrieve the tarball from an endpoint")
	return c
}

func (c *ImportPackageCmd) Run(cmd *cobra.Command, args []string) {
	// check a package name has been provided
	if len(args) < 1 {
		log.Fatal("at least the name of one tarball file to import is required")
	}
	// create a local registry
	r := registry.NewLocalRegistry()
	// import the tar archive(s)
	err := r.Import(args, c.creds)
	core.CheckErr(err, err.Error())
}
